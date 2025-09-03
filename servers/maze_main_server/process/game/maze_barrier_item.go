package game

import (
	"context"

	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeAttrItemAttrV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/excel/mazeskillattrcfg"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/buff"
	"maze_game_server/services/barrieritemservice"
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
func UseGoldCoinPile(ctx context.Context, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (items []*MazeCommon.MazeItem, err error) {
	logger := fklog.ContextAppLogger(ctx)
	items = make([]*MazeCommon.MazeItem, 0)
	cfg := GMazeConfigV8Cfg.GetWithCtx(ctx, constdef.GoldPileShow2RealCfgId)
	if cfg == nil {
		err = errors.New("金币堆兑换金币配置错误")
		logger.CtxError(ctx, "UseGoldCoinPile GMazeConfigV8Cfg.Get fail", zap.Error(err), zap.Int32("cfgId", constdef.GoldPileShow2RealCfgId))
		return nil, err
	}
	goldCoinId, ok := cfg.Value_map[itemId]
	if !ok {
		err = errors.New("金币堆无映射金币配置")
		logger.CtxError(ctx, "UseGoldCoinPile no target gold coin found", zap.Error(err), zap.Int32("itemId", itemId), zap.Any("Value_map", cfg.Value_map))
		return nil, err
	}
	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrierId)
	// Add item
	items = append(items, &MazeCommon.MazeItem{
		ItemId: proto.Int32(int32(goldCoinId)),
		Count:  proto.Int64(int64(barrierCfg.Item1_nums_per_pile) * count),
	})
	return
}

// UseQianghuashiPile 使用强化石堆，转换为强化石
func UseQianghuashiPile(ctx context.Context, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (items []*MazeCommon.MazeItem, err error) {
	logger := fklog.ContextAppLogger(ctx)
	items = make([]*MazeCommon.MazeItem, 0)
	cfg := GMazeConfigV8Cfg.GetWithCtx(ctx, constdef.StrengthenStonePileShow2RealCfgId)
	if cfg == nil {
		err = errors.New("强化石堆兑换强化石配置错误")
		logger.CtxError(ctx, "UseQianghuashiPile GMazeConfigV8Cfg.Get fail", zap.Error(err), zap.Int32("cfgId", constdef.StrengthenStonePileShow2RealCfgId))
		return nil, err
	}
	goldCoinId, ok := cfg.Value_map[itemId]
	if !ok {
		err = errors.New("强化石堆无映射强化石配置")
		logger.CtxError(ctx, "UseQianghuashiPile no target qianghuashi found", zap.Error(err), zap.Int32("itemId", itemId), zap.Any("Value_map", cfg.Value_map))
		return nil, err
	}
	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrierId)
	// Add item
	items = append(items, &MazeCommon.MazeItem{
		ItemId: proto.Int32(int32(goldCoinId)),
		Count:  proto.Int64(int64(barrierCfg.Item2_nums_per_pile) * count),
	})
	return
}

// TriggerTempBuff 触发三选一
func TriggerTempBuff(ctx context.Context, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	optionalBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetOptionalTempBuffList(ctx, userID, barrierId, 0, int32(MazeTempBuff.Type_USE_ITEM), areaId, areaIndex, 0)
	if err != nil {
		logger.CtxError(ctx, "TriggerTempBuff GetOptionalTempBuffList fail",
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
		Type:             MazeTempBuff.Type_USE_ITEM.Enum(),
		Level:            proto.Int32(optionalBuffInfo.Level),
	}
	// Push
	err = online.ClusterPush(ctx, userID, 10552, optionalTempBuffListID)
	if err != nil {
		logger.CtxError(ctx, "TriggerTempBuff Push fail",
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

// TriggerSkill 触发客户端使用技能
func TriggerSkill(ctx context.Context, userID uint64, barrierId int32, areaId int32, areaIndex int32, itemId int32, count int64, skillID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	triggerSkill := &MazeAIBattle.MazeUserTriggerSkillInfoID{}
	triggerSkill.SkillIds = append(triggerSkill.SkillIds, skillID)
	err = online.ClusterPush(ctx, userID, 10666, triggerSkill)
	if err != nil {
		logger.CtxError(ctx, "TriggerSkill ClusterPush fail",
			zap.Error(err),
			zap.Int32("barrierId", barrierId),
			zap.Int32("areaId", areaId),
			zap.Int32("areaIndex", areaIndex),
			zap.Int32("itemId", itemId),
			zap.Int64("count", count),
			zap.Int32("skillID", skillID),
		)
	}
	logger.CtxInfo(ctx, "TriggerSkill success",
		zap.Error(err),
		zap.Int32("barrierId", barrierId),
		zap.Int32("areaId", areaId),
		zap.Int32("areaIndex", areaIndex),
		zap.Int32("itemId", itemId),
		zap.Int64("count", count),
		zap.Int32("skillID", skillID),
	)
	return
}

func (g *Game) OnBarrierUseItemRQ_10550_10551(s *session.Session, req *MazeGame.BarrierUseItemRQ) (err error) {
	logger := fklog.ContextAppLogger(s.Context())
	defer fkprometheus.InfoPMT("OnBarrierUseItemRQ")()

	res := &MazeGame.BarrierUseItemRS{}
	ctx := s.Context()
	logger.CtxInfo(s.Context(), "OnBarrierUseItemRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(s.Context(), "OnBarrierUseItemRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.AreaId = req.AreaId
	res.AreaIndex = req.AreaIndex
	// res.ItemList = req.ItemList

	userId := uint64(s.UID())

	// 校验关卡
	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(s.Context(), req.GetBarrierId())
	if barrierCfg == nil {
		logger.CtxError(s.Context(), "OnBarrierUseItemRQ barrier not found", zap.Any("barrierId", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡配置不存在")
		return
	}

	var (
		subItems  = make([]*itemservice.ItemInfo, 0)
		subEquips = make([]*itemservice.ItemInfo, 0)
	)

	for _, item := range req.GetItemList() {
		if item.GetItemId() <= 0 || item.GetCount() < 0 {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ invalid item param",
				zap.Any("barrierId", req.GetBarrierId()),
				zap.Any("ItemList", req.GetItemList()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
		itemCfg := GMazeItemsV8Cfg.GetWithCtx(s.Context(), item.GetItemId())
		if itemCfg == nil {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ item not found",
				zap.Any("barrierId", req.GetBarrierId()),
				zap.Any("ItemList", req.GetItemList()),
				zap.Any("ItemID", item.GetItemId()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
		subEquips = append(subEquips, &itemservice.ItemInfo{ItemId: item.GetItemId(), Count: item.GetCount()})
	}
	for _, equip := range req.GetEquipList() {
		if equip.GetItemId() <= 0 || equip.GetCount() < 0 {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ invalid equip param",
				zap.Any("barrierId", req.GetBarrierId()),
				zap.Any("EquipList", req.GetEquipList()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
		subEquips = append(subEquips, &itemservice.ItemInfo{ItemId: equip.GetItemId(), Count: equip.GetCount()})
	}

	tradeNo := tradeno.GetTradeNum()
	items := make([]*MazeCommon.MazeItem, 0)

	// 扣减场中物品
	ok, err := barrieritemservice.GbarrierItemsService.TrySubBarrierItems(s.Context(), userId, req.GetBarrierId(), subItems, subEquips)
	if err != nil {
		logger.CtxError(s.Context(), "OnBarrierUseItemRQ UseGoldCoinPile fail",
			zap.Any("barrierId", req.GetBarrierId()),
			zap.Any("subItems", subItems),
			zap.Any("ItemList", req.GetItemList()),
			zap.Any("subEquips", subEquips),
			zap.Any("EquipList", req.GetEquipList()),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("检查掉落物品失败")
		return err
	}
	if !ok {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("掉落物品不存在")
		return nil
	}

	// 使用道具
	for _, item := range req.GetItemList() {
		itemCfg := GMazeItemsV8Cfg.GetWithCtx(s.Context(), item.GetItemId())
		switch itemCfg.Type {
		// 金币堆
		case constdef.ItemTypeGoldCoinPile:
			addItems, err := UseGoldCoinPile(s.Context(), userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ UseGoldCoinPile fail",
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
			addItems, err := UseQianghuashiPile(s.Context(), userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ UseQianghuashiPile fail",
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
			err = TriggerTempBuff(s.Context(), userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount())
			if err != nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ TriggerTempBuff fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return
			}
		// 技能道具
		case constdef.ItemTypeSkill:
			skillItemCfg := GMazeAttrItemAttrV8Cfg.GetWithCtx(s.Context(), item.GetItemId())
			if skillItemCfg == nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ GMazeAttrItemAttrV8Cfg.GetWithCtx fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return
			}
			skillCfg := mazeskillattrcfg.GetAttrSkill(skillItemCfg.Add_attr)
			if skillCfg == nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ mazeskillattrcfg.GetAttrSkill fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return
			}
			err = TriggerSkill(s.Context(), userId, req.GetBarrierId(), req.GetAreaId(), req.GetAreaIndex(), item.GetItemId(), item.GetCount(), skillCfg.Id)
			if err != nil {
				logger.CtxError(s.Context(), "OnBarrierUseItemRQ TriggerSkill fail",
					zap.Any("barrierId", req.GetBarrierId()),
					zap.Any("ItemID", item.GetItemId()),
					zap.Any("ItemList", req.GetItemList()),
				)
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具使用失败")
				return
			}
		}
	}

	addEquipMap := make(map[int32]int32)
	// 捡装备
	for _, equip := range req.GetEquipList() {
		var (
			count   = equip.GetCount()
			equipID = equip.GetItemId()
		)
		if equipID <= 0 || count <= 0 {
			logger.CtxWarn(s.Context(), "OnBarrierUseItemRQ TriggerSkill fail",
				zap.Any("barrierId", req.GetBarrierId()),
				zap.Any("equipID", equipID),
				zap.Any("count", count),
				zap.Any("EquipList", req.GetEquipList()),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效装备参数")
			return
		}
		addEquipMap[equipID] += int32(count)
	}

	// 处理需要加入背包的道具
	if len(items) > 0 {
		itemList := itemutil.ItemPb2ItemInfo(items)
		errInfo := itemservice.GlobalItemService.AddItem(ctx, userId, itemservice.ItemOpTypeUseItem, tradeNo, itemList...)
		if errInfo != nil {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", items))
		}
		// 保存到已获取的道具
		itemMap := make(map[int32]int64)
		for _, i := range items {
			itemMap[i.GetItemId()] += i.GetCount()
		}
		if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreRewardItem(s.Context(), userId, req.GetBarrierId(), itemMap); err != nil {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ SaveBarrierScoreRewardItem err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
	}
	// 返回新增道具列表
	res.ItemList = items

	if len(addEquipMap) > 0 {
		rs, err2 := addequip.AddEquipToBag(s.Context(), userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE), tradeNo, addEquipMap)
		if err2 != nil {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ addEquipToBag fail", zap.Error(err2), zap.Any("optype", MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", addEquipMap), zap.Any("rs", rs))
		}
		if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreRewardEquip(s.Context(), userId, req.GetBarrierId(), addEquipMap); err != nil {
			logger.CtxError(s.Context(), "OnBarrierUseItemRQ SaveBarrierScoreRewardEquip err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
	}

	return
}
