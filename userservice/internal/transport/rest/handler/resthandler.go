package resthandler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
	logindto "userservice/internal/transport/rest/handler/dto/login"
	regdto "userservice/internal/transport/rest/handler/dto/registration"
	handlmapper "userservice/internal/transport/rest/handler/mapper"
	handlvalidator "userservice/internal/transport/rest/handler/validator"
	logerr "userservice/internal/usecase/errors/login"
	regerr "userservice/internal/usecase/errors/registration"
	"userservice/internal/usecase/interfaces"

	"github.com/gin-gonic/gin"
)

type RestHandler struct {
	log       *slog.Logger
	cookieTTL time.Duration

	regUC interfaces.RegisterUserUsecase
	logUC interfaces.LoginUserUsecase
}

func NewRestHandler(log *slog.Logger, cookieTTL time.Duration, regUC interfaces.RegisterUserUsecase, logUC interfaces.LoginUserUsecase) *RestHandler {
	return &RestHandler{
		log:       log,
		cookieTTL: cookieTTL,
		regUC:     regUC,
		logUC:     logUC,
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

	in := handlmapper.RegRequestToInput(&regRequest)

	if ro, err := h.regUC.Execute(ctx.Request.Context(), in); err != nil {
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
		log.Info("registration request completed successfully")
		rr := handlmapper.RegOutputToResponse(ro)
		ctx.JSON(http.StatusOK, rr)
	}
}

func (h *RestHandler) Login(ctx *gin.Context) {
	const op = "resthandler.Login"
	log := h.log.With(slog.String("op", op))

	log.Info("start login request")

	var logRequest logindto.LoginRequest

	if err := ctx.ShouldBindJSON(&logRequest); err != nil {
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

	in := handlmapper.LogRequestToInput(&logRequest)

	if lo, err := h.logUC.Execute(ctx.Request.Context(), in); err != nil {
		if err != nil {
			if errors.Is(err, logerr.ErrUserNotFound) {
				log.Info("user not found")
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
			} else if errors.Is(err, logerr.ErrWrongPassword) {
				log.Info("wrong password")
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
			} else {
				log.Warn("an error occurred while executing the request", slog.String("error", err.Error()))
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}
	} else {
		log.Info("login request completed successfully")
		ctx.SetCookie("sessionId", lo.SessionId, int(h.cookieTTL.Seconds()), "/", "", false, true)
		lr := handlmapper.LogOutputToResponse(lo)
		ctx.JSON(http.StatusOK, gin.H{
			"user": lr,
		})
	}
}
