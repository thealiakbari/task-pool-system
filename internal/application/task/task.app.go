package service

import (
	"github.com/gin-gonic/gin"
	"github.com/thealiakbari/task-pool-system/internal/application/task/domain/dto"
	"github.com/thealiakbari/task-pool-system/internal/application/task/domain/transform"
	"github.com/thealiakbari/task-pool-system/internal/domain/task/pool"
	userInterface "github.com/thealiakbari/task-pool-system/internal/ports/inbound/task"
	"github.com/thealiakbari/task-pool-system/pkg/common/db"
	appErr "github.com/thealiakbari/task-pool-system/pkg/common/response"
)

type TaskHttpApp struct {
	userSvc          userInterface.TaskService
	poolWorkerHelper *pool.Pool
	db               db.DBWrapper
}

func NewTaskHttpApp(userSvc userInterface.TaskService, db db.DBWrapper, poolWorkerHelper *pool.Pool) TaskHttpApp {
	return TaskHttpApp{
		db:               db,
		userSvc:          userSvc,
		poolWorkerHelper: poolWorkerHelper,
	}
}

// MakeCreate
// @Schemes
// @Summary Create Task
// @Description This api for create task
// @Tags Task
// @Accept json
// @Produce json
// @Content-Type application/json
// @Param  body body dto.CreateTaskRequest true "Contains information to set data"
// @Success 201  {object}  dto.Task
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tasks [post]
func (t TaskHttpApp) MakeCreate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.CreateTaskRequest
		if err := ginCtx.ShouldBindJSON(&req); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}

		if err := req.Validate(ginCtx.Request.Context()); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EValidation,
			})
			return
		}

		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		pollEntityResp, err := t.userSvc.Create(ctx, transform.CreateTaskRequestToEntity(req))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		err = t.poolWorkerHelper.Submit(&pollEntityResp)
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.CreatedResponse(ginCtx, transform.TaskEntityToTaskDto(pollEntityResp))
	}
}

// MakeUpdate
// @Schemes
// @Summary Update Task
// @Description This api for update task
// @Tags Task
// @Accept json
// @Produce json
// @Content-Type application/json
// @Param id path string true "Task Id"
// @Param  body body dto.UpdateTaskRequest true "Contains information to set data"
// @Success 200  {object}  dto.Task
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tasks/{id} [put]
func (t TaskHttpApp) MakeUpdate() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var req dto.UpdateTaskRequest
		if err := ginCtx.ShouldBindJSON(&req); err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}

		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		updateReq, err := transform.UpdateTaskRequestToEntity(req, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EBadArg,
			})
			return
		}
		pollEntityResp, err := t.userSvc.Update(ctx, updateReq)
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.OKResponse(ginCtx, transform.TaskEntityToTaskDto(pollEntityResp))
	}
}

// MakeDelete
// @Schemes
// @Summary Delete Task
// @Description This api for delete task
// @Tags Task
// @Accept json
// @Produce json
// @Content-Type application/json
// @Param id path string true "Task Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tasks/{id} [delete]
func (t TaskHttpApp) MakeDelete() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
				appErr.HandelError(ginCtx, err)
				return
			}
		}()

		err = t.userSvc.Delete(ctx, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.NoContentResponse(ginCtx)
	}
}

// MakePurge
// @Schemes
// @Summary Purge Task
// @Description This api for purge task
// @Tags Task
// @Accept json
// @Produce json
// @Content-Type application/json
// @Param id path string true "Task Id"
// @Success 204
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tasks/purge/{id} [delete]
func (t TaskHttpApp) MakePurge() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tx, ctx, err := db.BeginTx(ginCtx.Request.Context(), t.db.DB)
		if err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		defer func() {
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					appErr.HandelError(ginCtx, &appErr.Error{
						Cause:   err,
						Message: err.Error(),
						Class:   appErr.EConflict,
					})
					return
				}
			}
		}()

		err = t.userSvc.Purge(ctx, ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		if err = tx.Commit().Error; err != nil {
			appErr.HandelError(ginCtx, &appErr.Error{
				Cause:   err,
				Message: err.Error(),
				Class:   appErr.EConflict,
			})
			return
		}

		appErr.NoContentResponse(ginCtx)
	}
}

// MakeGetById
// @Schemes
// @Summary Get Task By Id
// @Description This api for task by id
// @Tags Task
// @Accept json
// @Produce json
// @Content-Type application/json
// @Param id path string true "Task Id"
// @Success 200  {object} dto.Task
// @Failure 400  {object}  appErr.ErrSwaggerResponse
// @Failure 422  {object}  appErr.ErrValidationSwaggerResponse
// @Failure 500  {object}  appErr.ErrSwaggerResponse
// @Router /tasks/{id} [get]
func (t TaskHttpApp) MakeGetById() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		pollEntityResp, err := t.userSvc.GetByIdOrEmpty(ginCtx.Request.Context(), ginCtx.Param("id"))
		if err != nil {
			appErr.HandelError(ginCtx, err)
			return
		}

		appErr.OKResponse(ginCtx, transform.TaskEntityToTaskDto(pollEntityResp))
	}
}
