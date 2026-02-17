package app

import (
	"context"
	"database/sql"
	"log/slog"
	"projectservice/internal/config"
	userserviceclient "projectservice/internal/infrastructure/grpc/userservice"
	"projectservice/internal/infrastructure/postgres"
	"projectservice/internal/transport/rest"
	resthandler "projectservice/internal/transport/rest/handler"
	"projectservice/internal/usecase/implementations/createproject"
	"projectservice/internal/usecase/implementations/deleteproject"
	"projectservice/internal/usecase/implementations/getallprojects"
	"projectservice/internal/usecase/implementations/updateproject"
	"projectservice/pkg/logger"
)

type App struct {
	log    *slog.Logger
	cfg    *config.Config
	serv   *rest.RestServer
	db     *sql.DB
	client *userserviceclient.UserServiceClient
}

func NewApp() *App {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.LoggerConf.Level)
	db := mustLoadPostgres(cfg)

	client := userserviceclient.NewUserServiceClient(log, cfg.ConnectionsConf.UserServConnConf.Host, cfg.ConnectionsConf.UserServConnConf.Port)
	postgres := postgres.NewPostgres(db)

	createProjectUC := createproject.NewCreateProjectUC(log, postgres)
	deleteProjectUC := deleteproject.NewDeleteProjectUC(log, postgres)
	getAllProjectsUC := getallprojects.NewGetAllProjectsUC(log, postgres)
	updatePrijectUC := updateproject.NewUpdateProjectUC(log, postgres)

	handl := resthandler.NewHandler(log, createProjectUC, deleteProjectUC, getAllProjectsUC, updatePrijectUC)

	serv := mustLoadHttpServer(cfg, log, handl, client)

	return &App{
		log:    log,
		cfg:    cfg,
		serv:   serv,
		db:     db,
		client: client,
	}
}

func (a *App) Run() {
	a.serv.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout)
	defer cancel()
	a.serv.Stop(ctx)

	a.client.Stop()
	a.db.Close()
}
