package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"taskservice/internal/repository/projectrepository"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckAccessToProjectMiddleware(log *slog.Logger, projectRepository projectrepository.ProjectRepository, respTimeout time.Duration) gin.HandlerFunc {
	const op = "middlware.GetOwnerIdMiddleware"
	return func(ctx *gin.Context) {
		log := log.With(slog.String("op", op))

		key, ok := ctx.Get("userId")
		if !ok {
			log.Error("failed get userId")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			ctx.Abort()
			return
		}
		userId, ok := key.(uint32)
		if !ok {
			log.Error("invalid userId type")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			ctx.Abort()
			return
		}
		if userId == 0 {
			log.Error("invalid userId")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
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

		tctx, cancel := context.WithTimeout(ctx.Request.Context(), respTimeout)
		defer cancel()

		ownerId, err := projectRepository.GetOwnerId(tctx, uint32(projectId))
		if err != nil {
			log.Info("failed to get ownerId")
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(tctx.Err(), context.DeadlineExceeded) {
				ctx.JSON(http.StatusGatewayTimeout, gin.H{
					"error": "project service timeout",
				})
				ctx.Abort()
				return
			}
			errId, errStr := projectServiceGrpcErrorToHttp(err)
			ctx.JSON(errId, gin.H{
				"error": errStr,
			})
			ctx.Abort()
			return
		}

		if userId != ownerId {
			log.Info("access denied")
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
