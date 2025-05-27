package card

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/userriddlemonthlyredis"
	"maze_game_server/lib/log"
	"maze_game_server/pb/common/MazeCard"
	"time"

	"maze_game_server/lib/nano/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (*Card) GetMazeCardRQ_10430_10431(s *session.Session, req *MazeCard.GetMazeCardRQ) (err error) {
	defer fkprometheus.InfoPMT("GetMazeCardRQ")()
	start := time.Now()

	logger := log.Clone("Card", uint64(s.UID()), 0)
	res := &MazeCard.GetMazeCardRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetMazeCardRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId := uint64(s.UID())

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
