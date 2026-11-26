package usr

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"qolda/internal/config"
	"qolda/internal/domain/usr/model"
	"qolda/internal/errs"
)

var telegramRe = regexp.MustCompile(`^[A-Za-z0-9_]{5,32}$`)

type Usecase struct {
	svc        ServiceI
	sessionSvc SessionServiceI
}

func New(svc ServiceI, sessionSvc SessionServiceI) *Usecase {
	return &Usecase{svc: svc, sessionSvc: sessionSvc}
}

func (u *Usecase) GetMe(ctx context.Context) (*model.Main, error) {
	if !u.sessionSvc.CtxIsAuthed(ctx) {
		return nil, errs.NotAuthorized
	}
	result, _, err := u.svc.Get(ctx, u.sessionSvc.CtxGetId(ctx), true)
	if err != nil {
		return nil, fmt.Errorf("svc.Get: %w", err)
	}
	return result, nil
}

// Get returns a public profile: only the user themselves sees the e-mail
func (u *Usecase) Get(ctx context.Context, id int64) (*model.Main, error) {
	if id == 0 {
		return nil, errs.IdRequired
	}
	result, _, err := u.svc.Get(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("svc.Get: %w", err)
	}
	if result.Id != u.sessionSvc.CtxGetId(ctx) {
		result.Email = ""
	}
	return result, nil
}

func (u *Usecase) UpdateMe(ctx context.Context, obj *model.Edit) error {
	if !u.sessionSvc.CtxIsAuthed(ctx) {
		return errs.NotAuthorized
	}
	// only profile fields can be changed
	obj = &model.Edit{Name: obj.Name, Campus: obj.Campus, Telegram: obj.Telegram}
	if err := u.validateEdit(obj); err != nil {
		return err
	}
	if err := u.svc.Update(ctx, u.sessionSvc.CtxGetId(ctx), obj); err != nil {
		return fmt.Errorf("svc.Update: %w", err)
	}
	return nil
}

func (u *Usecase) validateEdit(obj *model.Edit) error {
	if obj.Name != nil {
		obj.Name = new(strings.TrimSpace(*obj.Name))
		if n := utf8.RuneCountInString(*obj.Name); n < 2 || n > 60 {
			return errs.Field("name", "2-60 characters")
		}
	}
	if obj.Campus != nil && *obj.Campus != "" && !slices.Contains(config.Conf.Campuses, *obj.Campus) {
		return errs.Field("campus", "unknown campus")
	}
	if obj.Telegram != nil {
		v := strings.TrimSpace(*obj.Telegram)
		v = strings.TrimPrefix(v, "https://t.me/")
		v = strings.TrimPrefix(v, "@")
		if v != "" && !telegramRe.MatchString(v) {
			return errs.Field("telegram", "5-32 letters, digits or _")
		}
		obj.Telegram = &v
	}
	return nil
}
