package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	sessionModel "qolda/internal/domain/session/model"
	sessionServiceP "qolda/internal/domain/session/service"
	"qolda/internal/errs"
)

type SessionServiceI interface {
	FromToken(token string) (*sessionModel.Session, error)
}

// withSession works like the gRPC interceptor: an invalid token means anonymous, the usecase decides
func withSession(r *http.Request, svc SessionServiceI) context.Context {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return sessionServiceP.WithContext(r.Context(), nil)
	}
	s, _ := svc.FromToken(token)
	return sessionServiceP.WithContext(r.Context(), s)
}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeErr uses the same {code, message, fields} format as grpc-gateway
func writeErr(w http.ResponseWriter, r *http.Request, err error) {
	code := errs.Code(err)
	level := slog.LevelInfo
	if errs.IsFailure(code) {
		level = slog.LevelError
	}
	slog.Log(r.Context(), level, "HTTP handler error", "error", err.Error(), "path", r.URL.Path)

	rep := map[string]any{"code": code, "message": code, "fields": map[string]string{}}
	if full, ok := errs.AsErrFull(err); ok {
		rep["message"], rep["fields"] = full.Desc, full.Fields
	}
	writeJson(w, errs.HttpStatus(code), rep)
}
