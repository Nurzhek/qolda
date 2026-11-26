package listing

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"qolda/internal/constant"
	"qolda/internal/domain/listing/model"
	domainListingServiceP "qolda/internal/domain/listing/service"
	sessionModel "qolda/internal/domain/session/model"
	sessionServiceP "qolda/internal/domain/session/service"
	usrModel "qolda/internal/domain/usr/model"
	"qolda/internal/errs"
)

// fakeRepo is an in-memory listing RepoDbI under the real domain service
type fakeRepo struct {
	items  map[int64]*model.Main
	lastLR *model.ListReq
}

func (f *fakeRepo) List(_ context.Context, pars *model.ListReq) ([]*model.Main, int64, error) {
	f.lastLR = pars
	result := make([]*model.Main, 0)
	for _, v := range f.items {
		if (pars.Status == nil || v.Status == *pars.Status) && (pars.OwnerId == nil || v.OwnerId == *pars.OwnerId) {
			result = append(result, new(*v))
		}
	}
	return result, int64(len(result)), nil
}

func (f *fakeRepo) Get(_ context.Context, id int64) (*model.Main, bool, error) {
	v, ok := f.items[id]
	if !ok {
		return nil, false, nil
	}
	return new(*v), true, nil
}

func (f *fakeRepo) Create(_ context.Context, obj *model.Edit) (int64, error) {
	id := int64(len(f.items) + 1)
	m := model.Main{Id: id}.Apply(obj)
	f.items[id] = &m
	return id, nil
}

func (f *fakeRepo) Update(_ context.Context, id int64, obj *model.Edit) error {
	m := f.items[id].Apply(obj)
	f.items[id] = &m
	return nil
}

type fakeUsr struct{}

func (fakeUsr) List(_ context.Context, pars *usrModel.ListReq) ([]*usrModel.Main, int64, error) {
	result := make([]*usrModel.Main, 0)
	for _, id := range *pars.IDs {
		result = append(result, &usrModel.Main{Id: id, Name: "user", Campus: "dormitory"})
	}
	return result, int64(len(result)), nil
}

func setup() (*Usecase, *fakeRepo) {
	repo := &fakeRepo{items: map[int64]*model.Main{}}
	svc := domainListingServiceP.New(repo)
	return New(svc, fakeUsr{}, sessionServiceP.New("s", time.Minute, time.Hour)), repo
}

func as(id int64, role string) context.Context {
	return sessionServiceP.WithContext(context.Background(), &sessionModel.Session{Id: id, Role: role})
}

func rent() *model.Edit {
	tomorrow := time.Now().AddDate(0, 0, 1).Format(dateLayout)
	return &model.Edit{
		Type: new("rent"), Title: new("Camping tent, 3 person"), Category: new("sports_outdoor"),
		Photos: new([]string{"/files/abc.jpg"}), Price: new(int64(1500)), Deposit: new(int64(10000)),
		AvailableFrom: new(tomorrow), AvailableTo: new(tomorrow),
	}
}

func TestCreateValidation(t *testing.T) {
	uc, repo := setup()
	ctx := as(1, constant.RoleStudent)

	_, err := uc.Create(context.Background(), rent())
	assert.ErrorIs(t, err, errs.NotAuthorized)

	cases := map[string]func(e *model.Edit){
		"type":           func(e *model.Edit) { e.Type = new("lease") },
		"title":          func(e *model.Edit) { e.Title = new("  ab ") },
		"category":       func(e *model.Edit) { e.Category = new("weapons") },
		"price":          func(e *model.Edit) { e.Price = new(int64(0)) },
		"available_from": func(e *model.Edit) { e.AvailableFrom = new("2026-02-30") },
		"available_to":   func(e *model.Edit) { e.AvailableTo = new("2000-01-01") },
		"photos":         func(e *model.Edit) { e.Photos = new([]string{"https://evil.com/x.png"}) },
		"campus":         func(e *model.Edit) { e.Campus = new("mars") },
	}
	for field, mutate := range cases {
		e := rent()
		mutate(e)
		_, err = uc.Create(ctx, e)
		full, ok := errs.AsErrFull(err)
		require.True(t, ok, field)
		assert.Contains(t, full.Fields, field)
	}

	sale := &model.Edit{Type: new("sale"), Title: new("Calculus textbook"), Category: new("books")}
	_, err = uc.Create(ctx, sale)
	assert.Contains(t, mustFields(t, err), "price")

	exchange := &model.Edit{Type: new("exchange"), Title: new("Kettle"), Category: new("kitchen"), Price: new(int64(999))}
	_, err = uc.Create(ctx, exchange)
	assert.Contains(t, mustFields(t, err), "wanted")

	// exchange: price and deposit are zeroed by normalization, owner and status come from the session
	exchange.Wanted = new("a desk lamp")
	exchange.OwnerId, exchange.Status = new(int64(99)), new("paused")
	id, err := uc.Create(ctx, exchange)
	require.NoError(t, err)
	got := repo.items[id]
	assert.Equal(t, int64(1), got.OwnerId)
	assert.Equal(t, constant.ListingStatusActive, got.Status)
	assert.Zero(t, got.Price)
	assert.Empty(t, got.AvailableFrom)
}

func TestOwnershipAndVisibility(t *testing.T) {
	uc, _ := setup()
	owner, other := as(1, constant.RoleStudent), as(2, constant.RoleStudent)

	id, err := uc.Create(owner, rent())
	require.NoError(t, err)

	// only the owner edits and changes status
	assert.ErrorIs(t, uc.Update(other, id, &model.Edit{Title: new("Mine now")}), errs.NoPermission)
	assert.ErrorIs(t, uc.SetStatus(other, id, "paused"), errs.NoPermission)
	assert.Equal(t, "invalid_request", errs.Code(uc.SetStatus(owner, id, "deleted")))

	// type does not change through Update
	require.NoError(t, uc.Update(owner, id, &model.Edit{Type: new("sale"), Title: new("Tent for 3")}))
	item, err := uc.Get(other, id)
	require.NoError(t, err)
	assert.Equal(t, "rent", item.Type)
	assert.Equal(t, "Tent for 3", item.Title)
	assert.Equal(t, "user", item.Owner.Name)

	// paused: hidden from others, visible to the owner and moderators
	require.NoError(t, uc.SetStatus(owner, id, "paused"))
	_, err = uc.Get(other, id)
	assert.ErrorIs(t, err, errs.ObjectNotFound)
	_, err = uc.Get(owner, id)
	assert.NoError(t, err)
	_, err = uc.Get(as(3, constant.RoleModerator), id)
	assert.NoError(t, err)

	// catalogue: someone else's status filter is ignored, mine returns own listings in any status
	page := func() *model.ListReq { return &model.ListReq{} }
	_, _, err = uc.List(other, page(), false)
	assert.ErrorIs(t, err, errs.IncorrectPageSize)

	lr := &model.ListReq{Status: new("paused")}
	lr.PageSize = 20
	items, _, err := uc.List(other, lr, false)
	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, constant.ListingStatusActive, *lr.Status)

	mine := &model.ListReq{}
	mine.PageSize = 20
	items, _, err = uc.List(owner, mine, true)
	require.NoError(t, err)
	assert.Len(t, items, 1)

	anon := &model.ListReq{}
	anon.PageSize = 20
	_, _, err = uc.List(context.Background(), anon, true)
	assert.ErrorIs(t, err, errs.NotAuthorized)
}

func mustFields(t *testing.T, err error) map[string]string {
	t.Helper()
	full, ok := errs.AsErrFull(err)
	require.True(t, ok, "expected validation error, got %v", err)
	return full.Fields
}
