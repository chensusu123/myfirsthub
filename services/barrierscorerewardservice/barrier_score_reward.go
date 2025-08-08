package barrierscorerewardservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/barrierscorerewardmodel"
)

// 保存积分装备奖励信息
func (s *service) SaveBarrierScoreRewardEquip(logger fklog.FKLogI, userId uint64, barrier int32, equipList map[int32]int32) (err error) {
	err = s.saveBarrierScoreReward(logger, userId, barrier, equipList, nil)
	if err != nil {
		return err
	}
	return nil
}

// 保存积分道具奖励信息
func (s *service) SaveBarrierScoreRewardItem(logger fklog.FKLogI, userId uint64, barrier int32, itemList map[int32]int32) (err error) {
	err = s.saveBarrierScoreReward(logger, userId, barrier, nil, itemList)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) saveBarrierScoreReward(logger fklog.FKLogI, userId uint64, barrier int32, equipList map[int32]int32, itemList map[int32]int32) (err error) {
	model, err := barrierscorerewardmodel.NewBarrierScoreRewardModel(logger, userId, barrier)
	if err != nil {
		logger.ErrorWF("SaveBarrierScoreReward NewBarrierScoreRewardModel err", zap.Error(err))
		return fmt.Errorf("获取用户关卡积分奖励信息失败")
	}
	if len(equipList) != 0 {
		for k, v := range equipList {
			exist := false
			for _, j := range model.EquipList {
				if k == j.EquipId {
					j.Count += v
					exist = true
				}
			}
			if !exist {
				model.EquipList = append(model.EquipList, &barrierscorerewardmodel.EquipRewardInfo{
					EquipId: k,
					Count:   v,
				})
			}
		}
	}
	if len(itemList) != 0 {
		for k, v := range itemList {
			exist := false
			for _, j := range model.ItemList {
				if k == j.ItemId {
					j.Count += v
					exist = true
				}
			}
			if !exist {
				model.ItemList = append(model.ItemList, &barrierscorerewardmodel.ItemRewardInfo{
					ItemId: k,
					Count:  v,
				})
			}
		}
	}

	err = model.Save(logger, userId, barrier)
	if err != nil {
		logger.ErrorWF("SaveBarrierScoreReward Save err", zap.Error(err))
		return fmt.Errorf("保存关卡奖励失败")
	}

	return nil
}
