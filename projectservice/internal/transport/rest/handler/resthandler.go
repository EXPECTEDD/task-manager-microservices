package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	projectdomain "projectservice/internal/domain/project"
	createdto "projectservice/internal/transport/rest/handler/dto/create"
	updatedto "projectservice/internal/transport/rest/handler/dto/update"
	handlmapper "projectservice/internal/transport/rest/handler/mapper"
	handlvalidator "projectservice/internal/transport/rest/handler/validator"
	createerr "projectservice/internal/usecase/error/createproject"
	deleteerr "projectservice/internal/usecase/error/deleteproject"
	getallerr "projectservice/internal/usecase/error/getallprojects"
	updateerr "projectservice/internal/usecase/error/updateproject"
	"projectservice/internal/usecase/interfaces"
	getallmodel "projectservice/internal/usecase/models/getallprojects"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	createProjUC interfaces.CreateProjectUsecase
	deleteProjUC interfaces.DeleteProjectUsecase
	getAllProjUC interfaces.GetAllProjectsUsecase
	updateProjUC interfaces.UpdateProjectUsecase
}

func NewHandler(
	log *slog.Logger,
	createProjUC interfaces.CreateProjectUsecase,
	deleteProjUC interfaces.DeleteProjectUsecase,
	getAllProjUC interfaces.GetAllProjectsUsecase,
	updateProjUC interfaces.UpdateProjectUsecase,
) *RestHandler {
	return &RestHandler{
		log:          log,
		createProjUC: createProjUC,
		deleteProjUC: deleteProjUC,
		getAllProjUC: getAllProjUC,
		updateProjUC: updateProjUC,
	}
}

func (h *RestHandler) Create(ctx *gin.Context) {
	const op = "resthandler.Create"

	userId := getUserId(ctx)
	if userId == 0 {
		h.log.Error("failed to get userId", slog.String("op", op))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	log := h.log.With(slog.String("op", op), slog.Int("userId", int(userId)))

	log.Info("starting create request")

	var req *createdto.CreateRequest

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

	in := handlmapper.CreateRequestToInput(req, uint32(userId))

	out, err := h.createProjUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, projectdomain.ErrInvalidName) {
			log.Info("invalid name")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if errors.Is(err, projectdomain.ErrInvalidOwnerId) {
			log.Info("invalid owner id")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if errors.Is(err, createerr.ErrAlreadyExists) {
			log.Info("project already exists")
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("cannot create new project", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("create request completed successfully")

	res := handlmapper.CreateOutputToResponse(out)
	ctx.JSON(http.StatusOK, res)
}

func (h *RestHandler) Delete(ctx *gin.Context) {
	const op = "resthandler.Delete"

	userId := getUserId(ctx)
	if userId == 0 {
		h.log.Error("failed to get userId", slog.String("op", op))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	log := h.log.With(slog.String("op", op), slog.Int("userId", int(userId)))

	projcetIdStr := ctx.Param("project_id")
	projectId, err := strconv.Atoi(projcetIdStr)
	if err != nil {
		log.Error("project_id in context has invalid type")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project_id type",
		})
		return
	}
	if projectId <= 0 {
		log.Error("invalid project_id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project_id value",
		})
		return
	}

	log.Info("starting delete request")

	in := handlmapper.DeleteRequestToInput(userId, uint32(projectId))

	out, err := h.deleteProjUC.Execute(ctx, in)
	if err != nil {
		if errors.Is(err, projectdomain.ErrInvalidName) {
			log.Info("invalid project id")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if errors.Is(err, deleteerr.ErrProjectNotFound) {
			log.Info("project not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("cannot delete project", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("delete request completed successfully")

	res := handlmapper.DeleteOutputToResponse(out)
	ctx.JSON(http.StatusOK, res)
}

func (h *RestHandler) GetAll(ctx *gin.Context) {
	const op = "resthandler.GetAll"

	userId := getUserId(ctx)
	if userId == 0 {
		h.log.Error("failed to get userId", slog.String("op", op))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	log := h.log.With(slog.String("op", op), slog.Int("userId", int(userId)))

	log.Info("starting get all request")

	in := getallmodel.NewGetAllProjectsInput(userId)

	out, err := h.getAllProjUC.Execute(ctx, in)
	if err != nil {
		if errors.Is(err, getallerr.ErrProjectsNotFound) {
			log.Info("projects not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("cannot get projects", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	log.Info("get all request completed successfully")

	res := handlmapper.GetAllOutputToResponse(out)
	ctx.JSON(http.StatusOK, res)
}

func (h *RestHandler) Update(ctx *gin.Context) {
	const op = "resthandler.Update"

	userId := getUserId(ctx)
	if userId == 0 {
		h.log.Error("failed to get userId", slog.String("op", op))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	log := h.log.With(slog.String("op", op), slog.Int("userId", int(userId)))

	projectIdStr := ctx.Param("project_id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		log.Error("project_id in context has invalid type")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project_id type",
		})
		return
	}
	if projectId <= 0 {
		log.Error("invalid project_id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project_id value",
		})
		return
	}

	log.Info("starting updated request")

	var req *updatedto.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("error with request data", slog.String("error", err.Error()))
		errMap, ok := handlvalidator.MapValidationErrors(err)
		if ok {
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

	in := handlmapper.UpdateRequestToInput(userId, uint32(projectId), req)

	out, err := h.updateProjUC.Execute(ctx.Request.Context(), in)
	if err != nil {
		if errors.Is(err, updateerr.ErrInvalidName) {
			log.Info("invalid new project name")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid new project name",
			})
		} else if errors.Is(err, updateerr.ErrProjectNotFound) {
			log.Info("project not found")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "project not found",
			})
		} else if errors.Is(err, updateerr.ErrProjectNameAlreadyExists) {
			log.Info("project name already exists")
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "project name already exists",
			})
		}
		return
	}

	log.Info("update request completed successfully")
	res := handlmapper.UpdateOutputToResponse(out)
	ctx.JSON(http.StatusOK, res)
}

func getUserId(ctx *gin.Context) uint32 {
	if val, exists := ctx.Get("userId"); exists {
		return val.(uint32)
	} else {
		return 0
	}
}
