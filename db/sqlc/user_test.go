package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      util.RandomBoolean(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.IsTeacher, user.IsTeacher)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func createRandomTeacher(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      true,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.IsTeacher, user.IsTeacher)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func createRandomStudent(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      false,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.IsTeacher, user.IsTeacher)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestGetByUserName(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	_, err = testQueries.GetUser(context.Background(), user.ID)
	require.Error(t, sql.ErrNoRows, err)
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

	for i := range users {
		testQueries.DeleteUser(context.Background(), users[i].ID)
	}
}

func TestListTeacher(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTeacher(t)
	}

	arg := ListTeachersOrStudentsParams{
		IsTeacher: true,
		Limit:     10,
		Offset:    0,
	}

	users, err := testQueries.ListTeachersOrStudents(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

	for i := range users {
		testQueries.DeleteUser(context.Background(), users[i].ID)
	}
}

func TestListStudent(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomStudent(t)
	}

	arg := ListTeachersOrStudentsParams{
		IsTeacher: false,
		Limit:     10,
		Offset:    0,
	}

	users, err := testQueries.ListTeachersOrStudents(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

	for i := range users {
		testQueries.DeleteUser(context.Background(), users[i].ID)
	}
}

func TestUpdateUserInfo(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserInfoParams{
		ID:          user1.ID,
		Username:    sql.NullString{String: "ngoc anh", Valid: true},
		Fullname:    sql.NullString{String: "Do Ngoc Anh", Valid: true},
		Email:       sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber: sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
	}

	user2, err := testQueries.UpdateUserInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, arg.Fullname, user2.Fullname)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.PhoneNumber, user2.PhoneNumber)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestUpdateUserName(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserInfoParams{
		ID:       user1.ID,
		Username: sql.NullString{String: "ngoc anh", Valid: true},
	}

	user2, err := testQueries.UpdateUserInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Username, user2.Username)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	arg := UpdateUserPasswordParams{
		ID:                user1.ID,
		HashedPassword:    hashPassword,
		PasswordChangedAt: time.Now(),
	}

	user2, err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)

	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, arg.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
