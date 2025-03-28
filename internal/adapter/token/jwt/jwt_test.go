package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/vlone310/bss/internal/adapter/token/maker"
	"github.com/vlone310/bss/testutil"
)

func TestJWTMacker(t *testing.T) {
	maker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	username := testutil.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(testutil.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgName(t *testing.T) {
	payload, err := maker.NewPayload(testutil.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, getClaims(payload))
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	jwtMaker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	payload, err = jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
