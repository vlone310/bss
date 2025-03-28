package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vlone310/bss/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := testutil.RandomString(6)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)

	err = CheckPasswordHash(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := testutil.RandomString(6)
	err = CheckPasswordHash(wrongPassword, hashedPassword1)
	require.EqualError(t, errors.Unwrap(err), bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
