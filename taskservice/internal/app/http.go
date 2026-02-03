package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"taskservice/internal/config"
	"taskservice/internal/transport/rest"

	"github.com/gin-gonic/gin"
)

func mustLoadRestServer(cfg *config.Config, log *slog.Logger) *rest.RestServer {
	gin.SetMode(cfg.RestConf.Mode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello World")
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.RestConf.Port),
		Handler:      router,
		WriteTimeout: cfg.RestConf.WriteTimeout,
		ReadTimeout:  cfg.RestConf.ReadTimeout,
	}

	return rest.NewRestServer(log, server)
}
