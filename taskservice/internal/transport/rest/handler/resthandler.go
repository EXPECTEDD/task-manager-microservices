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
	deleteerr "taskservice/internal/usecase/error/deletetask"
	getallerr "taskservice/internal/usecase/error/getalltasks"
	geterr "taskservice/internal/usecase/error/gettask"
	updatetaskerr "taskservice/internal/usecase/error/updatetask"
	"taskservice/internal/usecase/interfaces"
	deletemodel "taskservice/internal/usecase/models/deletetask"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
	getmodel "taskservice/internal/usecase/models/gettask"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	createUC interfaces.CreateTaskUsecase
	updateUC interfaces.UpdateTaskUsecase
	deleteUC interfaces.DeleteTaskUsecase
	getAllUC interfaces.GetAllTasksUsecase
	getUC    interfaces.GetTaskUsecase
}

func NewRestHandler(log *slog.Logger, createUC interfaces.CreateTaskUsecase, updateUC interfaces.UpdateTaskUsecase, deleteUC interfaces.DeleteTaskUsecase, getAllUC interfaces.GetAllTasksUsecase, getUC interfaces.GetTaskUsecase) *RestHandler {
	return &RestHandler{
		log:      log,
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		getAllUC: getAllUC,
		getUC:    getUC,
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

	log.Info("starting update request")

	taskIdStr := ctx.Param("task_id")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if taskId == 0 || err != nil {
		log.Warn("cannot get task id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		ctx.Abort()
		return
	}

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

	in := handlmapper.UpdateRequestToInput(&req, uint32(taskId), uint32(projectId))

	out, err := h.updateUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, updatetaskerr.ErrTaskNotFound) {
			log.Info("task not found")
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
	const op = "resthandler.Delete"

	log := h.log.With(slog.String("op", op))

	log.Info("starting delete request")

	taskIdStr := ctx.Param("task_id")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if taskId == 0 || err != nil {
		log.Warn("cannot get task id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		ctx.Abort()
		return
	}

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

	in := deletemodel.NewDeleteTaskInput(uint32(taskId), uint32(projectId))

	out, err := h.deleteUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, deleteerr.ErrTaskNotFound) {
			log.Info("task not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		} else {
			log.Warn("cannot delete task", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("delete task completed successfully")

	resp := handlmapper.DeleteOutputToResponse(out)
	ctx.JSON(http.StatusOK, resp)
}

func (h *RestHandler) GetAll(ctx *gin.Context) {
	const op = "resthandler.GetAll"

	log := h.log.With(slog.String("op", op))

	log.Info("starting get all tasks request")

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

	in := getallmodel.NewGetAllTasksInput(uint32(projectId))

	out, err := h.getAllUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, getallerr.ErrTasksNotFound) {
			log.Info("tasks not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "tasks not found",
			})
		} else {
			log.Warn("cannot get all tasks", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("get all tasks completed successfully")
	resp := handlmapper.GetAllOutputToResponse(out)
	ctx.JSON(http.StatusOK, resp)
}

func (h *RestHandler) Get(ctx *gin.Context) {
	const op = "resthandler.Get"

	log := h.log.With(slog.String("op", op))

	log.Info("starting get task request")

	taskIdStr := ctx.Param("task_id")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if taskId == 0 || err != nil {
		log.Warn("cannot get task id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		ctx.Abort()
		return
	}

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

	in := getmodel.NewGetTaskInput(uint32(taskId), uint32(projectId))

	out, err := h.getUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, geterr.ErrTaskNotFound) {
			log.Info("task not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		} else {
			log.Warn("cannot get task", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("get task completed successfully")
	resp := handlmapper.GetOutputToResponse(out)
	ctx.JSON(http.StatusOK, resp)
}
