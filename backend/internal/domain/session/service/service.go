package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"qolda/internal/domain/session/model"
	"qolda/internal/errs"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type ctxKey struct{}

type claims struct {
	Role string `json:"role,omitempty"`
	Typ  string `json:"typ"`
	jwt.RegisteredClaims
}

// Service is the only place that knows the token format (HS256 JWT).
type Service struct {
	secret     []byte
	accessTtl  time.Duration
	refreshTtl time.Duration
}

func New(secret string, accessTtl, refreshTtl time.Duration) *Service {
	return &Service{secret: []byte(secret), accessTtl: accessTtl, refreshTtl: refreshTtl}
}

func (s *Service) CreateTokens(usrId int64, role string) (*model.Tokens, error) {
	access, err := s.sign(usrId, role, tokenTypeAccess, s.accessTtl)
	if err != nil {
		return nil, fmt.Errorf("sign access: %w", err)
	}
	refresh, err := s.sign(usrId, "", tokenTypeRefresh, s.refreshTtl)
	if err != nil {
		return nil, fmt.Errorf("sign refresh: %w", err)
	}
	return &model.Tokens{Access: access, Refresh: refresh}, nil
}

// FromToken parses an access token. Invalid/expired -> errs.NotAuthorized.
func (s *Service) FromToken(token string) (*model.Session, error) {
	c, err := s.parse(token, tokenTypeAccess)
	if err != nil {
		return nil, err
	}
	id, _ := strconv.ParseInt(c.Subject, 10, 64)
	return &model.Session{Id: id, Role: c.Role}, nil
}

// ParseRefresh returns the user id from a refresh token.
// ponytail: refresh tokens can't be revoked (stateless); add a token table when "sign out everywhere" is needed
func (s *Service) ParseRefresh(token string) (int64, error) {
	c, err := s.parse(token, tokenTypeRefresh)
	if err != nil {
		return 0, err
	}
	id, _ := strconv.ParseInt(c.Subject, 10, 64)
	return id, nil
}

func (s *Service) sign(usrId int64, role, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Role: role,
		Typ:  typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(usrId, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}).SignedString(s.secret)
}

func (s *Service) parse(token, typ string) (*claims, error) {
	c := &claims{}
	_, err := jwt.ParseWithClaims(token, c, func(*jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil || c.Typ != typ {
		return nil, errs.NotAuthorized
	}
	if id, err := strconv.ParseInt(c.Subject, 10, 64); err != nil || id <= 0 {
		return nil, errs.NotAuthorized
	}
	return c, nil
}

// context

func WithContext(ctx context.Context, s *model.Session) context.Context {
	if s == nil {
		s = &model.Session{}
	}
	return context.WithValue(ctx, ctxKey{}, s)
}

// FromContext never returns nil
func FromContext(ctx context.Context) *model.Session {
	if s, ok := ctx.Value(ctxKey{}).(*model.Session); ok && s != nil {
		return s
	}
	return &model.Session{}
}

func (s *Service) CtxIsAuthed(ctx context.Context) bool {
	return FromContext(ctx).IsAuthed()
}

func (s *Service) CtxHasRole(ctx context.Context, roles ...string) bool {
	return FromContext(ctx).HasRole(roles...)
}

func (s *Service) CtxGetId(ctx context.Context) int64 {
	return FromContext(ctx).Id
}
