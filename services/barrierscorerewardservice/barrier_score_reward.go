package barrierscorerewardservice

import (
	"context"
	"fmt"
	"maze_game_server/model/barrierscorerewardmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 保存积分装备奖励信息
func (s *service) SaveBarrierScoreRewardEquip(ctx context.Context, userId uint64, barrier int32, equipList map[int32]int32) (err error) {
	err = s.SaveBarrierScoreReward(ctx, userId, barrier, equipList, nil)
	if err != nil {
		return err
	}
	return nil
}

// 保存积分道具奖励信息
func (s *service) SaveBarrierScoreRewardItem(ctx context.Context, userId uint64, barrier int32, itemList map[int32]int64) (err error) {
	err = s.SaveBarrierScoreReward(ctx, userId, barrier, nil, itemList)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SaveBarrierScoreReward(ctx context.Context, userId uint64, barrier int32, equipList map[int32]int32, itemList map[int32]int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barrierscorerewardmodel.NewBarrierScoreRewardModel(ctx, userId, barrier)
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

	err = model.Save(ctx, userId, barrier)
	if err != nil {
		logger.ErrorWF("SaveBarrierScoreReward Save err", zap.Error(err))
		return fmt.Errorf("保存关卡奖励失败")
	}

	return nil
}

func (s *service) GetBarrierScoreReward(ctx context.Context, userId uint64, barrier int32) (equipList map[int32]int32, itemList map[int32]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barrierscorerewardmodel.NewBarrierScoreRewardModel(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierScoreReward NewBarrierScoreRewardModel err", zap.Error(err))
		return nil, nil, fmt.Errorf("获取用户关卡积分奖励信息失败")
	}
	equipList = make(map[int32]int32)
	itemList = make(map[int32]int64)
	for _, v := range model.EquipList {
		equipList[v.EquipId] = v.Count
	}
	for _, v := range model.ItemList {
		itemList[v.ItemId] = v.Count
	}
	return equipList, itemList, nil
}

// 删除关卡积分奖励记录
func (s *service) DelBarrierScoreRewardItem(ctx context.Context, userId uint64, barrier int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	model := &barrierscorerewardmodel.BarrierScoreRewardModel{}
	err = model.Del(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "DelBarrierScoreRewardItem del err", zap.Error(err), zap.Int32("barrier", barrier))
		return err
	}
	logger.CtxInfo(ctx, "DelBarrierScoreRewardItem success", zap.Int32("barrier", barrier))
	return nil
}
