package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func createRandomMessage(t *testing.T, user1, user2 User) Message {
	arg := CreateMessageParams{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Content:    util.RandomString(100),
	}

	message, err := testQueries.CreateMessage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message)

	require.NotZero(t, message.ID)
	require.Equal(t, arg.FromUserID, message.FromUserID)
	require.Equal(t, arg.ToUserID, message.ToUserID)
	require.Equal(t, arg.Content, message.Content)
	require.False(t, message.IsRead)
	require.NotZero(t, message.CreatedAt)
	require.True(t, message.ReadAt.IsZero())

	return message
}

func TestCreateMessage(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	message := createRandomMessage(t, user1, user2)
	testQueries.DeleteMessage(context.Background(), message.ID)
	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestGetMessage(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	message1 := createRandomMessage(t, user1, user2)

	message2, err := testQueries.GetMessage(context.Background(), message1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, message2)

	require.Equal(t, message1.ID, message2.ID)
	require.Equal(t, message1.FromUserID, message2.FromUserID)
	require.Equal(t, message1.ToUserID, message2.ToUserID)
	require.Equal(t, message1.Content, message2.Content)
	require.Equal(t, message1.IsRead, message2.IsRead)
	require.WithinDuration(t, message1.CreatedAt, message2.CreatedAt, time.Second)
	require.WithinDuration(t, message1.ReadAt, message2.ReadAt, time.Second)

	testQueries.DeleteMessage(context.Background(), message1.ID)
	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestDeleteMessage(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	message := createRandomMessage(t, user1, user2)

	err := testQueries.DeleteMessage(context.Background(), message.ID)
	require.NoError(t, err)

	_, err = testQueries.GetMessage(context.Background(), message.ID)
	require.Error(t, sql.ErrNoRows, err)
}

func TestListMessages(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomMessage(t, user1, user2)
	}

	arg := ListMessagesParams{
		FromUserID: user1.ID,
		ToUserID:   user2.ID,
		Limit:      5,
		Offset:     0,
	}

	messages, err := testQueries.ListMessages(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	for _, message := range messages {
		require.NotEmpty(t, message)
	}

	for i := range messages {
		testQueries.DeleteMessage(context.Background(), messages[i].ID)
	}

	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestListMessagesFromUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomMessage(t, user1, user2)
	}

	arg := ListMessagesFromUserParams{
		FromUserID: user1.ID,
		Limit:      5,
		Offset:     0,
	}

	messages, err := testQueries.ListMessagesFromUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	for _, message := range messages {
		require.NotEmpty(t, message)
	}

	for i := range messages {
		testQueries.DeleteMessage(context.Background(), messages[i].ID)
	}

	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestListMessagesToUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomMessage(t, user1, user2)
	}

	arg := ListMessagesToUserParams{
		ToUserID: user2.ID,
		Limit:    5,
		Offset:   0,
	}

	messages, err := testQueries.ListMessagesToUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	for _, message := range messages {
		require.NotEmpty(t, message)
	}

	for i := range messages {
		testQueries.DeleteMessage(context.Background(), messages[i].ID)
	}

	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestUpdateMessage(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	message1 := createRandomMessage(t, user1, user2)

	arg := UpdateMessageParams{
		ID:      message1.ID,
		Content: util.RandomString(100),
	}

	message2, err := testQueries.UpdateMessage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message2)

	require.Equal(t, message1.ID, message2.ID)
	require.Equal(t, message1.FromUserID, message2.FromUserID)
	require.Equal(t, message1.ToUserID, message2.ToUserID)
	require.Equal(t, arg.Content, message2.Content)
	require.False(t, message2.IsRead)
	require.WithinDuration(t, message1.CreatedAt, message2.CreatedAt, time.Second)
	require.True(t, message1.ReadAt.IsZero())

	testQueries.DeleteMessage(context.Background(), message1.ID)
	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}

func TestUpdateMessageState(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	message1 := createRandomMessage(t, user1, user2)

	arg := UpdateMessageStateParams{
		ID:     message1.ID,
		IsRead: true,
		ReadAt: time.Now(),
	}

	message2, err := testQueries.UpdateMessageState(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message2)

	require.Equal(t, message1.ID, message2.ID)
	require.Equal(t, message1.FromUserID, message2.FromUserID)
	require.Equal(t, message1.ToUserID, message2.ToUserID)
	require.Equal(t, message1.Content, message2.Content)
	require.True(t, message2.IsRead)
	require.WithinDuration(t, message1.CreatedAt, message2.CreatedAt, time.Second)
	require.WithinDuration(t, arg.ReadAt, message2.ReadAt, time.Second)

	testQueries.DeleteMessage(context.Background(), message1.ID)
	testQueries.DeleteUser(context.Background(), user1.ID)
	testQueries.DeleteUser(context.Background(), user2.ID)
}
