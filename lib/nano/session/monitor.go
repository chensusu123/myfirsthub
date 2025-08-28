package session

import "context"

type Monitor interface {
	// OnCreate
	OnCreate(ctx context.Context, s *Session)

	// OnClose
	OnClose(ctx context.Context, s *Session, err error)
}
