// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dongocanh96/class_manager_go/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CloseHomework mocks base method.
func (m *MockStore) CloseHomework(arg0 context.Context, arg1 db.CloseHomeworkParams) (db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseHomework", arg0, arg1)
	ret0, _ := ret[0].(db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseHomework indicates an expected call of CloseHomework.
func (mr *MockStoreMockRecorder) CloseHomework(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseHomework", reflect.TypeOf((*MockStore)(nil).CloseHomework), arg0, arg1)
}

// CreateHomework mocks base method.
func (m *MockStore) CreateHomework(arg0 context.Context, arg1 db.CreateHomeworkParams) (db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHomework", arg0, arg1)
	ret0, _ := ret[0].(db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHomework indicates an expected call of CreateHomework.
func (mr *MockStoreMockRecorder) CreateHomework(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHomework", reflect.TypeOf((*MockStore)(nil).CreateHomework), arg0, arg1)
}

// CreateMessage mocks base method.
func (m *MockStore) CreateMessage(arg0 context.Context, arg1 db.CreateMessageParams) (db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", arg0, arg1)
	ret0, _ := ret[0].(db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockStoreMockRecorder) CreateMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockStore)(nil).CreateMessage), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateSolution mocks base method.
func (m *MockStore) CreateSolution(arg0 context.Context, arg1 db.CreateSolutionParams) (db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSolution", arg0, arg1)
	ret0, _ := ret[0].(db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSolution indicates an expected call of CreateSolution.
func (mr *MockStoreMockRecorder) CreateSolution(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSolution", reflect.TypeOf((*MockStore)(nil).CreateSolution), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// DeleteHomework mocks base method.
func (m *MockStore) DeleteHomework(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHomework", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHomework indicates an expected call of DeleteHomework.
func (mr *MockStoreMockRecorder) DeleteHomework(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHomework", reflect.TypeOf((*MockStore)(nil).DeleteHomework), arg0, arg1)
}

// DeleteMessage mocks base method.
func (m *MockStore) DeleteMessage(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockStoreMockRecorder) DeleteMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockStore)(nil).DeleteMessage), arg0, arg1)
}

// DeleteSession mocks base method.
func (m *MockStore) DeleteSession(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockStoreMockRecorder) DeleteSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockStore)(nil).DeleteSession), arg0, arg1)
}

// DeleteSolution mocks base method.
func (m *MockStore) DeleteSolution(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSolution", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSolution indicates an expected call of DeleteSolution.
func (mr *MockStoreMockRecorder) DeleteSolution(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSolution", reflect.TypeOf((*MockStore)(nil).DeleteSolution), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// GetByUsername mocks base method.
func (m *MockStore) GetByUsername(arg0 context.Context, arg1 sql.NullString) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername.
func (mr *MockStoreMockRecorder) GetByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockStore)(nil).GetByUsername), arg0, arg1)
}

