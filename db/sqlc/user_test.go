package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T) User { //remember to return Account model
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	// we can create a obj instance manually, e.g. tom, 1000, USD, or use util.Random that we created
	arg := CreateUserParams{
		UsernameOfUser: util.RandomUsername(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomUsername(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err) // this is to check if err == nil, fail if err != nil
	require.NotEmpty(t, user)

	require.Equal(t, arg.UsernameOfUser, user.UsernameOfUser)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UsernameOfUser)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UsernameOfUser, user2.UsernameOfUser)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
