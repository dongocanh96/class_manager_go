package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/dongocanh96/class_manager_go/db/mock"
	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/token"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateMessage(t *testing.T) {
	user1 := randomUser(t)
	user2 := randomUser(t)
	message := randomMessage(t, user1.ID, user2.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user2.ID).Times(1).Return(user2, nil)

				arg := db.CreateMessageParams{
					FromUserID: user1.ID,
					ToUserID:   user2.ID,
					Content:    message.Content,
				}

				store.EXPECT().CreateMessage(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(message, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMessage(t, recorder.Body, message)
			},
		},
		{
			name: "MissingRequestParamError",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SelfSendMessage",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user2.ID, user2.Username.String, user2.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToUserNotFound",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "GetUserError",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CreateMessageError",
			body: gin.H{
				"to_user_id": user2.ID,
				"content":    message.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user2.ID).
					Times(1).
					Return(user2, nil)
				store.EXPECT().CreateMessage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Message{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/messages/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

type eqUpdateMessageStateParamsMatcher struct {
	arg db.UpdateMessageStateParams
}

func (e eqUpdateMessageStateParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateMessageStateParams)
	if !ok {
		return false
	}

	if !util.CheckDate(e.arg.ReadAt, arg.ReadAt) {
		return false
	}

	if e.arg.ID != arg.ID {
		return false
	}

	if e.arg.IsRead != arg.IsRead {
		return false
	}

	return true
}

func (e eqUpdateMessageStateParamsMatcher) String() string {
	return fmt.Sprintf("matches read at %v", e.arg.ReadAt)
}

func EqUpdateMessageStateParams(arg db.UpdateMessageStateParams) gomock.Matcher {
	return eqUpdateMessageStateParamsMatcher{arg}
}

func TestGetMessage(t *testing.T) {
	user1 := randomUser(t)
	user2 := randomUser(t)
	message := randomMessage(t, user1.ID, user2.ID)

	testCases := []struct {
		name          string
		messageID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchMessage(t, recoder.Body, message)
			},
		},
		{
			name:      "ReadedOk",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user2.ID, user2.Username.String, user2.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)

				arg := db.UpdateMessageStateParams{
					ID:     message.ID,
					IsRead: true,
					ReadAt: time.Now(),
				}

				readMessage := db.Message{
					ID:         message.ID,
					FromUserID: user1.ID,
					ToUserID:   user2.ID,
					Content:    message.Content,
					IsRead:     true,
					CreatedAt:  message.CreatedAt,
					ReadAt:     arg.ReadAt,
				}

				store.EXPECT().
					UpdateMessageState(gomock.Any(), EqUpdateMessageStateParams(arg)).
					Times(1).
					Return(readMessage, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireStateMatchMessage(t, recoder.Body)
			},
		},
		{
			name:      "InvalidID",
			messageID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:      "UnAuthorizedUser",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, 0, "UnAuthorizedUser", user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:      "MessageNotFound",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Message{}, sql.ErrNoRows)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name:      "InternalServerError",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), message.ID).
					Times(1).
					Return(db.Message{}, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name:      "ReadGetInternalServerError",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user2.ID, user2.Username.String, user2.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)

				store.EXPECT().
					UpdateMessageState(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Message{}, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
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

			url := fmt.Sprintf("/messages/%d", tc.messageID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListMessage(t *testing.T) {
	user1 := randomUser(t)
	user2 := randomUser(t)

	n := 5
	messages := make([]db.Message, n)
	for i := 0; i < n; i++ {
		messages[i] = randomMessage(t, user1.ID, user2.ID)
	}

	type Query struct {
		toUserID int64
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				toUserID: user2.ID,
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user2.ID).Times(1).
					Return(user2, nil)
				arg := db.ListMessagesParams{
					FromUserID: user1.ID,
					ToUserID:   user2.ID,
					Limit:      int32(n),
					Offset:     0,
				}

				store.EXPECT().ListMessages(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(messages, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "InvalidToUser",
			query: Query{
				toUserID: 0,
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				toUserID: user2.ID,
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				toUserID: user2.ID,
				pageID:   1,
				pageSize: 1000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
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

			url := "/messages/list_messages"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("to_user_id", fmt.Sprintf("%d", tc.query.toUserID))
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	user1 := randomUser(t)
	user2 := randomUser(t)
	message := randomMessage(t, user1.ID, user2.ID)
	content := util.RandomString(200)

	testCases := []struct {
		name          string
		messageID     int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			messageID: message.ID,
			body: gin.H{
				"content": content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)

				arg := db.UpdateMessageParams{
					ID:      message.ID,
					Content: content,
					IsRead:  false,
				}

				updatedMessage := db.Message{
					ID:         message.ID,
					FromUserID: user1.ID,
					ToUserID:   user2.ID,
					Content:    content,
					IsRead:     false,
					CreatedAt:  message.CreatedAt,
					ReadAt:     message.ReadAt,
				}
				store.EXPECT().UpdateMessage(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedMessage, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUpdatedMessage(t, recorder.Body, content)
			},
		},
		{
			name:      "MessageNotFound",
			messageID: message.ID,
			body: gin.H{
				"content": content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Message{}, sql.ErrNoRows)
				store.EXPECT().UpdateMessage(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			messageID: message.ID,
			body: gin.H{
				"content": content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user2.ID, user2.Username.String, user2.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).Return(message, nil)
				store.EXPECT().UpdateMessage(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/messages/%d/update", tc.messageID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	user1 := randomUser(t)
	user2 := randomUser(t)
	message := randomMessage(t, user1.ID, user2.ID)

	testCases := []struct {
		name          string
		messageID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)

				store.EXPECT().DeleteMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "UnauthorizedUser",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user2.ID, user2.Username.String, user2.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(1).
					Return(message, nil)

				store.EXPECT().DeleteMessage(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			messageID: message.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user1.ID, user1.Username.String, user1.IsTeacher,
					time.Minute*15)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMessage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Message{}, sql.ErrNoRows)

				store.EXPECT().DeleteMessage(gomock.Any(), gomock.Eq(message.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/messages/%d/delete", tc.messageID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomInt(1, 100),
		Username:       sql.NullString{String: util.RandomString(6), Valid: true},
		HashedPassword: hashedPassword,
		Fullname:       sql.NullString{String: util.RandomString(6), Valid: true},
		Email:          sql.NullString{String: util.RandomEmail(), Valid: true},
		PhoneNumber:    sql.NullString{String: util.RandomPhoneNumber(), Valid: true},
		IsTeacher:      util.RandomBoolean(),
	}

	return
}

func randomMessage(t *testing.T, user1ID int64, user2ID int64) (message db.Message) {
	return db.Message{
		ID:         util.RandomInt(1, 100),
		FromUserID: user1ID,
		ToUserID:   user2ID,
		Content:    util.RandomString(20),
	}
}

func requireBodyMatchMessage(t *testing.T, body *bytes.Buffer, message db.Message) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMessage db.Message
	err = json.Unmarshal(data, &gotMessage)
	require.NoError(t, err)
	require.Equal(t, message.FromUserID, gotMessage.FromUserID)
	require.Equal(t, message.ToUserID, gotMessage.ToUserID)
	require.Equal(t, message.Content, gotMessage.Content)
}

func requireStateMatchMessage(t *testing.T, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMessage db.Message
	err = json.Unmarshal(data, &gotMessage)
	require.NoError(t, err)
	require.True(t, gotMessage.IsRead)
}

func requireBodyMatchUpdatedMessage(t *testing.T, body *bytes.Buffer, content string) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var updatedMessage db.Message
	err = json.Unmarshal(data, &updatedMessage)
	require.NoError(t, err)
	require.False(t, updatedMessage.IsRead)
	require.Equal(t, updatedMessage.Content, content)
}
