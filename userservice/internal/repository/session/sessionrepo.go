package session

import "context"

type SessionRepo interface {
	Save(ctx context.Context, sessionId string, userId uint32) error
}
