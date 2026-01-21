package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusRequestTimeout, "timeout")
}

func TimeoutMiddleware(to time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(to),
		timeout.WithResponse(timeoutResponse),
	)
}
