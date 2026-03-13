package projectserviceclient

import (
	"context"
	"fmt"
	"log/slog"
	projectservicev1 "taskservice/proto/projectservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProjectServiceClient struct {
	log    *slog.Logger
	conn   *grpc.ClientConn
	client projectservicev1.ProjectServiceClient
}

func NewProjectServiceClient(log *slog.Logger, host string, port uint32) *ProjectServiceClient {
	const op = "projectserviceclient.NewProjectServiceClient"

	log.Info("create new project service client", slog.String("op", op))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("cannot create project service client: " + err.Error())
	}

	client := projectservicev1.NewProjectServiceClient(conn)

	return &ProjectServiceClient{
		log:    log,
		conn:   conn,
		client: client,
	}
}

func (p *ProjectServiceClient) GetOwnerId(ctx context.Context, projectId uint32) (uint32, error) {
	const op = "projectserviceclient.GetOwnerId"

	p.log.Info("sending a request to verify the owner", slog.String("op", op))
	in := &projectservicev1.GetOwnerIdRequest{
		ProjectId: projectId,
	}

	out, err := p.client.GetOwnerId(ctx, in)
	if err != nil {
		p.log.Warn("owner verification error", slog.String("op", op))
		return 0, err
	}

	p.log.Info("owner check successful", slog.String("op", op))
	return out.OwnerId, nil
}

func (p *ProjectServiceClient) Stop() {
	const op = "projectserviceclient.Stop"
	p.log.Info("disconnecting the connection", slog.String("op", op))
	p.conn.Close()
}
