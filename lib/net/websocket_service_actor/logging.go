package websocket_service_actor

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
	slogzap "github.com/samber/slog-zap/v2"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func zapAdapterLogging(system *actor.ActorSystem) *slog.Logger {
	zapLogger := fklog.GetAppZapLogger()
	if zapLogger == nil {
		zapLogger, _ = zap.NewProduction()
	}

	logger := slog.New(slogzap.Option{Level: slog.LevelDebug, Logger: zapLogger, AddSource: true}.NewZapHandler())
	return logger.
		With("lib", "Proto.Actor").
		With("system", system.ID)
}
