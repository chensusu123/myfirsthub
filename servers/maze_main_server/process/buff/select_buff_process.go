package buff

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"
	"time"
)

func (b *Buff) SelectMazeTempBuffRQ_10437_10438(s *session.Session, req *MazeTempBuff.SelectMazeTempBuffRQ) (err error) {
	defer fkprometheus.InfoPMT("SelectMazeTempBuffRQ")()

	start := time.Now()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.InfoWF("SelectMazeTempBuffRQ start", zap.Any("req", req))
	res := &MazeTempBuff.SelectMazeTempBuffRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.Type = req.Type
	defer func() {
		err = s.Response(res)
		logger.InfoWF("SelectMazeTempBuffRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, buffId, buffType := uint64(s.UID()), req.GetStageId(), req.GetLevel(), req.GetBuffId(), int32(req.GetType())
	if userId == 0 || stageId == 0 || level == 0 || buffId == 0 {
		logger.WarnWF("SelectMazeTempBuffRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}

	if buffType != int32(MazeTempBuff.Type_UP_LEVEL) && buffType != int32(MazeTempBuff.Type_USE_ITEM) {
		logger.ErrorWF("SelectMazeTempBuffRQ buffType args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff类型参数错误")
		return nil
	}

	buffList, err := tempbuffservice.GlobalTempBuffService.SelectMazeTempBuff(ctx, userId, stageId, level, buffId, buffType)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.BuffList = buffInfo2MazeBuffInfo(buffList)
	return
}
