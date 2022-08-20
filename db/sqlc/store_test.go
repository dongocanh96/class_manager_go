package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserTx(t *testing.T) {
	store := NewStore(testDB)

	user1 := createRandomUser(t)

	n := 2

	username := util.RandomString(10)
	fullname := util.RandomString(10)
	email := util.RandomEmail()
	phoneNumber := util.RandomPhoneNumber()

	errs := make(chan error)
	results := make(chan UpdateUserInfoTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.UpdateUserInfoTx(context.Background(), UpdateUserInfoTxParams{
				ID:          user1.ID,
				Username:    sql.NullString{String: username, Valid: true},
				Fullname:    sql.NullString{String: fullname, Valid: true},
				Email:       sql.NullString{String: email, Valid: true},
				PhoneNumber: sql.NullString{String: phoneNumber, Valid: true},
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		user2 := result.User
		user3, err := store.GetUser(context.Background(), user1.ID)
		require.NoError(t, err)

		require.NotEmpty(t, user2)
		require.Equal(t, user3.ID, user2.ID)
		require.Equal(t, user3.Username, user2.Username)
		require.Equal(t, user3.HashedPassword, user2.HashedPassword)
		require.Equal(t, user3.PasswordChangedAt, user2.PasswordChangedAt)
		require.Equal(t, user3.Fullname, user2.Fullname)
		require.Equal(t, user3.PhoneNumber, user2.PhoneNumber)
	}

}
