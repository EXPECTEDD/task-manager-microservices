package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	projectdomain "projectservice/internal/domain/project"
	createdto "projectservice/internal/transport/rest/handler/dto/create"
	deletedto "projectservice/internal/transport/rest/handler/dto/delete"
	handlmapper "projectservice/internal/transport/rest/handler/mapper"
	handlvalidator "projectservice/internal/transport/rest/handler/validator"
	createerr "projectservice/internal/usecase/error/createproject"
	deleteerr "projectservice/internal/usecase/error/deleteproject"
	getallerr "projectservice/internal/usecase/error/getallprojects"
	"projectservice/internal/usecase/interfaces"
	getallmodel "projectservice/internal/usecase/models/getallprojects"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	createProjUC interfaces.CreateProjectUsecase
	deleteProjUC interfaces.DeleteProjectUsecase
	getAllProjUC interfaces.GetAllProjectsUsecase
}

func NewHandler(
	log *slog.Logger,
	createProjUC interfaces.CreateProjectUsecase,
	deleteProjUC interfaces.DeleteProjectUsecase,
	getAllProjUC interfaces.GetAllProjectsUsecase,
) *RestHandler {
	return &RestHandler{
		log:          log,
		createProjUC: createProjUC,
		deleteProjUC: deleteProjUC,
		getAllProjUC: getAllProjUC,
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
	}

	log := h.log.With(slog.String("op", op), slog.Int("userId", int(userId)))

	log.Info("starting delete request")

	var req *deletedto.DeleteRequest

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

	in := handlmapper.DeleteRequestToInput(req, userId)

	out, err := h.deleteProjUC.Execute(ctx, in)
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
	}

	log.Info("get all request completed successfully")

	res := handlmapper.GetAllOutputToResponse(out)
	ctx.JSON(http.StatusOK, res)
}

func getUserId(ctx *gin.Context) uint32 {
	if val, exists := ctx.Get("userId"); exists {
		return val.(uint32)
	} else {
		return 0
	}
}
