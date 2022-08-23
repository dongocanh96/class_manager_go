package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type createMessageRequest struct {
	FromUserId int64  `json:"from_user_id" binding:"required"`
	ToUserId   int64  `json:"to_user_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=5000"`
}

func (server *Server) createMessage(ctx *gin.Context) {
	var req createMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FromUserId == req.ToUserId {
		ctx.JSON(http.StatusBadRequest, errors.New("You can not send message to yourself!"))
		return
	}

	_, err := server.store.GetUser(ctx, req.ToUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateMessageParams{
		FromUserID: req.FromUserId,
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

	message, err := server.store.GetMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

type listMessagesRequestForm struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

type listMessagesRequestJSON struct {
	FromUserId int64 `json:"from_user_id" binding:"required,min=1"`
	ToUserId   int64 `json:"to_user_id" binding:"required,min=1"`
}

func (server *Server) listMessages(ctx *gin.Context) {
	var reqJSON listMessagesRequestJSON
	// if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	fmt.Println("vao day2")
	// 	return
	// }

	if err := ctx.ShouldBindBodyWith(&reqJSON, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		fmt.Println("vao day2")
		return
	}

	var reqForm listMessagesRequestForm
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMessagesParams{
		FromUserID: reqJSON.FromUserId,
		ToUserID:   reqJSON.ToUserId,
		Limit:      reqForm.PageSize,
		Offset:     (reqForm.PageID - 1) * reqForm.PageSize,
	}

	messages, err := server.store.ListMessages(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type listSendedMessageRequest struct {
	FromUserID int64 `form:"from_user_id" binding:"required,min=1"`
	PageID     int32 `form:"page_id" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listSendedMessage(ctx *gin.Context) {
	var req listSendedMessageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.FromUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListMessagesFromUserParams{
		FromUserID: req.FromUserID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
	}

	messages, err := server.store.ListMessagesFromUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

type listReceivedMessageRequest struct {
	ToUserID int64 `form:"to_user_id" binding:"required,min=1"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listReceivedMessages(ctx *gin.Context) {
	var req listReceivedMessageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.ToUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListMessagesToUserParams{
		ToUserID: req.ToUserID,
		Limit:    req.PageSize,
		Offset:   (req.PageID - 1) * req.PageSize,
	}

	messages, err := server.store.ListMessagesToUser(ctx, arg)
	if err != nil {
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

	_, err := server.store.GetMessage(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateMessageParams{
		ID:      reqURI.ID,
		Content: reqJSON.Content,
		IsRead:  false,
	}

	message, err := server.store.UpdateMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (server *Server) updateMessageState(ctx *gin.Context) {
	var req getMessageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateMessageStateParams{
		ID:     req.ID,
		IsRead: true,
		ReadAt: time.Now(),
	}

	message, err := server.store.UpdateMessageState(ctx, arg)
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

	_, err := server.store.GetMessage(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteMessage(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
