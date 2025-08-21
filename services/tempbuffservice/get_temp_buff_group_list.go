package tempbuffservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/excel/mazeenergyaffixfrontv8config"
	"maze_game_server/excel/mazeenergyaffixlvv8config"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) GetTempBuffGroupList(ctx context.Context, userId uint64, barrier int32) ([]*GroupInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrier)
	if err != nil {
		logger.ErrorWF("GetMazeTempBuffListRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	return s.getGroupList(logger, buffInfo)
}

type GroupInfo struct {
	BuffId  int32
	Count   int32
	GroupId int32
}

func (s *service) getGroupList(logger fklog.FKLogI, buffModel *tempbuffmodel.TempBuffInfoModel) ([]*GroupInfo, error) {
	groupCount := make(map[int32]int32)
	groupFirstAffixList := make([]*GroupInfo, 0)
	for _, i := range buffModel.SelectedBuff {
		buffConfig := mazeenergyaffixlvv8config.GetAffixConfig(i.BuffId)
		if buffConfig == nil {
			logger.WarnWF("GetTempBuffGroupList buffConfig is nil", zap.Int32("buffId", i.BuffId))
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
			frontConfig := mazeenergyaffixfrontv8config.GetMazeEnergyAffixFrontConfig(j)
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

	return groupFirstAffixList, nil
}
