package resthandler

import "log/slog"

type RestHandler struct {
	log *slog.Logger
}

func NewRestHandler(log *slog.Logger) *RestHandler {
	return &RestHandler{
		log: log,
	}
}
