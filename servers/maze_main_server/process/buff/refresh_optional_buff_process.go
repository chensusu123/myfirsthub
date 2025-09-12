package buff

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (b *Buff) RefreshOptionalMazeTempBuffListRQ_10439_10440(s *session.Session, req *MazeTempBuff.RefreshOptionalMazeTempBuffListRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "RefreshOptionalMazeTempBuffListRQ start", zap.Any("req", req))
	res := &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.AreaId = req.AreaId
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "RefreshOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	userId, barrierId, level, cost, areaId, buffType := uint64(s.UID()), req.GetStageId(), req.GetLevel(), req.GetCost(), req.GetAreaId(), int32(req.GetType())
	if userId == 0 || barrierId == 0 || level == 0 {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}
	if buffType != int32(MazeTempBuff.Type_UP_LEVEL) && buffType != int32(MazeTempBuff.Type_USE_ITEM) {
		logger.CtxError(ctx, "GetOptionalMazeTempBuffListRQ buffType args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff类型参数错误")
		return nil
	}
	if buffType == int32(MazeTempBuff.Type_UP_LEVEL) && areaId == 0 {
		logger.CtxError(ctx, "GetOptionalMazeTempBuffListRQ areaId error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("areaId参数错误")
		return nil
	}

	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.RefreshOptionalMazeTempBuffList(ctx, userId, barrierId, level, areaId, req.GetAttrMask(), cost)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.OptionalBuffInfo, err = OptionalBuffInfo2PbOptionalBuffInfo(ctx, userId, barrierId, optionalBuffInfo)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	return
}
