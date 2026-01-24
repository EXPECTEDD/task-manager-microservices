package interceptor

import (
	"context"
	"io"
	"log/slog"
	"testing"
	userservicev1 "userservice/proto/userservice"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlSuccess struct {
}

func (h *HandlSuccess) GetIdBySession(ctx context.Context, req *userservicev1.GetIdBySessionRequest) (*userservicev1.GetIdBySessionResponse, error) {
	return &userservicev1.GetIdBySessionResponse{
		UserId: 1,
	}, nil
}

type HandlPanic struct {
}

func (h *HandlPanic) GetIdBySession(ctx context.Context, req *userservicev1.GetIdBySessionRequest) (*userservicev1.GetIdBySessionResponse, error) {
	panic("panic")
}

func TestRecover(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	interceptor := RecoverInterceptor(log)

	handlSuccess := func(ctx context.Context, req any) (any, error) {
		return &userservicev1.GetIdBySessionResponse{
			UserId: 1,
		}, nil
	}

	resp, err := interceptor(context.Background(), &userservicev1.GetIdBySessionRequest{SessionId: "sessionId"}, &grpc.UnaryServerInfo{}, handlSuccess)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), resp.(*userservicev1.GetIdBySessionResponse).UserId)

	handlPanic := func(ctx context.Context, req any) (any, error) {
		panic("panic")
	}

	resp, err = interceptor(context.Background(), &userservicev1.GetIdBySessionRequest{SessionId: "sessionId"}, &grpc.UnaryServerInfo{}, handlPanic)

	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, s.Code())
	assert.Equal(t, "internal server error", s.Message())
	assert.Nil(t, resp)
}
