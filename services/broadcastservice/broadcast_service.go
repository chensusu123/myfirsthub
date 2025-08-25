package broadcastservice

import (
	"context"
	"errors"
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

type BroadcastUsers interface {
	BroadcastUserList(ctx context.Context, broadcastID int64) ([]uint64, error)
}

var gStandaloneConsumer *StandaloneConsumer = &StandaloneConsumer{
	nameLen: len("maze.broadcast.msg."),
}

func Register(broadcastUsers BroadcastUsers) {
	fkserver.AddBusiness(natsconsumer.NewNatsConsumerSvc("maze.broadcast.msg.consumer", natsService,
		natsconsumer.WithMsgProcessor(StandaloneConsumerBytesConsumer()),
	))
	gStandaloneConsumer.broadcastUsers = broadcastUsers
}

func StandaloneConsumerBytesConsumer() *simplenatsconsumer.SimpleConsumerProcessor {
	return simplenatsconsumer.New(
		simplenatsconsumer.WithSubject("maze.broadcast.msg.*"),
		simplenatsconsumer.WithProcessorFunc(gStandaloneConsumer.Processor),
		simplenatsconsumer.WithStandaloneConsumer(),
	)
}

type StandaloneConsumer struct {
	nameLen        int
	broadcastUsers BroadcastUsers
}

func (s *StandaloneConsumer) Processor(ctx context.Context, obj any, opts ...natsmsgoption.NatsMsgOption) (any, error) {
	data, _ := obj.([]byte)
	logger := fklog.ContextAppLogger(ctx)
	cfg := natsmsgoption.NewNatsMsgOpts(opts...)
	if len(cfg.Subject) <= s.nameLen {
		logger.CtxWarn(ctx, "broadcastservice Processor Subject is invalid",
			zap.String("subject", cfg.Subject),
		)
		return nil, nil
	}
	userID := cfg.Subject[s.nameLen:]
	packetType, data, err := clusterpaket.SplitClusterPacket(data)
	if err != nil {
		logger.CtxWarn(ctx, "broadcastservice Processor SplitClusterPacket failed",
			zap.String("subject", cfg.Subject),
			zap.Any("err", err),
		)
		return nil, nil
	}
	if packetType == 0 {
		logger.CtxWarn(ctx, "broadcastservice Processor PacketType is invalid",
			zap.String("subject", cfg.Subject),
		)
		return nil, nil
	}

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		logger.CtxWarn(ctx, "broadcastservice Processor ParseUint failed",
			zap.String("subject", cfg.Subject),
			zap.Any("err", err),
		)
		return nil, nil
	}

	userList, err := s.broadcastUsers.BroadcastUserList(ctx, int64(userIDUint))
	if err != nil {
		logger.CtxWarn(ctx, "broadcastservice Processor BroadcastUserList failed",
			zap.String("subject", cfg.Subject),
			zap.Any("err", err),
		)
		return nil, nil
	}
	for _, user := range userList {
		err = online.PushBytes(ctx, user, packetType, data)
		if err != nil && !errors.Is(err, online.ErrSessionNotFound) {
			logger.CtxWarn(ctx, "broadcastservice Processor PushBytes failed",
				zap.String("subject", cfg.Subject),
				zap.Any("err", err),
			)
		}
	}
	logger.CtxInfo(ctx, "broadcastservice Success",
		zap.Uint64("userIDUint", userIDUint),
		zap.Uint64("userIDUint", userIDUint),
		zap.Uint16("packetType", packetType),
		zap.Int("dataLen", len(data)),
	)
	return nil, nil
}

type normalBroadcast struct{}

// BroadcastUserList 普通广播，广播所有在线用户
func NewNormalBroadcast() *normalBroadcast {
	return &normalBroadcast{}
}

func (n *normalBroadcast) BroadcastUserList(ctx context.Context, broadcastID int64) ([]uint64, error) {
	userList := online.GetOnlineUsers()

	return userList, nil
}
