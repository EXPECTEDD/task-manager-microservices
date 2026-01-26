package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"userservice/internal/config"
	"userservice/internal/transport/rest"
	resthandler "userservice/internal/transport/rest/handler"
	"userservice/internal/transport/rest/middleware"

	"github.com/gin-gonic/gin"
)

func mustLoadHttpServer(cfg *config.Config, log *slog.Logger, handl *resthandler.RestHandler) *rest.RestServer {
	// GIN SETTINGS
	gin.SetMode(cfg.RestConf.Mode)
	router := gin.New()
	router.Use(middleware.TimeoutMiddleware(cfg.RestConf.RequestTimeout))
	router.Use(gin.Recovery())

	// REGISTER HTTP ROUTES
	router.POST("/user/registration", handl.Registration)
	router.POST("/user/login", handl.Login)

	// SERVER SETTING
	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.RestConf.Port),
		Handler:      router,
		WriteTimeout: cfg.RestConf.WriteTimeout,
		ReadTimeout:  cfg.RestConf.ReadTimeout,
	}

	restServer := rest.NewRestServer(log, serv)

	return restServer
}
