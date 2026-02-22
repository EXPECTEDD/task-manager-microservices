package interceptor

import (
	"context"
	"io"
	"log/slog"
	projectservicev1 "projectservice/proto/projectservice"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTimeout(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	interceptor := TimeoutInterceptor(log, 1*time.Millisecond)

	handlSuccess := func(ctx context.Context, req any) (any, error) {
		return &projectservicev1.GetOwnerIdResponse{
			OwnerId: 1,
		}, nil
	}

	resp, err := interceptor(context.Background(), &projectservicev1.GetOwnerIdRequest{ProjectId: 1}, &grpc.UnaryServerInfo{}, handlSuccess)
	require.NoError(t, err)
	require.Equal(t, uint32(1), resp.(*projectservicev1.GetOwnerIdResponse).OwnerId)

	handlTimeout := func(ctx context.Context, req any) (any, error) {
		time.Sleep(2 * time.Millisecond)
		return &projectservicev1.GetOwnerIdResponse{
			OwnerId: 1,
		}, nil
	}

	resp, err = interceptor(context.Background(), &projectservicev1.GetOwnerIdRequest{ProjectId: 1}, &grpc.UnaryServerInfo{}, handlTimeout)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.DeadlineExceeded, s.Code())
	require.Equal(t, "request time out", s.Message())
	require.Nil(t, resp)
}
