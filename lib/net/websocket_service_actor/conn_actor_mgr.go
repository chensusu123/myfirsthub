package websocket_service_actor

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
)

type ConnActorMgr struct {
	userToConnActor     map[uint64]*actor.PID
	sesstionToConnActor map[uint64]*actor.PID
	system              *actor.ActorSystem
	conns               *actor.PIDSet
}

func NewConnActorMgr(system *actor.ActorSystem) *ConnActorMgr {
	return &ConnActorMgr{
		system:              system,
		conns:               actor.NewPIDSet(),
		userToConnActor:     make(map[uint64]*actor.PID),
		sesstionToConnActor: make(map[uint64]*actor.PID),
	}
}

func (a *ConnActorMgr) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *ClientLogin:
		if _, ok := a.userToConnActor[msg.UserId]; ok {
			ctx.Logger().Warn("user already login", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
			return
		}
		ctx.Logger().Info("user login", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
		a.userToConnActor[msg.UserId] = msg.wPID
		ctx.Send(msg.wPID, &msg.UserId)
		// a.sesstionToConnActor[msg.SessionId] = msg.wPID
		a.conns.Add(msg.wPID)
	case *ClientLogout:
		ctx.Logger().Info("user logout", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
		pid, ok := a.userToConnActor[msg.UserId]
		if !ok {
			ctx.Logger().Warn("user not login", slog.Int64("userId", int64(msg.UserId)))
			return
		}
		if pid != msg.wPID {
			ctx.Logger().Warn("user not is self", slog.Int64("userId", int64(msg.UserId)))
			return
		}
		ctx.Logger().Info("user logout remove", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
		delete(a.userToConnActor, msg.UserId)
		// delete(a.sesstionToConnActor, msg.SessionId)
		a.conns.Remove(pid)
	case *ClientKick:
		ctx.Logger().Info("user kick", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
		pid, ok := a.userToConnActor[msg.UserId]
		if !ok {
			ctx.Logger().Warn("user not login", slog.Int64("userId", int64(msg.UserId)))
			return
		}
		if pid == msg.wPID {
			ctx.Logger().Warn("user is self", slog.Int64("userId", int64(msg.UserId)))
			return
		}
		ctx.Logger().Info("user kick old user ", slog.Int64("userId", int64(msg.UserId)), slog.Uint64("sessionId", msg.SessionId))
		delete(a.userToConnActor, msg.UserId)
		a.conns.Remove(pid)
	case *SendDataMsg:
		if pid, ok := a.userToConnActor[msg.UserId]; ok {
			a.system.Root.Send(pid, msg)
			return
		}
		ctx.Logger().Warn("SendDataMsg user not login", slog.Int64("userId", int64(msg.UserId)))
	}
}
