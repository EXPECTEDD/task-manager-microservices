package app

import (
	"context"
	"taskservice/internal/config"
	"taskservice/internal/transport/rest"
	"taskservice/pkg/logger"
)

type App struct {
	cfg        *config.Config
	restServer *rest.RestServer
}

func NewApp() *App {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.LoggerConf.Level)

	restServer := mustLoadRestServer(cfg, log)

	return &App{
		cfg:        cfg,
		restServer: restServer,
	}
}

func (a *App) Run() {
	a.restServer.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout)
	defer cancel()

	a.restServer.Stop(ctx)
}
