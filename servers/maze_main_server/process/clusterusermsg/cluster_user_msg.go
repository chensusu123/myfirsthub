package clusterusermsg

import (
	"context"
	"strconv"

	"maze_game_server/common/function/clusterpaket"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/simpleclient/simplenatsconsumer"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/natsconsumer"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/natsconsumer/natsmsgoption"
	"go.uber.org/zap"
)

const natsService = "maze.usermsg.route.nats"

func Register() {
	fkserver.AddBusiness(natsconsumer.NewNatsConsumerSvc("maze.user.msg.consumer", natsService,
		natsconsumer.WithMsgProcessor(StandaloneConsumerBytesConsumer()),
	))
}

func StandaloneConsumerBytesConsumer() *simplenatsconsumer.SimpleConsumerProcessor {
	sc := &StandaloneConsumer{}
	return simplenatsconsumer.New(
		simplenatsconsumer.WithSubject("maze.user.cluster.msg"),
		simplenatsconsumer.WithProcessorFunc(sc.Processor),
		simplenatsconsumer.WithSectionSubject(),
		simplenatsconsumer.WithStandaloneConsumer(),
	)
}

type StandaloneConsumer struct{}

func (s *StandaloneConsumer) Processor(ctx context.Context, obj any, opts ...natsmsgoption.NatsMsgOption) (any, error) {
	data, _ := obj.([]byte)
	logger := fklog.ContextAppLogger(ctx)
	cfg := natsmsgoption.NewNatsMsgOpts(opts...)

	packetType, data, err := clusterpaket.SplitClusterPacket(data)
	if err != nil {
		logger.CtxWarn(ctx, "clusterusermsg Processor SplitClusterPacket failed",
			zap.String("subject", cfg.Subject),
			zap.String("PrimaryKey", cfg.PrimaryKey),
			zap.Any("err", err),
		)
		return nil, nil
	}
	if packetType == 0 {
		logger.CtxWarn(ctx, "clusterusermsg Processor PacketType is invalid",
			zap.String("subject", cfg.Subject),
		)
		return nil, nil
	}

	userIDUint, err := strconv.ParseUint(cfg.PrimaryKey, 10, 64)
	if err != nil {
		logger.CtxWarn(ctx, "clusterusermsg Processor ParseUint failed",
			zap.String("subject", cfg.Subject),
			zap.String("PrimaryKey", cfg.PrimaryKey),
			zap.Any("err", err),
		)
		return nil, nil
	}
	if userIDUint == 0 {
		logger.CtxWarn(ctx, "clusterusermsg Processor PrimaryKey is invalid",
			zap.String("subject", cfg.Subject),
			zap.String("PrimaryKey", cfg.PrimaryKey),
		)
		return nil, nil
	}
	err = online.PushBytes(ctx, userIDUint, packetType, data)
	if err != nil {
		logger.CtxWarn(ctx, "clusterusermsg Processor PushBytes failed",
			zap.String("subject", cfg.Subject),
			zap.String("PrimaryKey", cfg.PrimaryKey),
			zap.Any("err", err),
		)
		return nil, nil
	}
	logger.CtxInfo(ctx, "clusterusermsg Success",
		zap.Uint64("userIDUint", userIDUint),
		zap.Uint16("packetType", packetType),
		zap.Int("dataLen", len(data)),
	)
	return nil, nil
}
