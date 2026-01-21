package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	tests := []struct {
		name   string
		timeMs time.Duration
		status int
		body   string
	}{
		{
			name:   "good test",
			timeMs: 1 * time.Millisecond,
			status: http.StatusOK,
			body:   "\"ok\"",
		},
		{
			name:   "timeout test",
			timeMs: 3 * time.Millisecond,
			status: http.StatusRequestTimeout,
			body:   "\"timeout\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(TimeoutMiddleware(2 * time.Millisecond))

			r.GET("/test", func(ctx *gin.Context) {
				time.Sleep(tt.timeMs)
				ctx.JSON(http.StatusOK, "ok")
			})

			serv := httptest.NewServer(r)
			defer serv.Close()

			resp, err := http.Get(serv.URL + "/test")
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.body, string(body))
			assert.Equal(t, tt.status, resp.StatusCode)
		})
	}
}
