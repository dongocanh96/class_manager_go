package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/token"
	"github.com/gin-gonic/gin"
)

type createMessageRequest struct {
	ToUserId int64  `json:"to_user_id" binding:"required,min=1"`
	Content  string `json:"content" binding:"required,max=5000"`
}

func (server *Server) createMessage(ctx *gin.Context) {
	var req createMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Userid == req.ToUserId {
		ctx.JSON(http.StatusBadRequest, errors.New("You can not send message to yourself!"))
		return
	}
	_, valid := server.validUser(ctx, req.ToUserId)
	if !valid {
		return
	}

	arg := db.CreateMessageParams{
		FromUserID: authPayload.Userid,
		ToUserID:   req.ToUserId,
		Content:    req.Content,
	}
	message, err := server.store.CreateMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

type getMessageRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getMessage(ctx *gin.Context) {
	var req getMessageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	message, err := server.store.GetMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != message.FromUserID && authPayload.Userid != message.ToUserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if authPayload.Userid == message.ToUserID && !message.IsRead {
		arg := db.UpdateMessageStateParams{
			ID:     req.ID,
			IsRead: true,
			ReadAt: time.Now(),
		}

		message, err = server.store.UpdateMessageState(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, message)
}

type listMessagesRequest struct {
	ToUserId int64 `form:"to_user_id" binding:"required,min=1"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listMessages(ctx *gin.Context) {
	var req listMessagesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	_, valid := server.validUser(ctx, req.ToUserId)
	if !valid {
		return
	}

	arg := db.ListMessagesParams{
		FromUserID: authPayload.Userid,
		ToUserID:   req.ToUserId,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}

	messages, err := server.store.ListMessages(ctx, arg)
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

type updateMessageRequest struct {
	Content string `json:"content" binding:"required,max=5000"`
}

func (server *Server) updateMessage(ctx *gin.Context) {
	var reqURI getMessageRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJSON updateMessageRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	message, err := server.store.GetMessage(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != message.FromUserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateMessageParams{
		ID:      reqURI.ID,
		Content: reqJSON.Content,
		IsRead:  false,
	}

	message, err = server.store.UpdateMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (server *Server) deleteMessage(ctx *gin.Context) {
	var req getMessageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	message, err := server.store.GetMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != message.FromUserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteMessage(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) validUser(ctx *gin.Context, userid int64) (db.User, bool) {
	user, err := server.store.GetUser(ctx, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return user, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return user, false
	}

	return user, true
}
