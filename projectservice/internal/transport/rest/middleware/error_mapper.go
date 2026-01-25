package middleware

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcErrorToHttp(err error) (int, string) {
	s, ok := status.FromError(err)
	if !ok {
		return http.StatusBadGateway, "upstream error"
	}

	switch s.Code() {
	case codes.NotFound:
		return http.StatusNotFound, "not found"
	case codes.Internal:
		return http.StatusBadGateway, "upstream error"
	default:
		return http.StatusBadGateway, "upstream error"
	}
}
