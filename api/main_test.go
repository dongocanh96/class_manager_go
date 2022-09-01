package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/token"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SignUpKeyForTeacher: "5WC7CnJ99KBhyPF",
		PrivateKeyLocation:  "../private.pem",
		PublicKeyLocation:   "../public.pem",
		AccessTokenDuration: time.Minute * 15,
		Asset:               "./asset/",
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func createAndSetAuthToken(t *testing.T, request *http.Request,
	tokenMaker token.JWTMaker, userid int64, username string, isteacher bool) {

	if len(username) == 0 || userid == 0 {
		return
	}

	token, err := tokenMaker.CreateToken(userid, username, isteacher, time.Minute*15)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit((m.Run()))
}
