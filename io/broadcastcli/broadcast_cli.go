package broadcastcli

import (
	"context"
	"encoding/binary"
	"fmt"

	"maze_game_server/common/function/clusterpaket"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/simpleclient/simplenatsproducer"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/nanotrace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var gNatsproducer database.NatsProducerI

func init() {
	natsService := "maze.usermsg.route.nats"
	natsproducer := simplenatsproducer.New(natsService, "BroadcastProducer")
	gNatsproducer = natsproducer
	serverdepend.RegisterDepend(natsproducer)
}

func Broadcast(ctx context.Context, broadcastID uint64, packetType uint16, v interface{}) (err error) {
	span := nanotrace.NewSimpleTrace("BroadcastMsg")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	cSpan := nanotrace.SpanFromContext(ctx)
	cSpan.SetAttributes(
		attribute.Int64("broadcast.id", int64(broadcastID)),
		attribute.Int("packet.id", int(packetType)),
		attribute.String("nats.subject", "maze.broadcast.cluster.msg.>"),
	)
	data, err := clusterpaket.MakeClusterPacket(packetType, v)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("maze.broadcast.cluster.msg.%s.%d", appconfig.GlobalConfig().Global.SectionID, broadcastID)
	err = gNatsproducer.Publish(ctx, subject, data)
	fklog.ContextAppLogger(ctx).CtxInfo(ctx, "broadcastcli publish to nats",
		zap.String("subject", subject),
		zap.Uint64("broadcastID", broadcastID),
		zap.Uint16("packetType", packetType),
		zap.Error(err))
	return err
}

// Deprecated: 仅仅供测试的时候使用
func BroadcastTest(ctx context.Context, broadcastID uint64, packetType uint16, data []byte) (err error) {
	span := nanotrace.NewSimpleTrace("BroadcastMsg")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	cSpan := nanotrace.SpanFromContext(ctx)
	cSpan.SetAttributes(
		attribute.Int64("broadcast.id", int64(broadcastID)),
		attribute.Int("packet.id", int(packetType)),
		attribute.String("nats.subject", "maze.user.msg.*"),
	)
	ret := make([]byte, 2+len(data))
	binary.LittleEndian.PutUint16(ret[:2], packetType)
	copy(ret[2:], data)
	subject := fmt.Sprintf("maze.broadcast.msg.%d", broadcastID)
	err = gNatsproducer.Publish(ctx, subject, ret)
	fklog.ContextAppLogger(ctx).CtxInfo(ctx, "broadcastcli publish to nats",
		zap.String("subject", subject),
		zap.Uint64("broadcastID", broadcastID),
		zap.Uint16("packetType", packetType),
		zap.Error(err))
	return err
}
