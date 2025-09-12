package buff

import (
	"context"
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (b *Buff) GetOptionalMazeTempBuffListRQ_10435_10436(s *session.Session, req *MazeTempBuff.GetOptionalMazeTempBuffListRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "GetOptionalMazeTempBuffListRQ start", zap.Any("req", req))
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
		logger.CtxInfo(ctx, "GetOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	userId, barrierId, level, buffType, areaId := uint64(s.UID()), req.GetStageId(), req.GetLevel(), int32(req.GetType()), req.GetAreaId()
	if userId == 0 || barrierId == 0 || level == 0 {
		logger.CtxError(ctx, "GetOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
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
	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetOptionalTempBuffList(ctx, userId, barrierId, level, int32(req.GetType()), areaId, req.GetAreaIndex(), req.GetAttrMask())
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

func OptionalBuffInfo2PbOptionalBuffInfo(ctx context.Context, userId uint64, barrierId int32, optionalBuffInfo *tempbuffservice.OptionalBuffInfo) (*MazeTempBuff.OptionalBuffInfo, error) {
	pb := &MazeTempBuff.OptionalBuffInfo{}
	pb.IsRefresh = proto.Int32(optionalBuffInfo.IsRefresh)
	pb.SelectBuffTime = proto.Int32(optionalBuffInfo.SelectBuffTime)
	pb.SelectBuffList = make([]*MazeTempBuff.MazeBuffInfo, 0, len(optionalBuffInfo.SelectBuffList))
	nowBuffGroupInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffGroupList(ctx, userId, barrierId)
	if err != nil {
		return nil, err
	}
	for _, buffInfo := range optionalBuffInfo.SelectBuffList {
		energyAffixCfg := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, buffInfo.BuffId)
		if energyAffixCfg == nil {
			return nil, fmt.Errorf("词条id%d缺少", buffInfo.BuffId)
		}

		var nowCount int32
		for _, buffGroupInfo := range nowBuffGroupInfo {
			if buffGroupInfo.GroupId == energyAffixCfg.Affix_group_id {
				nowCount = buffGroupInfo.Count
			}
		}

		pb.SelectBuffList = append(pb.SelectBuffList, &MazeTempBuff.MazeBuffInfo{
			BuffId:         proto.Int32(buffInfo.BuffId),
			Value:          proto.Int64(buffInfo.Value),
			Name:           proto.String(buffInfo.Name),
			Decs:           proto.String(buffInfo.Decs),
			IsRecommend:    proto.Int32(buffInfo.IsRecommend),
			BuffGroupCount: proto.Int32(nowCount),
		})
	}
	pb.Cost = make([]*MazeCommon.MazeItem, 0, len(optionalBuffInfo.Cost))
	for _, cost := range optionalBuffInfo.Cost {
		pb.Cost = append(pb.Cost, &MazeCommon.MazeItem{
			ItemId: proto.Int32(cost.ItemId),
			Count:  proto.Int64(cost.Count),
		})
	}
	return pb, nil
}
