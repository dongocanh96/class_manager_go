package api

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/token"
	"github.com/gin-gonic/gin"
)

type getSolutionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSolutionByID(ctx *gin.Context) {
	var req getSolutionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	solution, err := server.store.GetSolutionByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !authPayload.IsTeacher && authPayload.Userid != solution.UserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, solution)
}

type updateSolutionRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (server *Server) updateSolution(ctx *gin.Context) {
	var reqURI getSolutionRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateSolutionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	solution, err := server.store.GetSolutionByID(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != solution.UserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	savedPath := fmt.Sprintf("%s%s", server.config.Asset, req.File.Filename)
	err = ctx.SaveUploadedFile(req.File, savedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldSavedPath := solution.SavedPath

	arg := db.UpdateSolutionParams{
		ID:        reqURI.ID,
		FileName:  req.File.Filename,
		SavedPath: savedPath,
		UpdatedAt: time.Now(),
	}

	solution, err = server.store.UpdateSolution(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = os.Remove(oldSavedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, solution)
}

func (server *Server) deleteSolution(ctx *gin.Context) {
	var req getSolutionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	solution, err := server.store.GetSolutionByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != solution.UserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = os.Remove(solution.SavedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteSolution(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
