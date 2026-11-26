package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/mail"
	"strings"

	"qolda/internal/config"
	"qolda/internal/constant"
	sessionModel "qolda/internal/domain/session/model"
	usrModel "qolda/internal/domain/usr/model"
	"qolda/internal/errs"
)

type Usecase struct {
	usrSvc     UsrServiceI
	sessionSvc SessionServiceI
	google     GoogleI
}

func New(usrSvc UsrServiceI, sessionSvc SessionServiceI, google GoogleI) *Usecase {
	return &Usecase{usrSvc: usrSvc, sessionSvc: sessionSvc, google: google}
}

// DevLogin signs in by e-mail without an IdP (development/pilot). Same domain check as Google.
func (u *Usecase) DevLogin(ctx context.Context, email, name string) (*sessionModel.Tokens, error) {
	if !config.Conf.AuthDevLogin {
		return nil, errs.NoPermission
	}
	return u.login(ctx, email, name)
}

func (u *Usecase) Refresh(ctx context.Context, refreshToken string) (*sessionModel.Tokens, error) {
	usrId, err := u.sessionSvc.ParseRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	usr, found, err := u.usrSvc.Get(ctx, usrId, false)
	if err != nil {
		return nil, fmt.Errorf("usrSvc.Get: %w", err)
	}
	if !found {
		return nil, errs.NotAuthorized
	}
	return u.issue(usr)
}

// GoogleStart returns the Google URL and one-time state/verifier that the transport keeps until the callback
func (u *Usecase) GoogleStart() (authUrl, state, verifier string, err error) {
	if !u.google.Enabled() {
		return "", "", "", errs.NotImplemented
	}
	state = rand.Text()
	verifier = u.google.NewVerifier()
	return u.google.AuthCodeURL(state, verifier), state, verifier, nil
}

func (u *Usecase) GoogleLogin(ctx context.Context, code, verifier string) (*sessionModel.Tokens, error) {
	if !u.google.Enabled() {
		return nil, errs.NotImplemented
	}
	if code == "" || verifier == "" {
		return nil, errs.InvalidRequest
	}
	profile, err := u.google.Exchange(ctx, code, verifier)
	if err != nil {
		return nil, fmt.Errorf("google.Exchange: %w", err)
	}
	if !profile.EmailVerified {
		return nil, errs.EmailNotVerified
	}
	return u.login(ctx, profile.Email, profile.Name)
}

// login checks the domain, finds or creates the user and issues tokens
func (u *Usecase) login(ctx context.Context, email, name string) (*sessionModel.Tokens, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		return nil, errs.Field("email", "invalid e-mail")
	}
	if !DomainAllowed(email, config.Conf.AllowedDomains) {
		return nil, errs.DomainNotAllowed
	}

	usr, found, err := u.usrSvc.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("usrSvc.GetByEmail: %w", err)
	}
	if !found {
		name = strings.TrimSpace(name)
		if name == "" {
			name = email[:strings.IndexByte(email, '@')]
		}
		usr = &usrModel.Main{Email: email, Name: name, Role: constant.RoleStudent}
		usr.Id, err = u.usrSvc.Create(ctx, &usrModel.Edit{Email: &usr.Email, Name: &usr.Name})
		if err != nil {
			return nil, fmt.Errorf("usrSvc.Create: %w", err)
		}
	}
	return u.issue(usr)
}

func (u *Usecase) issue(usr *usrModel.Main) (*sessionModel.Tokens, error) {
	if usr.Blocked {
		return nil, errs.AccountBlocked
	}
	tokens, err := u.sessionSvc.CreateTokens(usr.Id, usr.Role)
	if err != nil {
		return nil, fmt.Errorf("sessionSvc.CreateTokens: %w", err)
	}
	return tokens, nil
}

// DomainAllowed: the address domain equals an allowed one or is its subdomain
// (student.nu.edu.kz for nu.edu.kz), but doesn't merely end with it (evilnu.edu.kz).
func DomainAllowed(email string, allowed []string) bool {
	at := strings.LastIndexByte(email, '@')
	if at < 0 {
		return false
	}
	domain := strings.ToLower(email[at+1:])
	for _, d := range allowed {
		d = strings.ToLower(strings.TrimSpace(d))
		if d != "" && (domain == d || strings.HasSuffix(domain, "."+d)) {
			return true
		}
	}
	return false
}
