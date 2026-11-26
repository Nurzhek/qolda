package file

import "context"

type StoreI interface {
	Save(ctx context.Context, ext string, data []byte) (string, error)
}

type SessionServiceI interface {
	CtxIsAuthed(ctx context.Context) bool
}
