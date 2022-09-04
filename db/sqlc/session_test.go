package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T, username sql.NullString) Session {
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     username.String,
		RefreshToken: util.RandomString(10),
		UserAgent:    util.RandomString(10),
		ClientIp:     util.RandomString(10),
		IsBlocked:    util.RandomBoolean(),
		ExpiresAt:    time.Now().Add(time.Hour * 24),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.WithinDuration(t, arg.ExpiresAt, session.ExpiresAt, time.Second)
	return session
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	session := createRandomSession(t, user.Username)
	testQueries.DeleteSession(context.Background(), session.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestGetSession(t *testing.T) {
	user := createRandomUser(t)
	session1 := createRandomSession(t, user.Username)

	session2, err := testQueries.GetSession(context.Background(), session1.ID)
	require.NoError(t, err)
	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, time.Second)
	require.WithinDuration(t, session1.CreateAt, session2.CreateAt, time.Second)

	testQueries.DeleteSession(context.Background(), session1.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestDeleteSession(t *testing.T) {
	user := createRandomUser(t)
	session := createRandomSession(t, user.Username)

	err := testQueries.DeleteSession(context.Background(), session.ID)
	require.NoError(t, err)

	_, err = testQueries.GetSession(context.Background(), session.ID)
	require.Error(t, err)
	testQueries.DeleteUser(context.Background(), user.ID)
}
