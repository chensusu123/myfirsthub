package card

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/userriddlemonthlyredis"
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCard"
	"go.uber.org/zap"
)

func GetMazeCardRQ(logger fknet.TCPContext, shardingID uint64, request, response proto.Message) error {
	defer fkprometheus.InfoPMT("GetMazeCardRQ")()
	start := time.Now()
	req := request.(*MazeCard.GetMazeCardRQ)
	res := response.(*MazeCard.GetMazeCardRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	defer func() {
		logger.InfoWF("GetMazeCardRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId := shardingID
	if userId == 0 {
		logger.WarnWF("GetMazeCardRQ args error")
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}

	// todo 检查用户是不是小程序用户

	// 获取月卡信息
	expirationTime, err := userriddlemonthlyredis.GetMazeCardExpirationTime(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMazeCardRQ GetMazeCardExpirationTime failed", zap.Error(err))
		return nil
	}

	// 检查月卡信息
	err = UpdateCard(logger, userId, expirationTime)
	if err != nil {
		logger.ErrorWF("GetMazeCardRQ UpdateCard failed", zap.Error(err))
	}

	if expirationTime < time.Now().Unix() {
		return nil
	}

	res.HaveCard = proto.Int32(1)
	res.ExpirationTime = proto.Int64(expirationTime)
	return nil
}

func UpdateCard(logger fklog.FKLogI, userId uint64, expirationTime int64) error {
	info, err := mazebuffinforedis.GetMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcMonthCard)
	if err != nil {
		logger.ErrorWF("UpdateCard GetMazeBuffBySrc failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	if info != nil {
		if expirationTime >= time.Now().Unix() {
			return nil
		}
		return DeleteMazeCard(logger, userId)
	}

	if expirationTime < time.Now().Unix() {
		return nil
	}

	return AddMazeCard(logger, userId, expirationTime)
}