// GetHomework mocks base method.
func (m *MockStore) GetHomework(arg0 context.Context, arg1 int64) (db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHomework", arg0, arg1)
	ret0, _ := ret[0].(db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHomework indicates an expected call of GetHomework.
func (mr *MockStoreMockRecorder) GetHomework(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHomework", reflect.TypeOf((*MockStore)(nil).GetHomework), arg0, arg1)
}

// GetMessage mocks base method.
func (m *MockStore) GetMessage(arg0 context.Context, arg1 int64) (db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", arg0, arg1)
	ret0, _ := ret[0].(db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockStoreMockRecorder) GetMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockStore)(nil).GetMessage), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetSolutionByID mocks base method.
func (m *MockStore) GetSolutionByID(arg0 context.Context, arg1 int64) (db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionByID", arg0, arg1)
	ret0, _ := ret[0].(db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionByID indicates an expected call of GetSolutionByID.
func (mr *MockStoreMockRecorder) GetSolutionByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionByID", reflect.TypeOf((*MockStore)(nil).GetSolutionByID), arg0, arg1)
}

// GetSolutionByProblemAndUser mocks base method.
func (m *MockStore) GetSolutionByProblemAndUser(arg0 context.Context, arg1 db.GetSolutionByProblemAndUserParams) (db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSolutionByProblemAndUser", arg0, arg1)
	ret0, _ := ret[0].(db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSolutionByProblemAndUser indicates an expected call of GetSolutionByProblemAndUser.
func (mr *MockStoreMockRecorder) GetSolutionByProblemAndUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSolutionByProblemAndUser", reflect.TypeOf((*MockStore)(nil).GetSolutionByProblemAndUser), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 int64) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserForUpdate mocks base method.
func (m *MockStore) GetUserForUpdate(arg0 context.Context, arg1 int64) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserForUpdate", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserForUpdate indicates an expected call of GetUserForUpdate.
func (mr *MockStoreMockRecorder) GetUserForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserForUpdate", reflect.TypeOf((*MockStore)(nil).GetUserForUpdate), arg0, arg1)
}

// ListHomeworks mocks base method.
func (m *MockStore) ListHomeworks(arg0 context.Context, arg1 db.ListHomeworksParams) ([]db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHomeworks", arg0, arg1)
	ret0, _ := ret[0].([]db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHomeworks indicates an expected call of ListHomeworks.
func (mr *MockStoreMockRecorder) ListHomeworks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHomeworks", reflect.TypeOf((*MockStore)(nil).ListHomeworks), arg0, arg1)
}

// ListHomeworksBySubject mocks base method.
func (m *MockStore) ListHomeworksBySubject(arg0 context.Context, arg1 db.ListHomeworksBySubjectParams) ([]db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHomeworksBySubject", arg0, arg1)
	ret0, _ := ret[0].([]db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHomeworksBySubject indicates an expected call of ListHomeworksBySubject.
func (mr *MockStoreMockRecorder) ListHomeworksBySubject(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHomeworksBySubject", reflect.TypeOf((*MockStore)(nil).ListHomeworksBySubject), arg0, arg1)
}

// ListHomeworksByTeacher mocks base method.
func (m *MockStore) ListHomeworksByTeacher(arg0 context.Context, arg1 db.ListHomeworksByTeacherParams) ([]db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHomeworksByTeacher", arg0, arg1)
	ret0, _ := ret[0].([]db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHomeworksByTeacher indicates an expected call of ListHomeworksByTeacher.
func (mr *MockStoreMockRecorder) ListHomeworksByTeacher(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHomeworksByTeacher", reflect.TypeOf((*MockStore)(nil).ListHomeworksByTeacher), arg0, arg1)
}

// ListMessages mocks base method.
func (m *MockStore) ListMessages(arg0 context.Context, arg1 db.ListMessagesParams) ([]db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMessages", arg0, arg1)
	ret0, _ := ret[0].([]db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMessages indicates an expected call of ListMessages.
func (mr *MockStoreMockRecorder) ListMessages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMessages", reflect.TypeOf((*MockStore)(nil).ListMessages), arg0, arg1)
}

// ListMessagesFromUser mocks base method.
func (m *MockStore) ListMessagesFromUser(arg0 context.Context, arg1 db.ListMessagesFromUserParams) ([]db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMessagesFromUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMessagesFromUser indicates an expected call of ListMessagesFromUser.
func (mr *MockStoreMockRecorder) ListMessagesFromUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMessagesFromUser", reflect.TypeOf((*MockStore)(nil).ListMessagesFromUser), arg0, arg1)
}

// ListMessagesToUser mocks base method.
func (m *MockStore) ListMessagesToUser(arg0 context.Context, arg1 db.ListMessagesToUserParams) ([]db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMessagesToUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMessagesToUser indicates an expected call of ListMessagesToUser.
func (mr *MockStoreMockRecorder) ListMessagesToUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMessagesToUser", reflect.TypeOf((*MockStore)(nil).ListMessagesToUser), arg0, arg1)
}

// ListSolutionsByProblem mocks base method.
func (m *MockStore) ListSolutionsByProblem(arg0 context.Context, arg1 db.ListSolutionsByProblemParams) ([]db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSolutionsByProblem", arg0, arg1)
	ret0, _ := ret[0].([]db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSolutionsByProblem indicates an expected call of ListSolutionsByProblem.
func (mr *MockStoreMockRecorder) ListSolutionsByProblem(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSolutionsByProblem", reflect.TypeOf((*MockStore)(nil).ListSolutionsByProblem), arg0, arg1)
}

// ListSolutionsByUser mocks base method.
func (m *MockStore) ListSolutionsByUser(arg0 context.Context, arg1 db.ListSolutionsByUserParams) ([]db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSolutionsByUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSolutionsByUser indicates an expected call of ListSolutionsByUser.
func (mr *MockStoreMockRecorder) ListSolutionsByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSolutionsByUser", reflect.TypeOf((*MockStore)(nil).ListSolutionsByUser), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// UpdateHomework mocks base method.
func (m *MockStore) UpdateHomework(arg0 context.Context, arg1 db.UpdateHomeworkParams) (db.Homework, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHomework", arg0, arg1)
	ret0, _ := ret[0].(db.Homework)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateHomework indicates an expected call of UpdateHomework.
func (mr *MockStoreMockRecorder) UpdateHomework(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHomework", reflect.TypeOf((*MockStore)(nil).UpdateHomework), arg0, arg1)
}

// UpdateMessage mocks base method.
func (m *MockStore) UpdateMessage(arg0 context.Context, arg1 db.UpdateMessageParams) (db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessage", arg0, arg1)
	ret0, _ := ret[0].(db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessage indicates an expected call of UpdateMessage.
func (mr *MockStoreMockRecorder) UpdateMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockStore)(nil).UpdateMessage), arg0, arg1)
}

// UpdateMessageState mocks base method.
func (m *MockStore) UpdateMessageState(arg0 context.Context, arg1 db.UpdateMessageStateParams) (db.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessageState", arg0, arg1)
	ret0, _ := ret[0].(db.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessageState indicates an expected call of UpdateMessageState.
func (mr *MockStoreMockRecorder) UpdateMessageState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessageState", reflect.TypeOf((*MockStore)(nil).UpdateMessageState), arg0, arg1)
}

// UpdateSolution mocks base method.
func (m *MockStore) UpdateSolution(arg0 context.Context, arg1 db.UpdateSolutionParams) (db.Solution, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSolution", arg0, arg1)
	ret0, _ := ret[0].(db.Solution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSolution indicates an expected call of UpdateSolution.
func (mr *MockStoreMockRecorder) UpdateSolution(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSolution", reflect.TypeOf((*MockStore)(nil).UpdateSolution), arg0, arg1)
}

// UpdateUserInfo mocks base method.
func (m *MockStore) UpdateUserInfo(arg0 context.Context, arg1 db.UpdateUserInfoParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserInfo", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserInfo indicates an expected call of UpdateUserInfo.
func (mr *MockStoreMockRecorder) UpdateUserInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserInfo", reflect.TypeOf((*MockStore)(nil).UpdateUserInfo), arg0, arg1)
}

// UpdateUserInfoTx mocks base method.
func (m *MockStore) UpdateUserInfoTx(arg0 context.Context, arg1 db.UpdateUserInfoTxParams) (db.UpdateUserInfoTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserInfoTx", arg0, arg1)
	ret0, _ := ret[0].(db.UpdateUserInfoTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserInfoTx indicates an expected call of UpdateUserInfoTx.
func (mr *MockStoreMockRecorder) UpdateUserInfoTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserInfoTx", reflect.TypeOf((*MockStore)(nil).UpdateUserInfoTx), arg0, arg1)
}

// UpdateUserPassword mocks base method.
func (m *MockStore) UpdateUserPassword(arg0 context.Context, arg1 db.UpdateUserPasswordParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockStoreMockRecorder) UpdateUserPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockStore)(nil).UpdateUserPassword), arg0, arg1)
}
