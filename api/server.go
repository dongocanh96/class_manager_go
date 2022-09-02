package api

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/token"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.JWTMaker
	router     *gin.Engine
}

func getKeyPairData(config util.Config) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	privateKeyByte, err := ioutil.ReadFile(config.PrivateKeyLocation)
	if err != nil {
		panic(err)
	}

	privatekey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		panic(err)
	}

	publicKeyByte, err := ioutil.ReadFile(config.PublicKeyLocation)
	if err != nil {
		panic(err)
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		panic(err)
	}

	return privatekey, publickey
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(getKeyPairData(config))
	if err != nil {
		return nil, fmt.Errorf("cannot create token %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: *tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("subject", validSubject)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	//user function
	router.POST("/users/create", server.createUser)
	router.POST("users/login", server.loginUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.PUT("/users/:id/update_info", server.updateUserInfo)
	authRoutes.PUT("/users/:id/update_password", server.updateUserPassword)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.GET("/users/:id/homeworks", server.listHomeworkByTeacher)
	authRoutes.GET("/users/:id/solutions", server.listSolutionsByUser)
	authRoutes.GET("/users/:id/sended_messages", server.listSendedMessage)
	authRoutes.GET("/users/:id/recieved_messages", server.listReceivedMessages)

	//homework function
	authRoutes.POST("/homeworks/create", server.createHomework)
	authRoutes.GET("/homeworks/:id", server.getHomework)
	authRoutes.GET("/homeworks", server.listHomework)
	authRoutes.GET("/homeworks/subject", server.listHomeworkBySubject)
	authRoutes.PUT("/homeworks/:id", server.updateHomework)
	authRoutes.PUT("/homeworks/:id/close", server.closeHomework)
	authRoutes.DELETE("/homeworks/:id", server.deleteHomework)
	authRoutes.POST("/homeworks/:id/solutions/create", server.createSolution)
	authRoutes.GET("homeworks/:id/solutions", server.listSolutionsByProblem)

	//solution function
	authRoutes.GET("/solutions/:id", server.getSolutionByID)
	authRoutes.GET("/solutions/by_homework_and_user", server.getSolutionByProblemAndUser)
	authRoutes.PUT("/solutions/:id", server.updateSolution)
	authRoutes.DELETE("/solutions/:id", server.deleteSolution)

	//message function
	authRoutes.POST("/messages/create", server.createMessage)
	authRoutes.GET("/messages/:id", server.getMessage)
	authRoutes.PUT("/messages/:id/update", server.updateMessage)
	authRoutes.GET("/messages/list_messages", server.listMessages)
	authRoutes.DELETE("/messages/:id/delete", server.deleteMessage)

	server.router = router
}
