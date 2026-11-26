package file

import (
	"context"
	"fmt"
	"net/http"

	"qolda/internal/errs"
)

const MaxSize = 8 << 20

// the type is detected from the content, not from the client's file name or header
var allowedTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

type Usecase struct {
	store      StoreI
	sessionSvc SessionServiceI
}

func New(store StoreI, sessionSvc SessionServiceI) *Usecase {
	return &Usecase{store: store, sessionSvc: sessionSvc}
}

func (u *Usecase) Upload(ctx context.Context, data []byte) (string, error) {
	if !u.sessionSvc.CtxIsAuthed(ctx) {
		return "", errs.NotAuthorized
	}
	if len(data) == 0 {
		return "", errs.Field("file", "required")
	}
	if len(data) > MaxSize {
		return "", errs.FileTooLarge
	}
	ext, ok := allowedTypes[http.DetectContentType(data)]
	if !ok {
		return "", errs.FileTypeNotAllowed
	}
	url, err := u.store.Save(ctx, ext, data)
	if err != nil {
		return "", fmt.Errorf("store.Save: %w", err)
	}
	return url, nil
}
