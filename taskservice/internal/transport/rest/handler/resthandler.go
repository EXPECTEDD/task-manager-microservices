package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	taskdomain "taskservice/internal/domain/task"
	createdto "taskservice/internal/transport/rest/handler/dto/create"
	handlmapper "taskservice/internal/transport/rest/handler/mapper"
	handlvalidator "taskservice/internal/transport/rest/handler/validator"
	"taskservice/internal/usecase/interfaces"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	createUC interfaces.CreateTaskUsecase
}

func NewRestHandler(log *slog.Logger, createUC interfaces.CreateTaskUsecase) *RestHandler {
	return &RestHandler{
		log:      log,
		createUC: createUC,
	}
}

func (h *RestHandler) Create(ctx *gin.Context) {
	const op = "resthandler.Create"

	log := h.log.With(slog.String("op", op))

	log.Info("starting create request")

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

	in := handlmapper.CreateRequestToInput(&req)

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

func (h *RestHandler) Delete(ctx *gin.Context) {
	panic("not implemented")
}

func (h *RestHandler) GetAll(ctx *gin.Context) {
	panic("not implemented")
}

func (h *RestHandler) Update(ctx *gin.Context) {
	panic("not implemented")
}

func (h *RestHandler) Get(ctx *gin.Context) {
	panic("not implemented")
}
