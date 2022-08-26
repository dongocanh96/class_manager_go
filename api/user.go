package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Username:          user.Username.String,
		Fullname:          user.Fullname.String,
		Email:             user.Email.String,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type createUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fullname    string `json:"fullname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	TeacherKey  string `json:"teacher_key"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var IsTeacher bool
	if req.TeacherKey == server.config.SignUpKeyForTeacher {
		IsTeacher = true
	} else {
		IsTeacher = false
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	arg := db.CreateUserParams{
		Username:       sql.NullString{String: req.Username, Valid: true},
		HashedPassword: hashPassword,
		Fullname:       sql.NullString{String: req.Fullname, Valid: true},
		Email:          sql.NullString{String: req.Email, Valid: true},
		PhoneNumber:    sql.NullString{String: req.PhoneNumber, Valid: true},
		IsTeacher:      IsTeacher,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userName := sql.NullString{String: req.Username, Valid: true}

	user, err := server.store.GetByUsername(ctx, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		req.Username,
		user.IsTeacher,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	length := len(users)
	rsps := make([]userResponse, length)
	iter := 0
	for _, user := range users {
		rsps[iter] = newUserResponse(user)
		iter++
	}

	ctx.JSON(http.StatusOK, rsps)
}

type listUserByRoleRequest struct {
	IsTeacher *bool `form:"is_teacher" binding:"required"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listUserByRole(ctx *gin.Context) {
	var req listUserByRoleRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUserByRoleParams{
		IsTeacher: *req.IsTeacher,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUserByRole(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	length := len(users)
	rsps := make([]userResponse, length)
	iter := 0
	for _, user := range users {
		rsps[iter] = newUserResponse(user)
		iter++
	}

	ctx.JSON(http.StatusOK, rsps)
}

type updateUserInfoRequest struct {
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (server *Server) updateUserInfo(ctx *gin.Context) {
	var reqJSON updateUserInfoRequest
	var reqURI getUserRequest

	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var validUsername, validFullname, validEmail, validPhoneNumber bool

	if reqJSON.Username != "" {
		validUsername = true
	} else {
		validUsername = false
	}

	if reqJSON.Fullname != "" {
		validFullname = true
	} else {
		validFullname = false
	}

	if reqJSON.Email != "" {
		validEmail = true
	} else {
		validEmail = false
	}

	if reqJSON.PhoneNumber != "" {
		validPhoneNumber = true
	} else {
		validPhoneNumber = false
	}

	arg := db.UpdateUserInfoTxParams{
		ID:          reqURI.ID,
		Username:    sql.NullString{String: reqJSON.Username, Valid: validUsername},
		Fullname:    sql.NullString{String: reqJSON.Fullname, Valid: validFullname},
		Email:       sql.NullString{String: reqJSON.Email, Valid: validEmail},
		PhoneNumber: sql.NullString{String: reqJSON.PhoneNumber, Valid: validPhoneNumber},
	}

	user, err := server.store.UpdateUserInfoTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user.User)
	ctx.JSON(http.StatusOK, rsp)
}

type updateUserPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (server *Server) updateUserPassword(ctx *gin.Context) {
	var reqJSON updateUserPasswordRequest
	var reqURI getUserRequest

	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(reqJSON.OldPassword, user.HashedPassword); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashPassword, err := util.HashPassword(reqJSON.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:                reqURI.ID,
		HashedPassword:    hashPassword,
		PasswordChangedAt: time.Now(),
	}

	user, err = server.store.UpdateUserPassword(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type deleteUserRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

type listSolutionsByUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listSolutionsByUser(ctx *gin.Context) {
	var reqForm listSolutionsByUserRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getUserRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListSolutionsByUserParams{
		UserID: reqURI.ID,
		Limit:  reqForm.PageSize,
		Offset: (reqForm.PageID - 1) * reqForm.PageSize,
	}

	solutions, err := server.store.ListSolutionsByUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, solutions)
}

type listSendedMessageRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listSendedMessage(ctx *gin.Context) {
	var reqForm listSendedMessageRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getUserRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMessagesFromUserParams{
		FromUserID: reqURI.ID,
		Limit:      reqForm.PageSize,
		Offset:     (reqForm.PageID - 1) * reqForm.PageSize,
	}

	messages, err := server.store.ListMessagesFromUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type listReceivedMessageRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listReceivedMessages(ctx *gin.Context) {
	var reqForm listReceivedMessageRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getUserRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMessagesToUserParams{
		ToUserID: reqURI.ID,
		Limit:    reqForm.PageSize,
		Offset:   (reqForm.PageID - 1) * reqForm.PageSize,
	}

	messages, err := server.store.ListMessagesToUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type listHomeworkByTeacherRequestURI struct {
	TeacherID int64 `uri:"id" binding:"required,min=1"`
}

type listHomeworkByTeacherRequestForm struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listHomeworkByTeacher(ctx *gin.Context) {
	var reqURI listHomeworkByTeacherRequestURI
	var reqForm listHomeworkByTeacherRequestForm

	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListHomeworksByTeacherParams{
		TeacherID: reqURI.TeacherID,
		Limit:     reqForm.PageSize,
		Offset:    (reqForm.PageID - 1) * reqForm.PageSize,
	}

	homeworks, err := server.store.ListHomeworksByTeacher(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homeworks)
}
