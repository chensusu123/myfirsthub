package buff

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"
	"time"
)

func (b *Buff) GetOptionalMazeTempBuffListRQ_10435_10436(s *session.Session, req *MazeTempBuff.GetOptionalMazeTempBuffListRQ) (err error) {
	defer fkprometheus.InfoPMT("GetOptionalMazeTempBuffListRQ")()

	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	res := &MazeTempBuff.GetOptionalMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.Type = req.Type
	res.AreaId = req.AreaId
	res.AreaIndex = req.AreaIndex

	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, buffType, areaId := uint64(s.UID()), req.GetStageId(), req.GetLevel(), int32(req.GetType()), req.GetAreaId()
	if userId == 0 || stageId == 0 || level == 0 {
		logger.WarnWF("GetOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
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

	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetOptionalTempBuffList(logger, userId, stageId, level, int32(req.GetType()), areaId, req.GetAreaIndex())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.OptionalBuffInfo = optionalBuffInfo2PbOptionalBuffInfo(optionalBuffInfo)
	return
}

func optionalBuffInfo2PbOptionalBuffInfo(optionalBuffInfo *tempbuffservice.OptionalBuffInfo) *MazeTempBuff.OptionalBuffInfo {
	pb := &MazeTempBuff.OptionalBuffInfo{}
	pb.IsRefresh = proto.Int32(optionalBuffInfo.IsRefresh)
	pb.SelectBuffTime = proto.Int32(optionalBuffInfo.SelectBuffTime)
	pb.SelectBuffList = make([]*MazeTempBuff.MazeBuffInfo, 0, len(optionalBuffInfo.SelectBuffList))
	for _, buffInfo := range optionalBuffInfo.SelectBuffList {
		pb.SelectBuffList = append(pb.SelectBuffList, &MazeTempBuff.MazeBuffInfo{
			BuffId:      proto.Int32(buffInfo.BuffId),
			Value:       proto.Int64(buffInfo.Value),
			Name:        proto.String(buffInfo.Name),
			Decs:        proto.String(buffInfo.Decs),
			IsRecommend: proto.Int32(buffInfo.IsRecommend),
		})
	}
	pb.Cost = make([]*MazeCommon.MazeItem, 0, len(optionalBuffInfo.Cost))
	for _, cost := range optionalBuffInfo.Cost {
		pb.Cost = append(pb.Cost, &MazeCommon.MazeItem{
			ItemId: proto.Int32(cost.ItemId),
			Count:  proto.Int64(cost.Count),
		})
	}
	return pb
}
