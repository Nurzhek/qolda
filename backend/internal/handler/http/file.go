package http

import (
	"errors"
	"io"
	"net/http"

	"qolda/internal/errs"
	usecase "qolda/internal/usecase/file"
)

type File struct {
	usecase    *usecase.Usecase
	sessionSvc SessionServiceI
}

func NewFile(uc *usecase.Usecase, sessionSvc SessionServiceI) *File {
	return &File{usecase: uc, sessionSvc: sessionSvc}
}

// Upload: POST /file, multipart field file -> {url}
func (h *File) Upload(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := withSession(r, h.sessionSvc)

	// headroom for multipart headers; the exact limit is checked by the usecase
	r.Body = http.MaxBytesReader(w, r.Body, usecase.MaxSize+64<<10)
	f, _, err := r.FormFile("file")
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			writeErr(w, r, errs.FileTooLarge)
			return
		}
		writeErr(w, r, errs.Field("file", "required"))
		return
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, usecase.MaxSize+1))
	if err != nil {
		writeErr(w, r, errs.FileTooLarge)
		return
	}

	url, err := h.usecase.Upload(ctx, data)
	if err != nil {
		writeErr(w, r, err)
		return
	}
	writeJson(w, http.StatusOK, map[string]string{"url": url})
}
