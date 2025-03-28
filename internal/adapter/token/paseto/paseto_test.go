package paseto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vlone310/bss/testutil"
)

func TestPasetoMacker(t *testing.T) {
	maker, err := NewPasetoMaker(testutil.RandomString(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(testutil.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(testutil.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
