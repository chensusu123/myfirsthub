package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
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
		if item.GetItemId() <= 0 || item.GetCount() <= 0 {
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

	equip := make(map[int32]int32)
	for _, eq := range req.GetEquipList() {
		if eq.GetItemId() <= 0 || eq.GetCount() <= 0 {
			logger.ErrorWF("OnBarrierPickItemRQ invalid equip param", zap.Any("barrierId", req.GetBarrierId()), zap.Any("EquipList", req.GetEquipList()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备参数错误")
			return
		}
		equip[eq.GetItemId()] += int32(eq.GetCount())
	}

	// 宝箱掉落装备
	tradeNo := gentradeno.GetTradeNum()

	if len(equip) > 0 {
		_, err = addequip.AddEquipToBagWithOpdata(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD), req.GetOpData(), tradeNo, equip)
		if err != nil {
			logger.ErrorWF("OnBarrierPickItemRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equip))
		}
	}

	// 处理需要加入背包的道具
	if len(req.GetItemList()) > 0 {
		errInfo := gentradeno.AddItemEx(logger, userId, 697, tradeNo, req.GetHeader(), req.GetItemList()...)
		if errInfo != nil {
			logger.ErrorWF("OnBarrierPickItemRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", req.GetItemList()))
		}
	}

	return
}
