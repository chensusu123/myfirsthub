package tempbuffservice

import (
	"context"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/excel/dollmappuzzlenewcfgex"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/model/tempbuffmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 进入关卡前检查关卡的buff情况，因为可能会有清除部分buff的情况
func (s *service) CheckTempBuff(ctx context.Context, userId uint64, barrierId int32, stage int32) (*tempbuffmodel.TempBuffInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	tempBuff, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "checkTempBuff GetMazeTempBuff", zap.Error(err))
		return nil, err
	}
	if tempBuff == nil || len(tempBuff.SelectedBuff) == 0 {
		logger.CtxInfo(ctx, "checkTempBuff not need delete buff")
		return nil, nil
	}
	// 已选择的buff不是0，就需要检查了
	//passArea, err := passareamodel.NewPassAreaModel(logger, userId, barrierId)
	//if err != nil {
	//	logger.ErrorWF("checkTempBuff GetBarrierPassArea fail", zap.Error(err))
	//	return nil, err
	//}
	passArea := dollmappuzzlenewcfgex.GetPassAreaInfos(barrierId, stage)
	deleteBuffIds := make([]int32, 0)
	j := 0
	for _, temp := range tempBuff.SelectedBuff {
		exist := false
		for _, i := range passArea {
			if temp.AreaId == i.AreaId && temp.AreaIndex == i.AreaIndex {
				exist = true
			}
		}
		if exist {
			tempBuff.SelectedBuff[j] = temp
			j++
		} else {
			deleteBuffIds = append(deleteBuffIds, temp.BuffId)
		}
	}
	tempBuff.SelectedBuff = tempBuff.SelectedBuff[:j]
	selectBuffCount := len(tempBuff.SelectedBuff)
	if len(deleteBuffIds) == 0 {
		logger.CtxInfo(ctx, "checkTempBuff deleteBuffIds==0 not need delete buff")
		return tempBuff, nil
	}

	// 有被清除掉的buff，那需要更新buff
	tempBuff.BuffSequence = &tempbuffmodel.BuffSequence{
		Level: int32(selectBuffCount) + 1,
	}
	var totalMap map[int32]int64
	totalMap, tempBuff.TotalBuff = s.GetTotalBuff(ctx, tempBuff.SelectedBuff)
	// 更新buff信息
	err = tempBuff.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "checkTempBuff SetMazeTempBuff failed", zap.Any("info", tempBuff), zap.Error(err))
		return nil, err
	}

	// 由于检查buff是在进入关卡时检查，那不需要重新计算buff的技能，因为进入关卡本身会计算
	// 推送buff变化信息
	msg := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
		UserId:  userId,
		StageId: barrierId,
		ChgType: 2,
		ChgDesc: "清除部分buff",
	}

	// 计算buff变化
	attrMap := make(map[int32]int64)
	for _, buffId := range deleteBuffIds {
		config := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, buffId)
		if config != nil {
			for id, value := range config.Add_attr {
				attrMap[id] += value
			}
		}
	}
	chgAttrs := make([]*mazetempbuffchgmsg.AttrChgInfo, 0, len(attrMap))
	for id, value := range attrMap {
		chgAttrs = append(chgAttrs, &mazetempbuffchgmsg.AttrChgInfo{
			AttrId: id,
			OldVal: totalMap[id] + value,
			CurVal: totalMap[id],
		})
	}

	msg.ChgAttrs = chgAttrs
	_ = mazetempbuffchgmsg.PushTempBuffChangeMsg(ctx, msg)

	// 同步到buff中心
	s.TempBuffChangeSync(ctx, logger, userId, tempBuff)

	logger.CtxInfo(ctx, "checkTempBuff delete buff success ", zap.Any("info", tempBuff), zap.Any("deleteBuffIds", deleteBuffIds))

	return tempBuff, nil
}
