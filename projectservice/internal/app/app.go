package app

import (
	"context"
	"database/sql"
	"log/slog"
	"projectservice/internal/config"
	userserviceclient "projectservice/internal/infrastructure/grpc/userservice"
	"projectservice/internal/infrastructure/postgres"
	grpcserver "projectservice/internal/transport/grpc"
	grpchandler "projectservice/internal/transport/grpc/handler"
	"projectservice/internal/transport/rest"
	resthandler "projectservice/internal/transport/rest/handler"
	"projectservice/internal/usecase/implementations/createproject"
	"projectservice/internal/usecase/implementations/deleteproject"
	"projectservice/internal/usecase/implementations/getallprojects"
	"projectservice/internal/usecase/implementations/updateproject"
	"projectservice/pkg/logger"
	"sync"
)

type App struct {
	log      *slog.Logger
	cfg      *config.Config
	restServ *rest.RestServer
	grpcServ *grpcserver.GRPCServer
	db       *sql.DB
	client   *userserviceclient.UserServiceClient
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
	updateProjectUC := updateproject.NewUpdateProjectUC(log, postgres)

	resthandl := resthandler.NewHandler(log, createProjectUC, deleteProjectUC, getAllProjectsUC, updateProjectUC)
	grpchandler := grpchandler.NewGRPCServer(log)

	restServ := mustLoadHttpServer(log, cfg, resthandl, client)
	grpcServ := mustLoadGRPCServer(log, cfg, grpchandler)

	return &App{
		log:      log,
		cfg:      cfg,
		restServ: restServ,
		grpcServ: grpcServ,
		db:       db,
		client:   client,
	}
}

func (a *App) Run() {
	go a.restServ.MustStart()
	a.grpcServ.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.restServ.Stop(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.grpcServ.Stop(ctx)
	}()

	wg.Wait()

	if err := a.db.Close(); err != nil {
		a.log.Error("db close failed", slog.String("error", err.Error()))
	}

	a.client.Stop()
}
