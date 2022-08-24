package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/dongocanh96/class_manager_go/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createSolutionRequest struct {
	ProblemID int64  `json:"problem_id" binding:"required,min=1"`
	UserID    int64  `json:"user_id" binding:"required,min=1"`
	FileName  string `json:"file_name" binding:"required"`
	SavedPath string `json:"saved_path" binding:"required"`
}

func (server *Server) createSolution(ctx *gin.Context) {
	var req createSolutionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateSolutionParams{
		ProblemID: req.ProblemID,
		UserID:    req.UserID,
		FileName:  req.FileName,
		SavedPath: req.SavedPath,
	}

	solution, err := server.store.CreateSolution(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, solution)
}

type getSolutionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getSolutionByID(ctx *gin.Context) {
	var req getSolutionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	solution, err := server.store.GetSolutionByID(ctx, req.ID)
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

type getSolutionsByProblemAndUser struct {
	ProblemID int64 `form:"problem_id" binding:"required,min=1"`
	UserID    int64 `form:"user_id" binding:"required,min=1"`
}

func (server *Server) getSolutionByProblemAndUser(ctx *gin.Context) {
	var req getSolutionsByProblemAndUser
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetSolutionByProblemAndUserParams{
		ProblemID: req.ProblemID,
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

type listSolutionsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listSolutions(ctx *gin.Context) {
	var req listHomeworkRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListSolutionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	solutions, err := server.store.ListSolutions(ctx, arg)
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

type updateSolutionRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	SavedPath string `json:"saved_path" binding:"required"`
}

func (server *Server) updateSolution(ctx *gin.Context) {
	var reqURI getSolutionRequest
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqJSON updateSolutionRequest
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateSolutionParams{
		ID:        reqURI.ID,
		FileName:  reqJSON.FileName,
		SavedPath: reqJSON.SavedPath,
		UpdatedAt: time.Now(),
	}

	solution, err := server.store.UpdateSolution(ctx, arg)
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

func (server *Server) deleteSolution(ctx *gin.Context) {
	var req getSolutionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteSolution(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
