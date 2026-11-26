package google

import (
	"context"

	googleModel "qolda/internal/service/google/model"
)

// OAuth is sign-in with Google (OIDC): authorization code + PKCE
type OAuth interface {
	Enabled() bool
	NewVerifier() string
	AuthCodeURL(state, verifier string) string
	Exchange(ctx context.Context, code, verifier string) (*googleModel.Profile, error)
}
