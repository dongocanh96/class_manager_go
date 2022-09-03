package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/dongocanh96/class_manager_go/db/mock"
	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)

}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {

	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserApi(t *testing.T) {
	student, studentPassword := randomStudentUser(t)
	teacher, teacherPassword := randomTeacherUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Create Student OK",
			body: gin.H{
				"username":     student.Username.String,
				"password":     studentPassword,
				"fullname":     student.Fullname.String,
				"email":        student.Email.String,
				"phone_number": student.PhoneNumber.String,
				"teacher_key":  "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    student.Username,
					Fullname:    student.Fullname,
					Email:       student.Email,
					PhoneNumber: student.PhoneNumber,
					IsTeacher:   student.IsTeacher,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, studentPassword)).
					Times(1).
					Return(student, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchStudent(t, recorder.Body, student)
			},
		},
		{
			name: "Create teacher OK",
			body: gin.H{
				"username":     teacher.Username.String,
				"password":     teacherPassword,
				"fullname":     teacher.Fullname.String,
				"email":        teacher.Email.String,
				"phone_number": teacher.PhoneNumber.String,
				"teacher_key":  "5WC7CnJ99KBhyPF",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    teacher.Username,
					Fullname:    teacher.Fullname,
					Email:       teacher.Email,
					PhoneNumber: teacher.PhoneNumber,
					IsTeacher:   true,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, teacherPassword)).
					Times(1).
					Return(teacher, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTeacher(t, recorder.Body, teacher)
			},
		},
		{
			name: "PasswordTooShort",
			body: gin.H{
				"username":     teacher.Username.String,
				"password":     "dfgh",
				"fullname":     teacher.Fullname.String,
				"email":        teacher.Email.String,
				"phone_number": teacher.PhoneNumber.String,
				"teacher_key":  "5WC7CnJ99KBhyPF",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateName",
			body: gin.H{
				"username":     teacher.Username.String,
				"password":     teacherPassword,
				"fullname":     teacher.Fullname.String,
				"email":        teacher.Email.String,
				"phone_number": teacher.PhoneNumber.String,
				"teacher_key":  "5WC7CnJ99KBhyPF",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidName",
			body: gin.H{
				"username":     "hellop#1",
				"password":     teacherPassword,
				"fullname":     teacher.Fullname.String,
				"email":        teacher.Email.String,
				"phone_number": teacher.PhoneNumber.String,
				"teacher_key":  "5WC7CnJ99KBhyPF",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user, password := randomStudentUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username.String,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username.String,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username.String,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#1",
				"password":  password,
				"full_name": user.Fullname.String,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomStudentUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomInt(1, 100),
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashedPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
	}

	return
}

func requireBodyMatchStudent(t *testing.T, body *bytes.Buffer, student db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser struct {
		ID                int64     `json:"id"`
		Username          string    `json:"username"`
		Fullname          string    `json:"fullname"`
		Email             string    `json:"email"`
		PhoneNumber       string    `json:"phone_number"`
		IsTeacher         bool      `json:"is_teacher"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
		CreatedAt         time.Time `json:"created_at"`
	}
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, student.Username.String, gotUser.Username)
	require.Equal(t, student.Fullname.String, gotUser.Fullname)
	require.Equal(t, student.Email.String, gotUser.Email)
	require.Equal(t, student.PhoneNumber.String, gotUser.PhoneNumber)
	require.Equal(t, false, gotUser.IsTeacher)

}

func randomTeacherUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashedPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      true,
	}

	return
}

func requireBodyMatchTeacher(t *testing.T, body *bytes.Buffer, teacher db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser struct {
		ID                int64     `json:"id"`
		Username          string    `json:"username"`
		Fullname          string    `json:"fullname"`
		Email             string    `json:"email"`
		PhoneNumber       string    `json:"phone_number"`
		IsTeacher         bool      `json:"is_teacher"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
		CreatedAt         time.Time `json:"created_at"`
	}
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, teacher.Username.String, gotUser.Username)
	require.Equal(t, teacher.Fullname.String, gotUser.Fullname)
	require.Equal(t, teacher.Email.String, gotUser.Email)
	require.Equal(t, teacher.PhoneNumber.String, gotUser.PhoneNumber)
	require.Equal(t, true, gotUser.IsTeacher)

}
