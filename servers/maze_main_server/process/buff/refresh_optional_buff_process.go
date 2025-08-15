package buff

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (b *Buff) RefreshOptionalMazeTempBuffListRQ_10439_10440(s *session.Session, req *MazeTempBuff.RefreshOptionalMazeTempBuffListRQ) (err error) {
	defer fkprometheus.InfoPMT("RefreshOptionalMazeTempBuffListRQ")()
	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	logger.InfoWF("RefreshOptionalMazeTempBuffListRQ start", zap.Any("req", req))
	res := &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.AreaId = req.AreaId
	defer func() {
		err = s.Response(res)
		logger.InfoWF("RefreshOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, cost, areaId, buffType := uint64(s.UID()), req.GetStageId(), req.GetLevel(), req.GetCost(), req.GetAreaId(), int32(req.GetType())
	if userId == 0 || stageId == 0 || level == 0 {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}
	if buffType != int32(MazeTempBuff.Type_UP_LEVEL) && buffType != int32(MazeTempBuff.Type_USE_ITEM) {
		logger.ErrorWF("GetOptionalMazeTempBuffListRQ buffType args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff类型参数错误")
		return nil
	}
	if buffType == int32(MazeTempBuff.Type_UP_LEVEL) && areaId == 0 {
		logger.WarnWF("GetOptionalMazeTempBuffListRQ areaId error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("areaId参数错误")
		return nil
	}

	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.RefreshOptionalMazeTempBuffList(logger, userId, stageId, level, areaId, req.GetAttrMask(), cost)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.OptionalBuffInfo = OptionalBuffInfo2PbOptionalBuffInfo(optionalBuffInfo)
	return
}
