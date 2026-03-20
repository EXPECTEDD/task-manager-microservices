package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	taskdomain "taskservice/internal/domain/task"
	createdto "taskservice/internal/transport/rest/handler/dto/create"
	updatedto "taskservice/internal/transport/rest/handler/dto/update"
	handlmapper "taskservice/internal/transport/rest/handler/mapper"
	handlvalidator "taskservice/internal/transport/rest/handler/validator"
	updatetaskerr "taskservice/internal/usecase/error/updatetask"
	"taskservice/internal/usecase/interfaces"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	createUC interfaces.CreateTaskUsecase
	updateUC interfaces.UpdateTaskUsecase
}

func NewRestHandler(log *slog.Logger, createUC interfaces.CreateTaskUsecase, updateUC interfaces.UpdateTaskUsecase) *RestHandler {
	return &RestHandler{
		log:      log,
		createUC: createUC,
		updateUC: updateUC,
	}
}

func (h *RestHandler) Create(ctx *gin.Context) {
	const op = "resthandler.Create"

	log := h.log.With(slog.String("op", op))

	log.Info("starting create request")

	projectIdStr := ctx.Param("project_id")
	projectId, err := strconv.ParseUint(projectIdStr, 10, 32)
	if projectId == 0 || err != nil {
		log.Warn("cannot get project id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		ctx.Abort()
		return
	}

	var req createdto.CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("error with request data", slog.String("error", err.Error()))
		if errMap, ok := handlvalidator.MapValidationErrors(err); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": errMap,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request body",
			})
		}
		return
	}

	in := handlmapper.CreateRequestToInput(&req, uint32(projectId))

	out, err := h.createUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, taskdomain.ErrInvalidProjectId) {
			log.Info("invalid project id")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if errors.Is(err, taskdomain.ErrInvalidDescription) {
			log.Info("invalid description")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("cannot create new task", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("create request completed successfully")

	resp := handlmapper.CreateOutputToResponse(out)
	ctx.JSON(http.StatusOK, resp)
}

func (h *RestHandler) Update(ctx *gin.Context) {
	const op = "resthandler.Update"

	log := h.log.With(slog.String("op", op))

	taskIdStr := ctx.Param("task_id")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if taskId == 0 || err != nil {
		log.Warn("cannot get project id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		ctx.Abort()
		return
	}

	var req updatedto.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("error with request data", slog.String("error", err.Error()))
		if errMap, ok := handlvalidator.MapValidationErrors(err); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": errMap,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request body",
			})
		}
		return
	}

	if req.NewDeadline == nil && req.NewDescription == nil {
		log.Info("empty body")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "nothing to update",
		})
		return
	}

	if req.NewDescription != nil && *req.NewDescription == "" {
		log.Info("invalid new description")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid new description",
		})
		return
	}

	in := handlmapper.UpdateRequestToInput(&req, uint32(taskId))

	out, err := h.updateUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, updatetaskerr.ErrTaskNotFound) {
			log.Info("project not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("cannot update task", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("update task completed successfully")

	resp := handlmapper.UpdateOutputToResponse(out)
	ctx.JSON(http.StatusOK, resp)
}

func (h *RestHandler) Delete(ctx *gin.Context) {
	panic("not implemented")
}

func (h *RestHandler) GetAll(ctx *gin.Context) {
	panic("not implemented")
}

func (h *RestHandler) Get(ctx *gin.Context) {
	panic("not implemented")
}
