package http

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"qolda/internal/errs"
	usecase "qolda/internal/usecase/auth"
)

const (
	oauthCookie  = "qolda_oauth"
	oauthPath    = "/api/auth/google"
	callbackPage = "/auth/callback"
)

// Auth is Google OAuth (authorization code + PKCE) via browser redirects
type Auth struct {
	usecase *usecase.Usecase
	secure  bool
}

func NewAuth(uc *usecase.Usecase, publicUrl string) *Auth {
	return &Auth{usecase: uc, secure: strings.HasPrefix(publicUrl, "https://")}
}

// GoogleStart — GET /auth/google?redirect=/path
func (h *Auth) GoogleStart(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	authUrl, state, verifier, err := h.usecase.GoogleStart()
	if err != nil {
		h.fail(w, r, err)
		return
	}
	v := url.Values{"s": {state}, "v": {verifier}, "r": {safeRedirect(r.URL.Query().Get("redirect"))}}
	h.setCookie(w, v.Encode(), 600)
	http.Redirect(w, r, authUrl, http.StatusFound)
}

// GoogleCallback — GET /auth/google/callback?code&state
func (h *Auth) GoogleCallback(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	c, err := r.Cookie(oauthCookie)
	h.setCookie(w, "", -1) // single use
	if err != nil {
		h.fail(w, r, errs.InvalidRequest)
		return
	}
	saved, err := url.ParseQuery(c.Value)
	q := r.URL.Query()
	if err != nil || q.Get("error") != "" || saved.Get("s") == "" || q.Get("state") != saved.Get("s") {
		h.fail(w, r, errs.InvalidRequest)
		return
	}

	tokens, err := h.usecase.GoogleLogin(r.Context(), q.Get("code"), saved.Get("v"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	// tokens go in the fragment: they stay out of server logs and Referer
	frag := url.Values{
		"access_token":  {tokens.Access},
		"refresh_token": {tokens.Refresh},
		"redirect":      {safeRedirect(saved.Get("r"))},
	}
	http.Redirect(w, r, callbackPage+"#"+frag.Encode(), http.StatusFound)
}

func (h *Auth) fail(w http.ResponseWriter, r *http.Request, err error) {
	code := errs.Code(err)
	switch errs.Err(code) {
	case errs.DomainNotAllowed, errs.EmailNotVerified, errs.AccountBlocked:
		slog.Info("google login rejected", "code", code)
	default:
		slog.Error("google login failed", "error", err)
		code = "oauth_failed"
	}
	http.Redirect(w, r, callbackPage+"#error="+code, http.StatusFound)
}

func (h *Auth) setCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     oauthCookie,
		Value:    value,
		Path:     oauthPath,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// safeRedirect allows only internal SPA paths, no open redirect to //evil.com
func safeRedirect(v string) string {
	if !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") || strings.HasPrefix(v, "/\\") {
		return "/"
	}
	return v
}
