package game

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierscorerewardservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnBarrierPickItemRQ_10527_10528(s *session.Session, req *MazeGame.BarrierPickItemRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierPickItemRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierPickItemRS{}

	logger.InfoWF("OnBarrierPickItemRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnBarrierPickItemRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.ItemList = req.ItemList
	res.EquipList = req.EquipList
	res.OpData = req.OpData

	userId := uint64(s.UID())

	for _, item := range req.GetItemList() {
		if item.GetItemId() <= 0 || item.GetCount() < 0 {
			logger.ErrorWF("OnBarrierPickItemRQ invalid item param", zap.Any("barrierId", req.GetBarrierId()), zap.Any("ItemList", req.GetItemList()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
		itemCfg := GMazeItemsV8Cfg.Get(item.GetItemId())
		if itemCfg == nil {
			logger.ErrorWF("OnBarrierPickItemRQ item not found", zap.Any("barrierId", req.GetBarrierId()), zap.Any("ItemList", req.GetItemList()), zap.Any("ItemID", item.GetItemId()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("道具参数错误")
			return
		}
	}

	//equip := make(map[int32]int32)
	//for _, eq := range req.GetEquipList() {
	//	if eq.GetItemId() <= 0 || eq.GetCount() <= 0 {
	//		logger.ErrorWF("OnBarrierPickItemRQ invalid equip param", zap.Any("barrierId", req.GetBarrierId()), zap.Any("EquipList", req.GetEquipList()))
	//		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备参数错误")
	//		return
	//	}
	//	equip[eq.GetItemId()] += int32(eq.GetCount())
	//}

	tradeNo := gentradeno.GetTradeNum()

	//if len(equip) > 0 {
	//	_, err = addequip.AddEquipToBagWithOpdata(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD), req.GetOpData(), tradeNo, equip)
	//	if err != nil {
	//		logger.ErrorWF("OnBarrierPickItemRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
	//			zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equip))
	//	}
	//}

	// 传过来的是显示的道具，需要转为实际的道具  例如：金币堆->金币
	realAddItemList := make([]*MazeCommon.MazeItem, 0)
	for _, item := range req.GetItemList() {
		var cfgId int32 = 0
		if item.GetItemId() == constdef.GoldPileItemCfgId {
			cfgId = constdef.GoldPileShow2RealCfgId
		} else if item.GetItemId() == constdef.StrengthenStonePileItemCfgId {
			cfgId = constdef.StrengthenStonePileShow2RealCfgId
		} else {
			logger.ErrorWF("OnBarrierPickItemRQ cfg not exist", zap.Int32("cfgId", item.GetItemId()))
			return errors.New("拾取道具错误")
		}
		cfg := GMazeConfigV8Cfg.Get(cfgId)
		if cfg == nil {
			logger.ErrorWF("OnBarrierPickItemRQ cfg not exist", zap.Int32("cfgId", item.GetItemId()))
			return errors.New("道具配置不存在")
		}
		realAddItemId := cfg.Value_map[item.GetItemId()]
		barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
		if barrierCfg == nil {
			logger.ErrorWF("OnBarrierPickItemRQ cfg not exist", zap.Int32("barrier", req.GetBarrierId()))
			return errors.New("关卡配置不存在")
		}
		var realCount int32 = 0
		// 找到显示道具实际加的数量
		if item.GetItemId() == constdef.GoldPileItemCfgId {
			realCount = barrierCfg.Item1_nums_per_pile
		} else if item.GetItemId() == constdef.StrengthenStonePileItemCfgId {
			realCount = barrierCfg.Item2_nums_per_pile
		}
		realAddItemList = append(realAddItemList, &MazeCommon.MazeItem{
			ItemId: proto.Int32(int32(realAddItemId)),
			Count:  proto.Int64(int64(realCount * int32(item.GetCount()))),
		})
	}

	// 处理需要加入背包的道具
	if len(realAddItemList) > 0 {
		errInfo := gentradeno.AddItemEx(logger, userId, 697, tradeNo, req.GetHeader(), realAddItemList...)
		if errInfo != nil {
			logger.ErrorWF("OnBarrierPickItemRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", realAddItemList))
		}
	}

	// 保存到已获取的道具
	itemMap := make(map[int32]int64)
	for _, i := range realAddItemList {
		itemMap[i.GetItemId()] += i.GetCount()
	}
	if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreRewardItem(logger, userId, req.GetBarrierId(), itemMap); err != nil {
		logger.ErrorWF("OnBarrierPickItemRQ SaveBarrierScoreRewardItem err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}

	return
}
