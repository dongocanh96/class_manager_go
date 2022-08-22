package api

import (
	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) *Server {
	server := &Server{
		store:  store,
		config: config,
	}
	router := gin.Default()

	router.POST("/users/create", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)
	router.GET("/users/teacher", server.listTeacherOrStudent)
	router.GET("/users/student", server.listTeacherOrStudent)
	router.PUT("/users/:id/update_info", server.updateUserInfo)
	router.PUT("/users/:id/update_password", server.updateUserPassword)
	router.DELETE("/users/:id", server.deleteUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
