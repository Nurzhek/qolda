package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	googleModel "qolda/internal/service/google/model"
)

const userInfoUrl = "https://openidconnect.googleapis.com/v1/userinfo"

type Service struct {
	conf *oauth2.Config
	// hd hints Google to show accounts of this domain; the domain is still checked on our side
	hostedDomain string
}

func New(clientId, clientSecret, redirectUrl, hostedDomain string) *Service {
	return &Service{
		conf: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "email", "profile"},
		},
		hostedDomain: hostedDomain,
	}
}

func (s *Service) Enabled() bool {
	return s.conf.ClientID != "" && s.conf.ClientSecret != ""
}

func (s *Service) NewVerifier() string {
	return oauth2.GenerateVerifier()
}

func (s *Service) AuthCodeURL(state, verifier string) string {
	opts := []oauth2.AuthCodeOption{
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("prompt", "select_account"),
	}
	if s.hostedDomain != "" {
		opts = append(opts, oauth2.SetAuthURLParam("hd", s.hostedDomain))
	}
	return s.conf.AuthCodeURL(state, opts...)
}

// Exchange trades the code for a token and reads the profile from userinfo. The token comes straight from Google
// over TLS with our client secret, so a separate ID token signature check is not needed.
func (s *Service) Exchange(ctx context.Context, code, verifier string) (*googleModel.Profile, error) {
	tok, err := s.conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("conf.Exchange: %w", err)
	}

	rep, err := s.conf.Client(ctx, tok).Get(userInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("userinfo request: %w", err)
	}
	defer rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo: bad status code %d", rep.StatusCode)
	}

	var info struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}
	if err = json.NewDecoder(rep.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("userinfo decode: %w", err)
	}

	return &googleModel.Profile{Email: info.Email, EmailVerified: info.EmailVerified, Name: info.Name}, nil
}
