package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type RestServer struct {
	log  *slog.Logger
	serv *http.Server
}

func NewRestServer(log *slog.Logger, serv *http.Server) *RestServer {
	return &RestServer{
		log:  log,
		serv: serv,
	}
}

func (r *RestServer) MustStart() {
	const op = "rest.MustStart"
	r.log.Info("starting http server", slog.String("op", op), slog.String("port", r.serv.Addr))

	if err := r.serv.ListenAndServe(); err != nil {
		if !errors.Is(http.ErrServerClosed, err) {
			panic("server bad start: " + err.Error())
		}
	}
}

func (r *RestServer) Stop(ctx context.Context) {
	const op = "rest.Stop"
	r.log.Info("start http server shutdown", slog.String("op", op))
	r.serv.Shutdown(ctx)
	r.log.Info("http server stopped", slog.String("op", op))
}
