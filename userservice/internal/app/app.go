package app

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
	"userservice/internal/config"
	bcrypthash "userservice/internal/infrastructure/bcrypt"
	"userservice/internal/infrastructure/postgres"
	myredis "userservice/internal/infrastructure/redis"
	uuidgen "userservice/internal/infrastructure/uuid"
	"userservice/internal/transport/rest"
	resthandler "userservice/internal/transport/rest/handler"
	"userservice/internal/usecase/implementations/login"
	"userservice/internal/usecase/implementations/registration"

	"github.com/redis/go-redis/v9"
)

type App struct {
	log        *slog.Logger
	restServer *rest.RestServer
	cfg        *config.Config
	db         *sql.DB
	client     *redis.Client
}

func NewApp(cfg *config.Config, log *slog.Logger) *App {
	db := mustLoadPostgres(cfg)
	client := mustLoadRedis(cfg)

	pos := postgres.NewPostgres(db)
	hasher := bcrypthash.NewBcryptHasher()
	redis := myredis.NewRedis(client, &cfg.RedisConf.TTL)
	idgen := uuidgen.NewUUIDGenerator()

	regUC := registration.NewRegUC(log, pos, hasher)
	logUC := login.NewLoginUC(log, pos, hasher, redis, idgen)

	handl := resthandler.NewRestHandler(log, &cfg.RestConf.CookieTTL, regUC, logUC)

	restServer := mustLoadHttpServer(cfg, log, handl)

	return &App{
		log:        log,
		restServer: restServer,
		cfg:        cfg,
		db:         db,
		client:     client,
	}
}

func (a *App) Run() {
	a.restServer.MustStart()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.RestConf.ShutdownTimeout*time.Second)
	defer cancel()

	a.restServer.Stop(ctx)

	a.db.Close()
	a.client.Close()
}
