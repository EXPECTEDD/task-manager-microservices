package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"projectservice/internal/config"
	sessionalidator "projectservice/internal/repository/sessionvalidator"
	"projectservice/internal/transport/rest"
	resthandler "projectservice/internal/transport/rest/handler"
	"projectservice/internal/transport/rest/middleware"

	"github.com/gin-gonic/gin"
)

func mustLoadHttpServer(cfg *config.Config, log *slog.Logger, sessionValid sessionalidator.SessionValidator) *rest.RestServer {
	gin.SetMode(cfg.RestConf.Mode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.GetSessionMiddleware(log))
	router.Use(middleware.SessionAuthMiddleware(log, sessionValid, cfg.ConnectionsConf.UserServConnConf.ResponseTimeout))
	router.Use(middleware.TimeoutMiddleware(cfg.RestConf.RequestTimeout))
	handl := resthandler.NewHandler(log)

	router.POST("/project/create", handl.Create)
	router.DELETE("/project/delete", handl.Delete)
	router.GET("/project/getall", handl.GetAll)

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.RestConf.Port),
		Handler:      router,
		ReadTimeout:  cfg.RestConf.ReadTimeout,
		WriteTimeout: cfg.RestConf.WriteTimeout,
	}

	return rest.NewRestServer(log, serv)
}
