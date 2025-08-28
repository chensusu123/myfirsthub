package tempbuffservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeEnergyAffixFrontV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/excel/mazeconfigv8"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) GetTempBuffGroupList(ctx context.Context, userId uint64, barrier int32) ([]*GroupInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "GetMazeTempBuffListRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	return s.getGroupList(ctx, logger, buffInfo)
}

type GroupInfo struct {
	BuffId  int32
	Count   int32
	GroupId int32
}

func (s *service) getGroupList(ctx context.Context, logger fklog.FKLogI, buffModel *tempbuffmodel.TempBuffInfoModel) ([]*GroupInfo, error) {
	groupCount := make(map[int32]int32)
	groupFirstAffixList := make([]*GroupInfo, 0)
	for _, i := range buffModel.SelectedBuff {
		buffConfig := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, i.BuffId)
		if buffConfig == nil {
			logger.CtxError(ctx, "GetTempBuffGroupList buffConfig is nil", zap.Int32("buffId", i.BuffId))
			return nil, fmt.Errorf("能力词条不存在")
		}
		groupCount[buffConfig.Affix_group_id] += 1
		for _, j := range buffConfig.Font_affix_condition {
			if j == 0 {
				// 找到词条组的第一个词条
				groupFirstAffixList = append(groupFirstAffixList, &GroupInfo{
					BuffId:  i.BuffId,
					GroupId: buffConfig.Affix_group_id,
				})
				continue
			}
			frontConfig := GMazeEnergyAffixFrontV8Cfg.GetWithCtx(ctx, j)
			if frontConfig.Extra_affix_group_id == 0 {
				continue
			}
			// 有额外词条组的
			groupCount[frontConfig.Extra_affix_group_id] += 1
		}
	}

	for _, i := range groupFirstAffixList {
		i.Count = groupCount[i.GroupId]
	}
	var specialBuffGroup *GroupInfo
	res := make([]*GroupInfo, 0, len(groupFirstAffixList))
	for _, i := range groupFirstAffixList {
		if i.GroupId == int32(mazeconfigv8.GetSpecialBuffGroupId(ctx)) {
			specialBuffGroup = i
		} else {
			res = append(res, i)
		}
	}

	difference := int(mazeconfigv8.GetMaxBuffGroupCount(ctx)) - len(res)
	if difference > 0 {
		for i := 0; i < difference; i++ {
			res = append(res, &GroupInfo{})
		}
	}
	if specialBuffGroup == nil {
		specialBuffGroup = &GroupInfo{}
	}
	res = append(res, specialBuffGroup)

	return res, nil
}
