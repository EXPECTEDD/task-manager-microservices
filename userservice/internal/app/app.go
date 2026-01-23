package app

import (
	"context"
	"database/sql"
	"log/slog"
	"userservice/internal/config"
	bcrypthash "userservice/internal/infrastructure/bcrypt"
	"userservice/internal/infrastructure/postgres"
	myredis "userservice/internal/infrastructure/redis"
	uuidgen "userservice/internal/infrastructure/uuid"
	grpcserv "userservice/internal/transport/grpc"
	grpchandler "userservice/internal/transport/grpc/handler"
	"userservice/internal/transport/rest"
	resthandler "userservice/internal/transport/rest/handler"
	"userservice/internal/usecase/implementations/authenticate"
	"userservice/internal/usecase/implementations/login"
	"userservice/internal/usecase/implementations/registration"
	"userservice/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type App struct {
	log        *slog.Logger
	restServer *rest.RestServer
	grpcServer *grpcserv.GRPCServer
	cfg        *config.Config
	db         *sql.DB
	client     *redis.Client
}

func NewApp() *App {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.LogConf.Level)

	db := mustLoadPostgres(&cfg)
	client := mustLoadRedis(&cfg)

	pos := postgres.NewPostgres(db)
	hasher := bcrypthash.NewBcryptHasher()
	redis := myredis.NewRedis(client, &cfg.RedisConf.TTL)
	idgen := uuidgen.NewUUIDGenerator()

	regUC := registration.NewRegUC(log, pos, hasher)
	logUC := login.NewLoginUC(log, pos, hasher, redis, idgen)
	authUC := authenticate.NewAuthUC(log, redis)

	resthandl := resthandler.NewRestHandler(log, &cfg.RestConf.CookieTTL, regUC, logUC)
	grpchandl := grpchandler.NewGRPCHandler(log, cfg.GrpcConf.Timeout, authUC)

	restServer := mustLoadHttpServer(&cfg, log, resthandl)
	grpcserv := mustLoadGRPCServer(&cfg, log, grpchandl)

	return &App{
		log:        log,
		restServer: restServer,
		grpcServer: grpcserv,
		cfg:        &cfg,
		db:         db,
		client:     client,
	}
}

func (a *App) Run() {
	go a.restServer.MustStart()
	a.grpcServer.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout)
	defer cancel()

	a.restServer.Stop(ctx)
	a.grpcServer.Stop()

	a.db.Close()
	a.client.Close()
}
