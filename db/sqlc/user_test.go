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
		Username:       util.RandomString(6),
		HashedPassword: hashPassword,
		Fullname:       util.RandomString(6),
		Email:          util.RandomEmail(),
		PhoneNumber:    util.RandomPhoneNumber(),
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
		Username:       util.RandomString(6),
		HashedPassword: hashPassword,
		Fullname:       util.RandomString(6),
		Email:          util.RandomEmail(),
		PhoneNumber:    util.RandomPhoneNumber(),
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
		Username:       util.RandomString(6),
		HashedPassword: hashPassword,
		Fullname:       util.RandomString(6),
		Email:          util.RandomEmail(),
		PhoneNumber:    util.RandomPhoneNumber(),
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

func TestUpdateUsername(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUsernameParams{
		ID:       user1.ID,
		Username: "ngoc anh",
	}

	user2, err := testQueries.UpdateUsername(context.Background(), arg)
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

func TestUpdateHashedPassword(t *testing.T) {
	user1 := createRandomUser(t)
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := UpdateHashedPasswordParams{
		ID:                user1.ID,
		HashedPassword:    newHashedPassword,
		PasswordChangedAt: time.Now(),
	}
	user2, err := testQueries.UpdateHashedPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, arg.PasswordChangedAt, user2.PasswordChangedAt, time.Second)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestUpdateFullName(t *testing.T) {
	user1 := createRandomUser(t)
	fullname := util.RandomString(6)

	arg := UpdateFullnameParams{
		ID:       user1.ID,
		Fullname: fullname,
	}

	user2, err := testQueries.UpdateFullname(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Fullname, user2.Fullname)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestUpdateEmail(t *testing.T) {
	user1 := createRandomUser(t)
	email := util.RandomEmail()

	arg := UpdateEmailParams{
		ID:    user1.ID,
		Email: email,
	}

	user2, err := testQueries.UpdateEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Email, user2.Email)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.PhoneNumber, user2.PhoneNumber)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}

func TestUpdatePhoneNumber(t *testing.T) {
	user1 := createRandomUser(t)
	phoneNumber := util.RandomPhoneNumber()

	arg := UpdatePhoneNumberParams{
		ID:          user1.ID,
		PhoneNumber: phoneNumber,
	}

	user2, err := testQueries.UpdatePhoneNumber(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.PhoneNumber, user2.PhoneNumber)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.IsTeacher, user2.IsTeacher)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	testQueries.DeleteUser(context.Background(), user1.ID)
}
