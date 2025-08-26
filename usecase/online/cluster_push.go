package online

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"maze_game_server/common/function/clusterpaket"
	"maze_game_server/lib/codec"
	"maze_game_server/lib/nano/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/client/natsproduceroption"
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

	natsproducer := simplenatsproducer.New(natsService, "NatsProducerDemo")
	gNatsproducer = natsproducer
	serverdepend.RegisterDepend(natsproducer)
}

// MakeNormalPushData generate a normal push data
func MakeNormalPushData(packetType uint16, payload interface{}, isBytes bool) *session.NormalPushData {
	return &session.NormalPushData{
		Mid:     codec.ToMessageID(uint32(time.Now().Unix()), 0, packetType),
		Payload: payload,
		IsBytes: isBytes,
	}
}

// ClusterPush push data to cluster
// 注意： 如果用户不在当前分片。 则会往其他分片广播，由其他分片发送给用户
func ClusterPush(ctx context.Context, userID uint64, packetType uint16, v interface{}) (err error) {
	s, found := monitor.online.Load(userID)
	if !found {
		span := nanotrace.NewSimpleTrace("UserMsgPushToOtherServer")
		ctx = span.Start(ctx)
		defer span.Finish(ctx)
		cSpan := nanotrace.SpanFromContext(ctx)
		subject := fmt.Sprintf("maze.user.cluster.msg.%s.%d", appconfig.GlobalConfig().Global.SectionID, userID)
		cSpan.SetAttributes(
			attribute.Int64("enduser.id", int64(userID)),
			attribute.Int("packet.id", int(packetType)),
			attribute.String("nats.subject", "maze.user.cluster.msg.>"),
		)
		data, err := clusterpaket.MakeClusterPacket(packetType, v)
		if err != nil {
			return err
		}

		err = gNatsproducer.Publish(ctx, subject, data, natsproduceroption.WithSkipSelfConsumer())
		fklog.ContextAppLogger(ctx).CtxInfo(ctx, "ClusterPush publish to nats",
			zap.String("subject", subject),
			zap.Uint64("userID", userID),
			zap.Uint16("packetType", packetType),
			zap.Error(err))
		return err
	}
	return s.(*session.Session).ResponseMID(ctx, codec.ToMessageID(uint32(time.Now().Unix()), 0, packetType), v)
}

// PushBytes push bytes data to user
func PushBytes(ctx context.Context, userID uint64, packetType uint16, data []byte) (err error) {
	span := nanotrace.NewSimpleTrace("PushBytesToUser")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	cSpan := nanotrace.SpanFromContext(ctx)
	cSpan.SetAttributes(
		attribute.Int64("enduser.id", int64(userID)),
		attribute.Int("packet.id", int(packetType)),
	)
	s, found := monitor.online.Load(userID)
	if !found {
		cSpan.AddEvent("PushBytes not found")
		return ErrSessionNotFound
	}

	return s.(*session.Session).Push(ctx, "", MakeNormalPushData(packetType, data, true))
}

// Deprecated: 仅仅供测试的时候使用
// PushToClusterTest push data to cluster
func PushToClusterTest(ctx context.Context, userID uint64, packetType uint16, data []byte) (err error) {
	span := nanotrace.NewSimpleTrace("PushToClusterTest")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	ret := make([]byte, 2+len(data))
	binary.LittleEndian.PutUint16(ret[:2], packetType)
	copy(ret[2:], data)
	fklog.ContextAppLogger(ctx).CtxInfo(ctx, "PushToClusterTest",
		zap.Uint64("userID", userID),
		zap.Uint16("packetType", packetType))
	return gNatsproducer.Publish(ctx, fmt.Sprintf("maze.user.msg.%d", userID), ret)
}
