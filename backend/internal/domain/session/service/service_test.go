package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"qolda/internal/domain/session/model"
	"qolda/internal/errs"
)

func TestTokens(t *testing.T) {
	svc := New("secret", time.Minute, time.Hour)

	tokens, err := svc.CreateTokens(42, "moderator")
	require.NoError(t, err)

	s, err := svc.FromToken(tokens.Access)
	require.NoError(t, err)
	assert.Equal(t, &model.Session{Id: 42, Role: "moderator"}, s)

	id, err := svc.ParseRefresh(tokens.Refresh)
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)

	// tokens of different types are not interchangeable
	_, err = svc.FromToken(tokens.Refresh)
	assert.ErrorIs(t, err, errs.NotAuthorized)
	_, err = svc.ParseRefresh(tokens.Access)
	assert.ErrorIs(t, err, errs.NotAuthorized)

	// foreign signature
	_, err = New("other", time.Minute, time.Hour).FromToken(tokens.Access)
	assert.ErrorIs(t, err, errs.NotAuthorized)

	// expired
	expired, err := New("secret", -time.Second, time.Hour).CreateTokens(42, "student")
	require.NoError(t, err)
	_, err = svc.FromToken(expired.Access)
	assert.ErrorIs(t, err, errs.NotAuthorized)

	_, err = svc.FromToken("garbage")
	assert.ErrorIs(t, err, errs.NotAuthorized)
}

func TestContext(t *testing.T) {
	svc := New("secret", time.Minute, time.Hour)
	ctx := context.Background()

	assert.False(t, svc.CtxIsAuthed(ctx), "no session means anonymous")

	ctx = WithContext(ctx, &model.Session{Id: 7, Role: "student"})
	assert.True(t, svc.CtxIsAuthed(ctx))
	assert.Equal(t, int64(7), svc.CtxGetId(ctx))
	assert.True(t, svc.CtxHasRole(ctx, "student", "admin"))
	assert.False(t, svc.CtxHasRole(ctx, "moderator"))
}
