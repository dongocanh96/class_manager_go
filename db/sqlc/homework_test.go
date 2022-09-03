package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func createRandomHomework(t *testing.T, teacherID int64, subject string) Homework {
	arg := CreateHomeworkParams{
		TeacherID: teacherID,
		Subject:   subject,
		Title:     util.RandomString(10),
		FileName:  util.RandomString(10),
		SavedPath: util.RandomString(10),
	}

	homework, err := testQueries.CreateHomework(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homework)

	require.NotZero(t, homework.ID)
	require.Equal(t, arg.TeacherID, homework.TeacherID)
	require.Equal(t, arg.Subject, homework.Subject)
	require.Equal(t, arg.Title, homework.Title)
	require.Equal(t, arg.FileName, homework.FileName)
	require.Equal(t, arg.SavedPath, homework.SavedPath)
	require.False(t, homework.IsClosed)

	require.NotZero(t, homework.CreatedAt)
	require.True(t, homework.UpdatedAt.IsZero())
	require.True(t, homework.ClosedAt.IsZero())

	return homework
}

func TestCreateHomework(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestGetHomework(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()
	homework1 := createRandomHomework(t, teacher.ID, subject)

	homework2, err := testQueries.GetHomework(context.Background(), homework1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, homework2)

	require.Equal(t, homework1.ID, homework2.ID)
	require.Equal(t, homework1.TeacherID, homework2.TeacherID)
	require.Equal(t, homework1.Subject, homework2.Subject)
	require.Equal(t, homework1.Title, homework2.Title)
	require.Equal(t, homework1.FileName, homework2.FileName)
	require.Equal(t, homework1.SavedPath, homework2.SavedPath)
	require.Equal(t, homework1.IsClosed, homework2.IsClosed)
	require.WithinDuration(t, homework1.CreatedAt, homework2.CreatedAt, time.Second)
	require.WithinDuration(t, homework1.UpdatedAt, homework2.UpdatedAt, time.Second)
	require.WithinDuration(t, homework1.ClosedAt, homework2.ClosedAt, time.Second)

	testQueries.DeleteHomework(context.Background(), homework2.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestDeleteHomework(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	err := testQueries.DeleteHomework(context.Background(), homework.ID)
	require.NoError(t, err)

	_, err = testQueries.GetHomework(context.Background(), homework.ID)
	require.Error(t, sql.ErrNoRows, err)
}

func TestListHomeWorskByTeacher(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()

	for i := 0; i < 5; i++ {
		createRandomHomework(t, teacher.ID, subject)
	}

	arg := ListHomeworksByTeacherParams{
		TeacherID: teacher.ID,
		Limit:     5,
		Offset:    0,
	}

	homeworks, err := testQueries.ListHomeworksByTeacher(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homeworks)

	for _, homework := range homeworks {
		require.NotEmpty(t, homework)
	}

	for i := range homeworks {
		testQueries.DeleteHomework(context.Background(), homeworks[i].ID)
	}

	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestListHomeworksBySubject(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()

	for i := 0; i < 5; i++ {
		createRandomHomework(t, teacher.ID, subject)
	}

	arg := ListHomeworksBySubjectParams{
		Subject: subject,
		Limit:   5,
		Offset:  0,
	}

	homeworks, err := testQueries.ListHomeworksBySubject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homeworks)

	for _, homework := range homeworks {
		require.NotEmpty(t, homework)
	}

	for i := range homeworks {
		testQueries.DeleteHomework(context.Background(), homeworks[i].ID)
	}

	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestListHomeworks(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()

	for i := 0; i < 5; i++ {
		createRandomHomework(t, teacher.ID, subject)
	}

	arg := ListHomeworksParams{
		Limit:  5,
		Offset: 0,
	}

	homeworks, err := testQueries.ListHomeworks(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homeworks)

	for _, homework := range homeworks {
		require.NotEmpty(t, homework)
	}

	for i := range homeworks {
		testQueries.DeleteHomework(context.Background(), homeworks[i].ID)
	}

	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestUpdateHomework(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()
	homework1 := createRandomHomework(t, teacher.ID, subject)

	arg := UpdateHomeworkParams{
		ID:        homework1.ID,
		FileName:  util.RandomString(10),
		SavedPath: util.RandomString(10),
		UpdatedAt: time.Now(),
	}

	homework2, err := testQueries.UpdateHomework(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homework2)

	require.Equal(t, homework1.ID, homework2.ID)
	require.Equal(t, homework1.TeacherID, homework2.TeacherID)
	require.Equal(t, homework1.Subject, homework2.Subject)
	require.WithinDuration(t, homework1.CreatedAt, homework2.CreatedAt, time.Second)
	require.WithinDuration(t, homework1.ClosedAt, homework2.ClosedAt, time.Second)

	require.Equal(t, homework1.Title, homework2.Title)
	require.Equal(t, arg.FileName, homework2.FileName)
	require.Equal(t, arg.SavedPath, homework2.SavedPath)
	require.WithinDuration(t, arg.UpdatedAt, homework2.UpdatedAt, time.Second)

	testQueries.DeleteHomework(context.Background(), homework1.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
}

func TestCloseHomework(t *testing.T) {
	teacher := createRandomTeacher(t)
	subject := util.RandomSubject()
	homework1 := createRandomHomework(t, teacher.ID, subject)

	arg := CloseHomeworkParams{
		ID:       homework1.ID,
		IsClosed: true,
		ClosedAt: time.Now(),
	}

	homework2, err := testQueries.CloseHomework(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, homework2)

	require.Equal(t, homework1.ID, homework2.ID)
	require.Equal(t, homework1.TeacherID, homework2.TeacherID)
	require.Equal(t, homework1.Subject, homework2.Subject)
	require.Equal(t, homework1.Title, homework2.Title)
	require.Equal(t, homework1.FileName, homework2.FileName)
	require.Equal(t, homework1.SavedPath, homework2.SavedPath)
	require.WithinDuration(t, homework1.CreatedAt, homework2.CreatedAt, time.Second)

	require.True(t, homework2.IsClosed)
	require.WithinDuration(t, arg.ClosedAt, homework2.ClosedAt, time.Second)

	testQueries.DeleteHomework(context.Background(), homework2.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
}
