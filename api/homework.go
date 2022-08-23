package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/gin-gonic/gin"
)

type createHomeworkRequest struct {
	TeacherID int64  `json:"teacher_id" binding:"required,min=1"`
	Subject   string `json:"subject" binding:"required,max=256"`
	Title     string `json:"title" binding:"required,max=256"`
	FileName  string `json:"file_name" binding:"required"`
	SavedPath string `json:"saved_path" binding:"required"`
}

func (server *Server) createHomework(ctx *gin.Context) {
	var req createHomeworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateHomeworkParams{
		TeacherID: req.TeacherID,
		Subject:   req.Subject,
		Title:     req.Title,
		FileName:  req.FileName,
		SavedPath: req.SavedPath,
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

	arg := db.ListHomeworksParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	homeworks, err := server.store.ListHomeworks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, homeworks)
}

type listHomeworkByTeacherRequest struct {
	TeacherID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listHomeworkByTeacher(ctx *gin.Context) {
	var reqURI listHomeworkByTeacherRequest
	var reqForm listHomeworkRequest

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
	FileName  string `json:"file_name" binding:"required"`
	SavedPath string `json:"saved_path" binding:"required"`
}

func (server *Server) updateHomework(ctx *gin.Context) {
	var reqURI getHomeworkRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJSON updateHomeworkRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetHomework(ctx, reqURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateHomeworkParams{
		ID:        reqURI.ID,
		FileName:  reqJSON.FileName,
		SavedPath: reqJSON.SavedPath,
		UpdatedAt: time.Now(),
	}

	homework, err := server.store.UpdateHomework(ctx, arg)
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

	_, err := server.store.GetHomework(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CloseHomeworkParams{
		ID:       req.ID,
		IsClosed: true,
		ClosedAt: time.Now(),
	}

	homework, err := server.store.CloseHomework(ctx, arg)
	if err != nil {
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

	_, err := server.store.GetHomework(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

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
