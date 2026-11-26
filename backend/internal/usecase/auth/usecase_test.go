package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"qolda/internal/config"
	sessionServiceP "qolda/internal/domain/session/service"
	usrModel "qolda/internal/domain/usr/model"
	"qolda/internal/errs"
	googleModel "qolda/internal/service/google/model"
)

// fakeUsr — in-memory UsrServiceI
type fakeUsr struct {
	byId map[int64]*usrModel.Main
}

func (f *fakeUsr) Get(_ context.Context, id int64, _ bool) (*usrModel.Main, bool, error) {
	u, ok := f.byId[id]
	return u, ok, nil
}

func (f *fakeUsr) GetByEmail(_ context.Context, email string) (*usrModel.Main, bool, error) {
	for _, u := range f.byId {
		if u.Email == email {
			return u, true, nil
		}
	}
	return nil, false, nil
}

func (f *fakeUsr) Create(_ context.Context, obj *usrModel.Edit) (int64, error) {
	id := int64(len(f.byId) + 1)
	f.byId[id] = &usrModel.Main{Id: id, Email: *obj.Email, Name: *obj.Name, Role: "student"}
	return id, nil
}

type fakeGoogle struct {
	profile *googleModel.Profile
}

func (f *fakeGoogle) Enabled() bool       { return true }
func (f *fakeGoogle) NewVerifier() string { return "verifier" }
func (f *fakeGoogle) AuthCodeURL(state, _ string) string {
	return "https://accounts.google.com/?state=" + state
}
func (f *fakeGoogle) Exchange(context.Context, string, string) (*googleModel.Profile, error) {
	return f.profile, nil
}

func setup(t *testing.T) (*Usecase, *fakeUsr, *fakeGoogle, *sessionServiceP.Service) {
	t.Helper()
	prevDomains, prevDev := config.Conf.AllowedDomains, config.Conf.AuthDevLogin
	t.Cleanup(func() { config.Conf.AllowedDomains, config.Conf.AuthDevLogin = prevDomains, prevDev })
	config.Conf.AllowedDomains = []string{"nu.edu.kz"}
	config.Conf.AuthDevLogin = true

	usr := &fakeUsr{byId: map[int64]*usrModel.Main{}}
	google := &fakeGoogle{}
	session := sessionServiceP.New("secret", time.Minute, time.Hour)
	return New(usr, session, google), usr, google, session
}

func TestDomainAllowed(t *testing.T) {
	allowed := []string{"nu.edu.kz"}
	assert.True(t, DomainAllowed("a@nu.edu.kz", allowed))
	assert.True(t, DomainAllowed("a@student.nu.edu.kz", allowed))
	assert.True(t, DomainAllowed("a@NU.EDU.KZ", allowed))
	assert.False(t, DomainAllowed("a@gmail.com", allowed))
	assert.False(t, DomainAllowed("a@evilnu.edu.kz", allowed))
	assert.False(t, DomainAllowed("a@nu.edu.kz.evil.com", allowed))
	assert.False(t, DomainAllowed("nu.edu.kz", allowed))
}

func TestDevLogin(t *testing.T) {
	uc, usr, _, session := setup(t)
	ctx := context.Background()

	// first sign-in creates a student; the name comes from the local part
	tokens, err := uc.DevLogin(ctx, "  Aidana.S@NU.edu.kz ", "")
	require.NoError(t, err)
	s, err := session.FromToken(tokens.Access)
	require.NoError(t, err)
	assert.Equal(t, "student", s.Role)
	require.Len(t, usr.byId, 1)
	assert.Equal(t, "aidana.s@nu.edu.kz", usr.byId[s.Id].Email)
	assert.Equal(t, "aidana.s", usr.byId[s.Id].Name)

	// signing in again returns the same user
	_, err = uc.DevLogin(ctx, "aidana.s@nu.edu.kz", "")
	require.NoError(t, err)
	assert.Len(t, usr.byId, 1)

	// foreign domain and garbage
	_, err = uc.DevLogin(ctx, "someone@gmail.com", "")
	assert.ErrorIs(t, err, errs.DomainNotAllowed)
	_, err = uc.DevLogin(ctx, "Name <x@nu.edu.kz>", "")
	assert.Equal(t, "invalid_request", errs.Code(err))
	assert.Len(t, usr.byId, 1, "a rejection does not create a user")

	// blocked
	usr.byId[s.Id].Blocked = true
	_, err = uc.DevLogin(ctx, "aidana.s@nu.edu.kz", "")
	assert.ErrorIs(t, err, errs.AccountBlocked)

	// dev sign-in disabled
	config.Conf.AuthDevLogin = false
	_, err = uc.DevLogin(ctx, "other@nu.edu.kz", "")
	assert.ErrorIs(t, err, errs.NoPermission)
}

func TestRefresh(t *testing.T) {
	uc, usr, _, session := setup(t)
	ctx := context.Background()

	tokens, err := uc.DevLogin(ctx, "a@nu.edu.kz", "A")
	require.NoError(t, err)

	// the role in the new access token comes from the DB, not the old token
	usr.byId[1].Role = "moderator"
	next, err := uc.Refresh(ctx, tokens.Refresh)
	require.NoError(t, err)
	s, err := session.FromToken(next.Access)
	require.NoError(t, err)
	assert.Equal(t, "moderator", s.Role)

	_, err = uc.Refresh(ctx, tokens.Access)
	assert.ErrorIs(t, err, errs.NotAuthorized, "an access token does not refresh the session")

	usr.byId[1].Blocked = true
	_, err = uc.Refresh(ctx, tokens.Refresh)
	assert.ErrorIs(t, err, errs.AccountBlocked)
}

func TestGoogleLogin(t *testing.T) {
	uc, usr, google, _ := setup(t)
	ctx := context.Background()

	authUrl, state, verifier, err := uc.GoogleStart()
	require.NoError(t, err)
	assert.NotEmpty(t, state)
	assert.Equal(t, "verifier", verifier)
	assert.Contains(t, authUrl, state)

	google.profile = &googleModel.Profile{Email: "b@nu.edu.kz", EmailVerified: false, Name: "B"}
	_, err = uc.GoogleLogin(ctx, "code", verifier)
	assert.ErrorIs(t, err, errs.EmailNotVerified)

	// the IdP lets any Google account in: we check the domain ourselves
	google.profile = &googleModel.Profile{Email: "b@gmail.com", EmailVerified: true, Name: "B"}
	_, err = uc.GoogleLogin(ctx, "code", verifier)
	assert.ErrorIs(t, err, errs.DomainNotAllowed)

	google.profile = &googleModel.Profile{Email: "b@nu.edu.kz", EmailVerified: true, Name: "Bekzat"}
	_, err = uc.GoogleLogin(ctx, "code", verifier)
	require.NoError(t, err)
	assert.Equal(t, "Bekzat", usr.byId[1].Name)

	_, err = uc.GoogleLogin(ctx, "", verifier)
	assert.ErrorIs(t, err, errs.InvalidRequest)
}
