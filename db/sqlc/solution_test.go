package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dongocanh96/class_manager_go/util"
	"github.com/stretchr/testify/require"
)

func createRandomSolution(t *testing.T, userID, homeworkID int64) Solution {
	arg := CreateSolutionParams{
		ProblemID: homeworkID,
		UserID:    userID,
		FileName:  util.RandomString(6),
		SavedPath: util.RandomString(6),
	}

	solution, err := testQueries.CreateSolution(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solution)

	require.NotZero(t, solution.ID)
	require.Equal(t, arg.ProblemID, solution.ProblemID)
	require.Equal(t, arg.UserID, solution.UserID)
	require.Equal(t, arg.FileName, solution.FileName)
	require.Equal(t, arg.SavedPath, solution.SavedPath)
	require.NotZero(t, solution.SubmitedAt)
	require.True(t, solution.UpdatedAt.IsZero())

	return solution
}

func TestCreateSolution(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	solution := createRandomSolution(t, user.ID, homework.ID)

	testQueries.DeleteSolution(context.Background(), solution.ID)
	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestGetSolutionByID(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	solution1 := createRandomSolution(t, user.ID, homework.ID)

	solution2, err := testQueries.GetSolutionByID(context.Background(), solution1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, solution2)

	require.Equal(t, solution1.ID, solution2.ID)
	require.Equal(t, solution1.ProblemID, solution2.ProblemID)
	require.Equal(t, solution1.UserID, solution2.UserID)
	require.Equal(t, solution1.FileName, solution2.FileName)
	require.Equal(t, solution1.SavedPath, solution2.SavedPath)
	require.WithinDuration(t, solution1.SubmitedAt, solution2.SubmitedAt, time.Second)
	require.WithinDuration(t, solution1.UpdatedAt, solution2.UpdatedAt, time.Second)

	testQueries.DeleteSolution(context.Background(), solution1.ID)
	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestDeleteSolution(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	solution := createRandomSolution(t, user.ID, homework.ID)

	err := testQueries.DeleteSolution(context.Background(), solution.ID)
	require.NoError(t, err)

	_, err = testQueries.GetSolutionByID(context.Background(), solution.ID)
	require.Error(t, sql.ErrNoRows, err)

	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestGetSolutionByProblemAndUser(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	solution1 := createRandomSolution(t, user.ID, homework.ID)

	arg := GetSolutionByProblemAndUserParams{
		ProblemID: homework.ID,
		UserID:    user.ID,
	}

	solution2, err := testQueries.GetSolutionByProblemAndUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solution2)

	require.Equal(t, solution1.ID, solution2.ID)
	require.Equal(t, solution1.ProblemID, solution2.ProblemID)
	require.Equal(t, solution1.UserID, solution2.UserID)
	require.Equal(t, solution1.FileName, solution2.FileName)
	require.Equal(t, solution1.SavedPath, solution2.SavedPath)
	require.WithinDuration(t, solution1.SubmitedAt, solution2.SubmitedAt, time.Second)
	require.WithinDuration(t, solution1.UpdatedAt, solution2.UpdatedAt, time.Second)

	testQueries.DeleteSolution(context.Background(), solution1.ID)
	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestListSolutionsByUser(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	for i := 0; i < 5; i++ {
		createRandomSolution(t, user.ID, homework.ID)
	}

	arg := ListSolutionsByUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}

	solutions, err := testQueries.ListSolutionsByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solutions)

	for _, solution := range solutions {
		require.NotEmpty(t, solution)
	}

	for i := range solutions {
		testQueries.DeleteSolution(context.Background(), solutions[i].ID)
	}

	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestListSolutionsByProblem(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	for i := 0; i < 5; i++ {
		createRandomSolution(t, user.ID, homework.ID)
	}

	arg := ListSolutionsByProblemParams{
		ProblemID: homework.ID,
		Limit:     5,
		Offset:    0,
	}

	solutions, err := testQueries.ListSolutionsByProblem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solutions)

	for _, solution := range solutions {
		require.NotEmpty(t, solution)
	}

	for i := range solutions {
		testQueries.DeleteSolution(context.Background(), solutions[i].ID)
	}

	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestListSolutions(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	for i := 0; i < 5; i++ {
		createRandomSolution(t, user.ID, homework.ID)
	}

	arg := ListSolutionsParams{
		Limit:  5,
		Offset: 0,
	}

	solutions, err := testQueries.ListSolutions(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solutions)

	for _, solution := range solutions {
		require.NotEmpty(t, solution)
	}

	for i := range solutions {
		testQueries.DeleteSolution(context.Background(), solutions[i].ID)
	}

	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}

func TestUpdateSolution(t *testing.T) {
	teacher := createRandomTeacher(t)
	user := createRandomUser(t)
	subject := util.RandomSubject()
	homework := createRandomHomework(t, teacher.ID, subject)

	solution1 := createRandomSolution(t, user.ID, homework.ID)

	arg := UpdateSolutionParams{
		ID:        solution1.ID,
		FileName:  util.RandomString(10),
		SavedPath: util.RandomString(10),
		UpdatedAt: time.Now(),
	}

	solution2, err := testQueries.UpdateSolution(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, solution2)

	require.Equal(t, solution1.ID, solution2.ID)
	require.Equal(t, solution1.ProblemID, solution2.ProblemID)
	require.Equal(t, solution1.UserID, solution2.UserID)
	require.Equal(t, arg.FileName, solution2.FileName)
	require.Equal(t, arg.SavedPath, solution2.SavedPath)
	require.WithinDuration(t, solution1.SubmitedAt, solution2.SubmitedAt, time.Second)
	require.WithinDuration(t, arg.UpdatedAt, solution2.UpdatedAt, time.Second)

	testQueries.DeleteSolution(context.Background(), solution1.ID)
	testQueries.DeleteHomework(context.Background(), homework.ID)
	testQueries.DeleteUser(context.Background(), teacher.ID)
	testQueries.DeleteUser(context.Background(), user.ID)
}
