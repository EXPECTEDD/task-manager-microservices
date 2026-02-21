package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"taskservice/internal/config"
	"taskservice/internal/repository/sessionvalidator"
	"taskservice/internal/transport/rest"
	resthandler "taskservice/internal/transport/rest/handler"
	"taskservice/internal/transport/rest/middleware"

	"github.com/gin-gonic/gin"
)

func mustLoadRestServer(cfg *config.Config, log *slog.Logger, handl *resthandler.RestHandler, sessionValid sessionvalidator.SessionValidator) *rest.RestServer {
	gin.SetMode(cfg.RestConf.Mode)
	router := gin.New()

	group := router.Group("/")

	group.Use(gin.Recovery())
	group.Use(middleware.GetSessionMiddleware(log))
	group.Use(middleware.SessionAuthMiddleware(log, sessionValid, cfg.ConnectionsConf.UserServConnConf.ResponseTimeout))
	group.Use(middleware.TimeoutMiddleware(cfg.RestConf.RequestTimeout))

	group.POST("/task/create", handl.Create)
	group.DELETE("task/delete", handl.Delete)
	group.GET("/task/getall/:project_id", handl.GetAll)
	group.PATCH("/task/update/:task_id", handl.Update)
	group.GET("/task/get/:task_id", handl.Get)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.RestConf.Port),
		Handler:      router,
		WriteTimeout: cfg.RestConf.WriteTimeout,
		ReadTimeout:  cfg.RestConf.ReadTimeout,
	}

	return rest.NewRestServer(log, server)
}
