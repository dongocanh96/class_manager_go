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
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
)

type createHomeworkRequest struct {
	Subject string                `form:"subject" binding:"required,subject"`
	Title   string                `form:"title" binding:"required,max=256"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}

func (server *Server) createHomework(ctx *gin.Context) {
	var req createHomeworkRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !authPayload.IsTeacher {
		err := errors.New("user is not teacher!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	savedPath := fmt.Sprintf("%s%s", server.config.Asset, req.File.Filename)
	err := ctx.SaveUploadedFile(req.File, savedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateHomeworkParams{
		TeacherID: authPayload.Userid,
		Subject:   req.Subject,
		Title:     req.Title,
		FileName:  req.File.Filename,
		SavedPath: savedPath,
	}

	homework, err := server.store.CreateHomework(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homework)
}

type getHomeworkRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getHomework(ctx *gin.Context) {
	var req getHomeworkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	homework, err := server.store.GetHomework(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homework)
}

type listHomeworkRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listHomework(ctx *gin.Context) {
	var req listHomeworkRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListHomeworksParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	homeworks, err := server.store.ListHomeworks(ctx, arg)
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

type listHomeworkBySubjectRequest struct {
	Subject  string `form:"subject" binding:"required,max=256"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listHomeworkBySubject(ctx *gin.Context) {
	var req listHomeworkBySubjectRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !util.IsSupportedSubject(req.Subject) {
		ctx.JSON(http.StatusBadRequest, errors.New("subject is not supported"))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListHomeworksBySubjectParams{
		Subject: req.Subject,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	homeworks, err := server.store.ListHomeworksBySubject(ctx, arg)
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

type updateHomeworkRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (server *Server) updateHomework(ctx *gin.Context) {
	var reqURI getHomeworkRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqForm updateHomeworkRequest
	if err := ctx.ShouldBind(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	homework, err := server.store.GetHomework(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if homework.IsClosed {
		err := errors.New("this homework is closed!")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if homework.TeacherID != authPayload.Userid {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	savedPath := fmt.Sprintf("%s%s", server.config.Asset, reqForm.File.Filename)
	err = ctx.SaveUploadedFile(reqForm.File, savedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldSavedPath := homework.SavedPath

	arg := db.UpdateHomeworkParams{
		ID:        reqURI.ID,
		FileName:  reqForm.File.Filename,
		SavedPath: savedPath,
		UpdatedAt: time.Now(),
	}

	homework, err = server.store.UpdateHomework(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = os.Remove(oldSavedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homework)

}

func (server *Server) closeHomework(ctx *gin.Context) {
	var req getHomeworkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	homework, err := server.store.GetHomework(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.Userid != homework.TeacherID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.CloseHomeworkParams{
		ID:       req.ID,
		IsClosed: true,
		ClosedAt: time.Now(),
	}

	homework, err = server.store.CloseHomework(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homework)
}

func (server *Server) deleteHomework(ctx *gin.Context) {
	var req getHomeworkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	homework, err := server.store.GetHomework(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if homework.TeacherID != authPayload.Userid {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = os.Remove(homework.SavedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteHomework(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type createSolutionRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (server *Server) createSolution(ctx *gin.Context) {
	var req createSolutionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getHomeworkRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	homework, err := server.store.GetHomework(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if homework.IsClosed {
		err := errors.New("homework is closed!")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	savedPath := fmt.Sprintf("%s%s", server.config.Asset, req.File.Filename)
	err = ctx.SaveUploadedFile(req.File, savedPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateSolutionParams{
		ProblemID: reqURI.ID,
		UserID:    authPayload.Userid,
		FileName:  req.File.Filename,
		SavedPath: savedPath,
	}

	solution, err := server.store.CreateSolution(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if authPayload.IsTeacher {
		argCloseHomework := db.CloseHomeworkParams{
			ID:       reqURI.ID,
			IsClosed: true,
			ClosedAt: time.Now(),
		}
		_, err := server.store.CloseHomework(ctx, argCloseHomework)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, solution)
}

type listSolutionsByProblemRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listSolutionsByProblem(ctx *gin.Context) {
	var reqForm listSolutionsByProblemRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getHomeworkRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !authPayload.IsTeacher {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListSolutionsByProblemParams{
		ProblemID: reqURI.ID,
		Limit:     reqForm.PageSize,
		Offset:    (reqForm.PageID - 1) * reqForm.PageSize,
	}

	solutions, err := server.store.ListSolutionsByProblem(ctx, arg)
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

type getSolutionsByUserRequest struct {
	UserID int64 `json:"userid" binding:"required,min=1"`
}

func (server *Server) getSolutionByProblemAndUser(ctx *gin.Context) {
	var req getSolutionsByUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqURI getHomeworkRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !authPayload.IsTeacher && authPayload.Userid != req.UserID {
		err := errors.New("permission denied!")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.GetSolutionByProblemAndUserParams{
		ProblemID: reqURI.ID,
		UserID:    req.UserID,
	}

	solution, err := server.store.GetSolutionByProblemAndUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, solution)
}
