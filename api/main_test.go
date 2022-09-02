package api

import (
	"os"
	"testing"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit((m.Run()))
}
