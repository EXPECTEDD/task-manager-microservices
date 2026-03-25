package app

import (
	"context"
	"database/sql"
	"log/slog"
	"taskservice/internal/config"
	projectserviceclient "taskservice/internal/infrastructure/grpc/projectservice"
	userserviceclient "taskservice/internal/infrastructure/grpc/userservice"
	"taskservice/internal/infrastructure/postgres"
	"taskservice/internal/transport/rest"
	resthandler "taskservice/internal/transport/rest/handler"
	createuc "taskservice/internal/usecase/implementations/createtask"
	deleteuc "taskservice/internal/usecase/implementations/deletetask"
	getalluc "taskservice/internal/usecase/implementations/getalltasks"
	getuc "taskservice/internal/usecase/implementations/gettask"
	updateuc "taskservice/internal/usecase/implementations/updatetask"
	"taskservice/pkg/logger"
)

type App struct {
	log                  *slog.Logger
	cfg                  *config.Config
	restServer           *rest.RestServer
	userServiceClient    *userserviceclient.UserServiceClient
	projectServiceClient *projectserviceclient.ProjectServiceClient
	db                   *sql.DB
}

func NewApp() *App {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.LoggerConf.Level)

	db := mustLoadPostgres(cfg)

	postgres := postgres.NewPostgres(db)

	createUC := createuc.NewCreateTaskUC(log, postgres)
	updateUC := updateuc.NewUpdateTaskUC(log, postgres)
	deleteUC := deleteuc.NewDeleteTaskUC(log, postgres)
	getAllUC := getalluc.NewGetAllTasksUC(log, postgres)
	getUC := getuc.NewGetTaskUC(log, postgres)

	userServiceClient := userserviceclient.NewUserServiceClient(log, cfg.ConnectionsConf.UserServConnConf.Host, cfg.ConnectionsConf.UserServConnConf.Port)
	projectServiceClient := projectserviceclient.NewProjectServiceClient(log, cfg.ConnectionsConf.ProjServConnConf.Host, cfg.ConnectionsConf.ProjServConnConf.Port)
	handl := resthandler.NewRestHandler(log, createUC, updateUC, deleteUC, getAllUC, getUC)

	restServer := mustLoadRestServer(cfg, log, handl, userServiceClient, projectServiceClient)

	return &App{
		log:                  log,
		cfg:                  cfg,
		restServer:           restServer,
		userServiceClient:    userServiceClient,
		projectServiceClient: projectServiceClient,
		db:                   db,
	}
}

func (a *App) Run() {
	a.restServer.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout)
	defer cancel()

	a.restServer.Stop(ctx)

	if err := a.db.Close(); err != nil {
		a.log.Error("db close failed", slog.String("error", err.Error()))
	}

	a.userServiceClient.Stop()
	a.projectServiceClient.Stop()
}
