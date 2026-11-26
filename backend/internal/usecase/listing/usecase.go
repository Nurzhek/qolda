package listing

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/samber/lo"

	"qolda/internal/config"
	"qolda/internal/constant"
	"qolda/internal/domain/listing/model"
	usrModel "qolda/internal/domain/usr/model"
	"qolda/internal/errs"
	"qolda/internal/util"
)

const (
	dateLayout = "2006-01-02"
	maxAmount  = 100_000_000 // KZT
	maxPhotos  = 6
	photoUrl   = "/files/"
)

type Usecase struct {
	svc        ServiceI
	usrSvc     UsrServiceI
	sessionSvc SessionServiceI
}

func New(svc ServiceI, usrSvc UsrServiceI, sessionSvc SessionServiceI) *Usecase {
	return &Usecase{svc: svc, usrSvc: usrSvc, sessionSvc: sessionSvc}
}

// List is the catalogue. Others' listings are visible only when active; mine returns own listings in any status.
func (u *Usecase) List(ctx context.Context, pars *model.ListReq, mine bool) ([]*model.Main, int64, error) {
	if err := util.RequirePageSize(pars.ListParams, 100); err != nil {
		return nil, 0, err
	}
	if mine {
		if !u.sessionSvc.CtxIsAuthed(ctx) {
			return nil, 0, errs.NotAuthorized
		}
		pars.OwnerId = new(u.sessionSvc.CtxGetId(ctx))
	} else if pars.Status == nil || !u.sessionSvc.CtxHasRole(ctx, constant.RoleModerator, constant.RoleAdmin) {
		pars.Status = new(constant.ListingStatusActive)
	}
	if pars.Query != nil {
		pars.Query = lo.EmptyableToPtr(strings.TrimSpace(*pars.Query))
	}
	for name, v := range map[string]*string{"free_from": pars.FreeFrom, "free_to": pars.FreeTo} {
		if v != nil && !isDate(*v) {
			return nil, 0, errs.Field(name, "YYYY-MM-DD")
		}
	}

	items, tCount, err := u.svc.List(ctx, pars)
	if err != nil {
		return nil, 0, fmt.Errorf("svc.List: %w", err)
	}
	if err = u.attachOwners(ctx, items); err != nil {
		return nil, 0, err
	}
	return items, tCount, nil
}

func (u *Usecase) Get(ctx context.Context, id int64) (*model.Main, error) {
	if id == 0 {
		return nil, errs.IdRequired
	}
	item, _, err := u.svc.Get(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("svc.Get: %w", err)
	}
	// an inactive listing does not exist for outsiders
	if item.Status != constant.ListingStatusActive && !u.canManage(ctx, item) &&
		!u.sessionSvc.CtxHasRole(ctx, constant.RoleModerator, constant.RoleAdmin) {
		return nil, errs.ObjectNotFound
	}
	if err = u.attachOwners(ctx, []*model.Main{item}); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *Usecase) Create(ctx context.Context, obj *model.Edit) (int64, error) {
	if !u.sessionSvc.CtxIsAuthed(ctx) {
		return 0, errs.NotAuthorized
	}
	obj.OwnerId = new(u.sessionSvc.CtxGetId(ctx))
	obj.Status = new(constant.ListingStatusActive)

	if err := validate(model.Main{}.Apply(trim(obj)), true); err != nil {
		return 0, err
	}
	newId, err := u.svc.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("svc.Create: %w", err)
	}
	return newId, nil
}

// Update never changes type, owner or status
func (u *Usecase) Update(ctx context.Context, id int64, obj *model.Edit) error {
	current, err := u.getManaged(ctx, id)
	if err != nil {
		return err
	}
	obj.OwnerId, obj.Status = nil, nil
	obj.Type = new(current.Type) // so fields get normalized by type

	if err = validate(current.Apply(trim(obj)), false); err != nil {
		return err
	}
	if err = u.svc.Update(ctx, id, obj); err != nil {
		return fmt.Errorf("svc.Update: %w", err)
	}
	return nil
}

func (u *Usecase) SetStatus(ctx context.Context, id int64, status string) error {
	if !slices.Contains(constant.ListingStatuses, status) {
		return errs.Field("status", "unknown status")
	}
	if _, err := u.getManaged(ctx, id); err != nil {
		return err
	}
	if err := u.svc.Update(ctx, id, &model.Edit{Status: &status}); err != nil {
		return fmt.Errorf("svc.Update: %w", err)
	}
	return nil
}

