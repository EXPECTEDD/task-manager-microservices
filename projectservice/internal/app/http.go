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

func mustLoadHttpServer(cfg *config.Config, log *slog.Logger, handl *resthandler.RestHandler, sessionValid sessionalidator.SessionValidator) *rest.RestServer {
	gin.SetMode(cfg.RestConf.Mode)
	router := gin.New()

	group := router.Group("/")

	group.Use(gin.Recovery())
	group.Use(middleware.GetSessionMiddleware(log))
	group.Use(middleware.SessionAuthMiddleware(log, sessionValid, cfg.ConnectionsConf.UserServConnConf.ResponseTimeout))
	group.Use(middleware.TimeoutMiddleware(cfg.RestConf.RequestTimeout))

	group.POST("/project/create", handl.Create)
	group.DELETE("/project/delete/:project_id", handl.Delete)
	group.GET("/project/getall", handl.GetAll)
	group.PATCH("/project/update/:project_id", handl.Update)

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.RestConf.Port),
		Handler:      router,
		ReadTimeout:  cfg.RestConf.ReadTimeout,
		WriteTimeout: cfg.RestConf.WriteTimeout,
	}

	return rest.NewRestServer(log, serv)
}
