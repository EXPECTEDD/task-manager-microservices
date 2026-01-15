package rest

import (
	"context"
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
	r.log.Info("starting the server", slog.String("op", op), slog.String("port", r.serv.Addr))

	if err := r.serv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			panic("server bad start")
		}
	}
}

func (r *RestServer) Stop(ctx context.Context) {
	const op = "rest.Stop"
	r.log.Info("start server shutdown")
	r.serv.Shutdown(ctx)
	r.log.Info("server stopped")
}
