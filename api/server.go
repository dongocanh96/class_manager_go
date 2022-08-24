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
	router.GET("/users/:id/solutions", server.listSolutionsByUser)
	router.GET("/users/:id/sended_messages", server.listSendedMessage)
	router.GET("/users/:id/recieved_messages", server.listReceivedMessages)

	//homework function
	router.POST("/homeworks/create", server.createHomework)
	router.GET("/homeworks/:id", server.getHomework)
	router.GET("/homeworks", server.listHomework)
	router.GET("/homeworks/teacher/:id", server.listHomeworkByTeacher)
	router.GET("/homeworks/subject", server.listHomeworkBySubject)
	router.PUT("/homeworks/:id", server.updateHomework)
	router.PUT("/homeworks/:id/close", server.closeHomework)
	router.DELETE("/homeworks/:id", server.deleteHomework)
	router.GET("homeworks/:id/solutions", server.listSolutionsByProblem)

	//solution function
	router.POST("/solutions/create", server.createSolution)
	router.GET("/solutions/:id", server.getSolutionByID)
	router.GET("/solutions/", server.listSolutions)
	router.GET("/solutions/by_homework_and_user", server.getSolutionByProblemAndUser)
	router.PUT("/solutions/:id", server.updateSolution)
	router.DELETE("/solutions/:id", server.deleteSolution)

	//message function
	router.POST("/messages/create", server.createMessage)
	router.GET("/messages/:id", server.getMessage)
	router.PUT("/messages/:id/update", server.updateMessage)
	router.PUT("/messages/:id/change_state", server.updateMessageState)
	router.GET("/messages/list_messages", server.listMessages)
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
