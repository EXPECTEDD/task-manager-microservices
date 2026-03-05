package integration

import (
	"context"
	projectservicev1 "projectservice/proto/projectservice"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestGetOwnerId_Success_Integration(t *testing.T) {
	userId, email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)
	projectId := createProject(t, sessionId, "NewProj")

	conn, err := grpc.NewClient(projectServiceConn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := projectservicev1.NewProjectServiceClient(conn)

	req := &projectservicev1.GetOwnerIdRequest{ProjectId: projectId}

	resp, err := client.GetOwnerId(context.Background(), req)
	sErr, _ := status.FromError(err)
	require.Equal(t, codes.OK, sErr.Code())
	require.Equal(t, userId, resp.OwnerId)
}

func TestGetOwnerId_ProjectNotFound_Integration(t *testing.T) {
	_, email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)
	projectId := createProject(t, sessionId, "NewProj")

	conn, err := grpc.NewClient(projectServiceConn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := projectservicev1.NewProjectServiceClient(conn)

	req := &projectservicev1.GetOwnerIdRequest{ProjectId: projectId + 1}

	resp, err := client.GetOwnerId(context.Background(), req)
	sErr, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, sErr.Code())
	require.Equal(t, (*projectservicev1.GetOwnerIdResponse)(nil), resp)
}

func TestGetOwnerId_InvalidProjectId_Integration(t *testing.T) {
	_, email, pass := registrationUser(t)
	sessionId := loginUser(t, email, pass)
	_ = createProject(t, sessionId, "NewProj")

	conn, err := grpc.NewClient(projectServiceConn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := projectservicev1.NewProjectServiceClient(conn)

	req := &projectservicev1.GetOwnerIdRequest{ProjectId: 0}

	resp, err := client.GetOwnerId(context.Background(), req)
	sErr, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, sErr.Code())
	require.Equal(t, (*projectservicev1.GetOwnerIdResponse)(nil), resp)
}
