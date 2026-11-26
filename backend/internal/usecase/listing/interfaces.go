package listing

import (
	"context"

	"qolda/internal/domain/listing/model"
	usrModel "qolda/internal/domain/usr/model"
)

type ServiceI interface {
	List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error)
	Get(ctx context.Context, id int64, errNE bool) (*model.Main, bool, error)
	Create(ctx context.Context, obj *model.Edit) (int64, error)
	Update(ctx context.Context, id int64, obj *model.Edit) error
}

type UsrServiceI interface {
	List(ctx context.Context, pars *usrModel.ListReq) ([]*usrModel.Main, int64, error)
}

type SessionServiceI interface {
	CtxIsAuthed(ctx context.Context) bool
	CtxGetId(ctx context.Context) int64
	CtxHasRole(ctx context.Context, roles ...string) bool
}
