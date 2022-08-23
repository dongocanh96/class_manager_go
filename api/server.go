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

	//user function
	router.POST("/users/create", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)
	router.GET("/users/teacher", server.listTeacherOrStudent)
	router.GET("/users/student", server.listTeacherOrStudent)
	router.PUT("/users/:id/update_info", server.updateUserInfo)
	router.PUT("/users/:id/update_password", server.updateUserPassword)
	router.DELETE("/users/:id", server.deleteUser)
	router.GET("users/:id/sended_messages", server.listSendedMessage)
	router.GET("users/:id/recieved_messages", server.listRecievedMessages)

	//message function
	router.POST("/messages/create", server.createMessage)
	router.GET("/messages/:id", server.getMessage)
	router.PUT("/messages/:id/update", server.updateMessage)
	router.PUT("/message/:id/change_state", server.updateMessageState)
	router.GET("/messages/list_message", server.listMessages)
	router.DELETE("/messages/:id/delete", server.deleteMessage)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
