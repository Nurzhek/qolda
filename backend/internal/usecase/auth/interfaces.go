package auth

import (
	"context"

	sessionModel "qolda/internal/domain/session/model"
	usrModel "qolda/internal/domain/usr/model"
	googleModel "qolda/internal/service/google/model"
)

type UsrServiceI interface {
	Get(ctx context.Context, id int64, errNE bool) (*usrModel.Main, bool, error)
	GetByEmail(ctx context.Context, email string) (*usrModel.Main, bool, error)
	Create(ctx context.Context, obj *usrModel.Edit) (int64, error)
}

type SessionServiceI interface {
	CreateTokens(usrId int64, role string) (*sessionModel.Tokens, error)
	ParseRefresh(token string) (int64, error)
}

type GoogleI interface {
	Enabled() bool
	NewVerifier() string
	AuthCodeURL(state, verifier string) string
	Exchange(ctx context.Context, code, verifier string) (*googleModel.Profile, error)
}
