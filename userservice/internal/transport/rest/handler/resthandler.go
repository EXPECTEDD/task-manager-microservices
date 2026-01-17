package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	regdto "userservice/internal/transport/rest/handler/dto/registration"
	handlmapper "userservice/internal/transport/rest/handler/mapper"
	handlvalidator "userservice/internal/transport/rest/handler/validator"
	regerr "userservice/internal/usecase/errors/registration"
	"userservice/internal/usecase/interfaces"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log *slog.Logger

	regUC interfaces.RegistrationUsecase
}

func NewRestHandler(log *slog.Logger, regUC interfaces.RegistrationUsecase) *RestHandler {
	return &RestHandler{
		log:   log,
		regUC: regUC,
	}
}

func (h *RestHandler) Registration(ctx *gin.Context) {
	const op = "resthandler.Registration"
	log := h.log.With(slog.String("op", op))

	log.Info("start registration request")

	var regRequest regdto.RegistrationRequest

	if err := ctx.ShouldBindJSON(&regRequest); err != nil {
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

	ri := handlmapper.RegRequestToInput(&regRequest)

	if ro, err := h.regUC.RegUser(ctx.Request.Context(), ri); err != nil {
		if errors.Is(err, regerr.ErrUserAlreadyExists) {
			log.Info("user already exists")
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Warn("an error occurred while executing the request", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
	} else {
		log.Info("request completed successfully")
		rr := handlmapper.RegOutputToResponse(ro)
		ctx.JSON(http.StatusOK, rr)
	}
}

func (h *RestHandler) Login(ctx *gin.Context) {
	panic("not implemented")
}
