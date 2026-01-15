package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userservice/internal/config"
	"userservice/internal/transport/rest"
	resthandler "userservice/internal/transport/rest/handler"
	"userservice/internal/transport/rest/middleware"
	"userservice/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.MustLoad()
	log := logger.SetupLogger(config.LogConf.Level)

	handl := resthandler.NewRestHandler(log)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(middleware.TimeoutMiddleware(config.RestConf.RequestTimeout))
	router.Use(gin.Recovery())

	router.POST("/registration", handl.Registration)
	router.POST("/login", handl.Login)

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.RestConf.Port),
		Handler:      router,
		WriteTimeout: config.RestConf.WriteTimeout,
		ReadTimeout:  config.RestConf.ReadTimeout,
	}

	restServer := rest.NewRestServer(log, serv)

	go restServer.MustStart()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	<-quitCh

	ctx, cancel := context.WithTimeout(context.Background(), config.RestConf.ShutdownTimeout*time.Second)
	defer cancel()

	restServer.Stop(ctx)
}
