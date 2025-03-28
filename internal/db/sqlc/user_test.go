package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vlone310/bss/testutil"
	"github.com/vlone310/bss/util"
)

func createRandomUser(t *testing.T) User {
	t.Helper()

	hashedPassword, err := util.HashPassword(testutil.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       testutil.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       testutil.RandomOwner(),
		Email:          testutil.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt.Time, user2.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
