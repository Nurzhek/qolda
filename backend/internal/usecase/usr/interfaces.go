package usr

import (
	"context"

	"qolda/internal/domain/usr/model"
)

type ServiceI interface {
	Get(ctx context.Context, id int64, errNE bool) (*model.Main, bool, error)
	Update(ctx context.Context, id int64, obj *model.Edit) error
}

type SessionServiceI interface {
	CtxIsAuthed(ctx context.Context) bool
	CtxGetId(ctx context.Context) int64
}
