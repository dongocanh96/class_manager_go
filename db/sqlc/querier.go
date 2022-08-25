// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CloseHomework(ctx context.Context, arg CloseHomeworkParams) (Homework, error)
	CreateHomework(ctx context.Context, arg CreateHomeworkParams) (Homework, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateSolution(ctx context.Context, arg CreateSolutionParams) (Solution, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteHomework(ctx context.Context, id int64) error
	DeleteMessage(ctx context.Context, id int64) error
	DeleteSolution(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetByUsername(ctx context.Context, username sql.NullString) (User, error)
	GetHomework(ctx context.Context, id int64) (Homework, error)
	GetMessage(ctx context.Context, id int64) (Message, error)
	GetSolutionByID(ctx context.Context, id int64) (Solution, error)
	GetSolutionByProblemAndUser(ctx context.Context, arg GetSolutionByProblemAndUserParams) (Solution, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserForUpdate(ctx context.Context, id int64) (User, error)
	ListHomeworks(ctx context.Context, arg ListHomeworksParams) ([]Homework, error)
	ListHomeworksBySubject(ctx context.Context, arg ListHomeworksBySubjectParams) ([]Homework, error)
	ListHomeworksByTeacher(ctx context.Context, arg ListHomeworksByTeacherParams) ([]Homework, error)
	ListMessages(ctx context.Context, arg ListMessagesParams) ([]Message, error)
	ListMessagesFromUser(ctx context.Context, arg ListMessagesFromUserParams) ([]Message, error)
	ListMessagesToUser(ctx context.Context, arg ListMessagesToUserParams) ([]Message, error)
	ListSolutions(ctx context.Context, arg ListSolutionsParams) ([]Solution, error)
	ListSolutionsByProblem(ctx context.Context, arg ListSolutionsByProblemParams) ([]Solution, error)
	ListSolutionsByUser(ctx context.Context, arg ListSolutionsByUserParams) ([]Solution, error)
	ListTeachersOrStudents(ctx context.Context, arg ListTeachersOrStudentsParams) ([]User, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateHomework(ctx context.Context, arg UpdateHomeworkParams) (Homework, error)
	UpdateMessage(ctx context.Context, arg UpdateMessageParams) (Message, error)
	UpdateMessageState(ctx context.Context, arg UpdateMessageStateParams) (Message, error)
	UpdateSolution(ctx context.Context, arg UpdateSolutionParams) (Solution, error)
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
}

var _ Querier = (*Queries)(nil)
