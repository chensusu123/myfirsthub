package game

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/servers/maze_main_server/process/buff"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/itemservice"
	"maze_game_server/services/tempbuffservice"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// UseGoldCoinPile 使用金币堆，转换为金币
func UseGoldCoinPile(logger fklog.FKLogI, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (items []*MazeCommon.MazeItem, err error) {
	items = make([]*MazeCommon.MazeItem, 0)
	cfg := GMazeConfigV8Cfg.Get(constdef.GoldPileShow2RealCfgId)
	if cfg == nil {
		err = errors.New("金币堆兑换金币配置错误")
		logger.ErrorWF("UseGoldCoinPile GMazeConfigV8Cfg.Get fail", zap.Error(err), zap.Int32("cfgId", constdef.GoldPileShow2RealCfgId))
		return nil, err
	}
	goldCoinId, ok := cfg.Value_map[itemId]
	if !ok {
		err = errors.New("金币堆无映射金币配置")
		logger.ErrorWF("UseGoldCoinPile no target gold coin found", zap.Error(err), zap.Int32("itemId", itemId), zap.Any("Value_map", cfg.Value_map))
		return nil, err
	}
	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	// Add item
	items = append(items, &MazeCommon.MazeItem{
		ItemId: proto.Int32(int32(goldCoinId)),
		Count:  proto.Int64(int64(barrierCfg.Item1_nums_per_pile) * count),
	})
	return
}

// UseQianghuashiPile 使用强化石堆，转换为强化石
func UseQianghuashiPile(logger fklog.FKLogI, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (items []*MazeCommon.MazeItem, err error) {
	items = make([]*MazeCommon.MazeItem, 0)
	cfg := GMazeConfigV8Cfg.Get(constdef.StrengthenStonePileShow2RealCfgId)
	if cfg == nil {
		err = errors.New("强化石堆兑换强化石配置错误")
		logger.ErrorWF("UseQianghuashiPile GMazeConfigV8Cfg.Get fail", zap.Error(err), zap.Int32("cfgId", constdef.StrengthenStonePileShow2RealCfgId))
		return nil, err
	}
	goldCoinId, ok := cfg.Value_map[itemId]
	if !ok {
		err = errors.New("强化石堆无映射强化石配置")
		logger.ErrorWF("UseQianghuashiPile no target qianghuashi found", zap.Error(err), zap.Int32("itemId", itemId), zap.Any("Value_map", cfg.Value_map))
		return nil, err
	}
	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	// Add item
	items = append(items, &MazeCommon.MazeItem{
		ItemId: proto.Int32(int32(goldCoinId)),
		Count:  proto.Int64(int64(barrierCfg.Item2_nums_per_pile) * count),
	})
	return
}

// TriggerTempBuff 触发三选一
func TriggerTempBuff(logger fklog.FKLogI, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (err error) {
	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetOptionalTempBuffList(context.TODO(), userID, barrierId, 0, int32(MazeTempBuff.Type_USE_ITEM), areaId, areaIndex, 0)
	if err != nil {
		logger.ErrorWF("TriggerTempBuff GetOptionalTempBuffList fail",
			zap.Error(err),
			zap.Int32("barrierId", barrierId),
			zap.Int32("areaId", areaId),
			zap.Int32("areaIndex", areaIndex),
			zap.Int32("itemId", itemId),
			zap.Int64("count", count),
		)
		return err
	}
	optionalTempBuffListID := &MazeTempBuff.OptionalMazeTempBuffListID{
		StageId:          proto.Int32(barrierId),
		AreaId:           proto.Int32(areaId),
		AreaIndex:        proto.Int32(areaIndex),
		OptionalBuffInfo: buff.OptionalBuffInfo2PbOptionalBuffInfo(optionalBuffInfo),
	}
	// Push
	err = online.Push(logger, userID, 10552, optionalTempBuffListID)
	if err != nil {
		logger.ErrorWF("TriggerTempBuff Push fail",
			zap.Error(err),
			zap.Int32("barrierId", barrierId),
			zap.Int32("areaId", areaId),
			zap.Int32("areaIndex", areaIndex),
			zap.Int32("itemId", itemId),
			zap.Int64("count", count),
			zap.Any("optionalTempBuffListID", optionalTempBuffListID),
		)
	}
	return
}

func (g *Game) OnBarrierUseItemRQ_10550_10551(s *session.Session, req *MazeGame.BarrierUseItemRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierUseItemRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierUseItemRS{}

	logger.InfoWF("OnBarrierUseItemRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnBarrierUseItemRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.AreaId = req.AreaId
	res.AreaIndex = req.AreaIndex
	// res.ItemList = req.ItemList

	userId := uint64(s.UID())

	// 校验关卡
	barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if barrierCfg == nil {
		logger.ErrorWF("OnBarrierUseItemRQ barrier not found", zap.Any("barrierId", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡配置不存在")
		return
	}

	for _, item := range req.GetItemList() {
		if item.GetItemId() <= 0 || item.GetCount() < 0 {
			logger.ErrorWF("OnBarrierUseItemRQ invalid item param",
				zap.Any("barrierId",
					req.GetBarrierId()),
				zap.Any("ItemList", req.GetItemList()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
		itemCfg := GMazeItemsV8Cfg.Get(item.GetItemId())
		if itemCfg == nil {
			logger.ErrorWF("OnBarrierUseItemRQ item not found",
				zap.Any("barrierId", req.GetBarrierId()),
				zap.Any("ItemList", req.GetItemList()),
				zap.Any("ItemID", item.GetItemId()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
	}

	tradeNo := tradeno.GetTradeNum()
	items := make([]*MazeCommon.MazeItem, 0)

	for _, item := range req.GetItemList() {
		itemCfg := GMazeItemsV8Cfg.Get(item.GetItemId())
		switch itemCfg.Type {
		// 金币堆
		case constdef.ItemTypeGoldCoinPile:
			addItems, err := UseGoldCoinPile(logger, userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("OnBarrierUseItemRQ UseGoldCoinPile fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return err
			}
			items = append(items, addItems...)
		// 强化石堆
		case constdef.ItemTypeQianghuashiPile:
			addItems, err := UseQianghuashiPile(logger, userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("OnBarrierUseItemRQ UseQianghuashiPile fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return err
			}
			items = append(items, addItems...)
		// 三选一
		case constdef.ItemTypeTempBuff:
			err = TriggerTempBuff(logger, userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("OnBarrierUseItemRQ TriggerTempBuff fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return
			}
		case constdef.ItemTypeEnergyPotion:
			//
		}
	}

	// 处理需要加入背包的道具
	if len(items) > 0 {
		itemList := itemutil.ItemPb2ItemInfo(items)
		errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), userId, itemservice.ItemOpTypeUseItem, tradeNo, itemList...)
		if errInfo != nil {
			logger.ErrorWF("OnBarrierUseItemRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", items))
		}
	}

	// 返回新增道具列表
	res.ItemList = items

	// 保存到已获取的道具
	itemMap := make(map[int32]int64)
	for _, i := range items {
		itemMap[i.GetItemId()] += i.GetCount()
	}
	if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreRewardItem(context.TODO(), userId, req.GetBarrierId(), itemMap); err != nil {
		logger.ErrorWF("OnBarrierUseItemRQ SaveBarrierScoreRewardItem err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}

	return
}