func (u *Usecase) getManaged(ctx context.Context, id int64) (*model.Main, error) {
	if !u.sessionSvc.CtxIsAuthed(ctx) {
		return nil, errs.NotAuthorized
	}
	if id == 0 {
		return nil, errs.IdRequired
	}
	item, _, err := u.svc.Get(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("svc.Get: %w", err)
	}
	if !u.canManage(ctx, item) {
		return nil, errs.NoPermission
	}
	return item, nil
}

func (u *Usecase) canManage(ctx context.Context, item *model.Main) bool {
	return u.sessionSvc.CtxIsAuthed(ctx) && item.OwnerId == u.sessionSvc.CtxGetId(ctx)
}

// attachOwners loads owner cards for a result page in one query
func (u *Usecase) attachOwners(ctx context.Context, items []*model.Main) error {
	if len(items) == 0 {
		return nil
	}
	ids := lo.Uniq(lo.Map(items, func(v *model.Main, _ int) int64 { return v.OwnerId }))
	owners, _, err := u.usrSvc.List(ctx, &usrModel.ListReq{IDs: &ids})
	if err != nil {
		return fmt.Errorf("usrSvc.List: %w", err)
	}
	byId := lo.KeyBy(owners, func(v *usrModel.Main) int64 { return v.Id })
	for _, item := range items {
		if o, ok := byId[item.OwnerId]; ok {
			item.Owner = &model.Owner{Id: o.Id, Name: o.Name, Campus: o.Campus}
		}
	}
	return nil
}

func trim(obj *model.Edit) *model.Edit {
	for _, p := range []**string{&obj.Title, &obj.Description, &obj.Wanted} {
		if *p != nil {
			*p = new(strings.TrimSpace(**p))
		}
	}
	return obj
}

// validate checks listing type rules (domain requirements) on the resulting state
func validate(m model.Main, forCreate bool) error {
	if !slices.Contains(constant.ListingTypes, m.Type) {
		return errs.Field("type", "rent, exchange or sale")
	}
	if n := utf8.RuneCountInString(m.Title); n < 3 || n > 120 {
		return errs.Field("title", "3-120 characters")
	}
	if utf8.RuneCountInString(m.Description) > 4000 {
		return errs.Field("description", "at most 4000 characters")
	}
	if !constant.IsCategory(m.Category) {
		return errs.Field("category", "unknown category")
	}
	if m.Campus != "" && !slices.Contains(config.Conf.Campuses, m.Campus) {
		return errs.Field("campus", "unknown campus")
	}
	if len(m.Photos) > maxPhotos {
		return errs.Field("photos", "at most 6 photos")
	}
	for _, p := range m.Photos {
		// only our own uploads: foreign URLs mean tracking and content spoofing
		if !strings.HasPrefix(p, photoUrl) || strings.Contains(p, "..") || len(p) > 200 {
			return errs.Field("photos", "upload photos through /file")
		}
	}

	switch m.Type {
	case constant.ListingTypeRent:
		if m.Price <= 0 || m.Price > maxAmount {
			return errs.Field("price", "price per day is required")
		}
		if m.Deposit < 0 || m.Deposit > maxAmount {
			return errs.Field("deposit", "invalid deposit")
		}
		if !isDate(m.AvailableFrom) {
			return errs.Field("available_from", "YYYY-MM-DD")
		}
		if !isDate(m.AvailableTo) {
			return errs.Field("available_to", "YYYY-MM-DD")
		}
		if m.AvailableTo < m.AvailableFrom {
			return errs.Field("available_to", "must not be before available_from")
		}
		if forCreate && m.AvailableTo < time.Now().Format(dateLayout) {
			return errs.Field("available_to", "must not be in the past")
		}
	case constant.ListingTypeSale:
		if m.Price <= 0 || m.Price > maxAmount {
			return errs.Field("price", "price is required")
		}
	case constant.ListingTypeExchange:
		if n := utf8.RuneCountInString(m.Wanted); n < 3 || n > 200 {
			return errs.Field("wanted", "3-200 characters")
		}
	}
	return nil
}

func isDate(v string) bool {
	_, err := time.Parse(dateLayout, v)
	return err == nil
}
