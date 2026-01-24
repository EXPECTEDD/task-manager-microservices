package integration

import (
	"context"
	"testing"
	userservicev1 "userservice/proto/userservice"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuthenticate_Success_Integration(t *testing.T) {
	email, pass := registrateUser(t)
	sessionId := loginUser(t, email, pass)

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := userservicev1.NewUserServiceClient(conn)

	in := &userservicev1.GetIdBySessionRequest{
		SessionId: sessionId,
	}

	resp, err := client.GetIdBySession(context.Background(), in)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAuthenticate_SessionNotFound_Integration(t *testing.T) {
	sessionId := "sessiondId"

	conn, err := grpc.NewClient("localhost:44045", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := userservicev1.NewUserServiceClient(conn)

	in := &userservicev1.GetIdBySessionRequest{
		SessionId: sessionId,
	}

	out, err := client.GetIdBySession(context.Background(), in)
	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
	require.Equal(t, "session not found", s.Message())
	require.Nil(t, out)
}
