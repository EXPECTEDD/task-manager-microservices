package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"taskservice/internal/repository/sessionvalidator"
	"time"

	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware(log *slog.Logger, sessionValid sessionvalidator.SessionValidator, respTimeout time.Duration) gin.HandlerFunc {
	const op = "middleware.SessionAuthMiddleware"
	return func(ctx *gin.Context) {
		sessionId, ok := ctx.Get("sessionId")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			ctx.Abort()
			return
		}

		tctx, cancel := context.WithTimeout(ctx.Request.Context(), respTimeout)
		defer cancel()

		userId, err := sessionValid.GetIdBySession(tctx, sessionId.(string))
		if err != nil {
			log.Info("failed to get userID", slog.String("op", op))
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(tctx.Err(), context.DeadlineExceeded) {
				ctx.JSON(http.StatusGatewayTimeout, gin.H{
					"error": "user service timeout",
				})
				ctx.Abort()
				return
			}
			errId, errStr := userServiceGrpcErrorToHttp(err)
			ctx.JSON(errId, gin.H{
				"error": errStr,
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}
