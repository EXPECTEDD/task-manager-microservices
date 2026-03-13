package userserviceclient

import (
	"context"
	"fmt"
	"log/slog"
	userservicev1 "taskservice/proto/userservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	log    *slog.Logger
	conn   *grpc.ClientConn
	client userservicev1.UserServiceClient
}

func NewUserServiceClient(log *slog.Logger, host string, port uint32) *UserServiceClient {
	const op = "userserviceclient.NewUserServiceClient"

	log.Info("create new user service client", slog.String("op", op))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("cannot create new user service client: " + err.Error())
	}

	client := userservicev1.NewUserServiceClient(conn)

	return &UserServiceClient{
		log:    log,
		conn:   conn,
		client: client,
	}
}

func (u *UserServiceClient) GetIdBySession(ctx context.Context, sessionId string) (uint32, error) {
	const op = "userserviceclient.GetIdBySession"

	u.log.Info("sending a request to verify a session", slog.String("op", op))
	in := &userservicev1.GetIdBySessionRequest{
		SessionId: sessionId,
	}

	res, err := u.client.GetIdBySession(ctx, in)
	if err != nil {
		u.log.Warn("session validity check error", slog.String("op", op))
		return 0, err
	}

	u.log.Info("session is valid", slog.String("op", op))
	return res.UserId, nil
}

func (u *UserServiceClient) Stop() {
	const op = "userserviceclient.Stop"
	u.log.Info("disconnecting the connection", slog.String("op", op))
	u.conn.Close()
}
