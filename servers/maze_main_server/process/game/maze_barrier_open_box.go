package game

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/barrierservice"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnBarrierOpenBoxRQ_10445_10446(s *session.Session, req *MazeGame.BarrierOpenBoxRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierOpenBoxRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierOpenBoxRS{}

	logger.InfoWF("OnBarrierOpenBoxRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnBarrierOpenBoxRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.BoxId = req.BoxId
	res.OpData = req.OpData

	userId := uint64(s.UID())

	// 关卡中打开宝箱
	kongfu, equips, items, errinfo := barrierservice.Global.OpenBox(logger, userId, req.GetBarrierId(), int32(req.GetBoxId()))
	if errinfo.GetErrCode() != errors.NO_ERROR_CODE {
		res.ErrInfo = errinfo
		logger.ErrorWF("OnBarrierOpenBoxRQ OpenBox fail", zap.Error(fmt.Errorf("OpenBox: %s", errinfo.GetErrMsg())), zap.Any("boxId", req.GetBoxId()))
		return
	}

	// 通关值
	res.Kongfu = proto.Int32(kongfu)

	tradeNo := tradeno.GetTradeNum()

	// TODO 使用equip,item服务增加奖励

	// 怪物掉落装备
	if len(equips) > 0 {
		_, err = addequip.AddEquipToBagWithOpdata(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD), req.GetOpData(), tradeNo, equips)
		if err != nil {
			logger.ErrorWF("OnBarrierOpenBoxRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equips))
		}
	}

	var bagItems []*MazeCommon.MazeItem
	// 增加掉落物品返回
	for itemID, count := range items {
		if itemID > 0 {
			itemCfg := GMazeItemsV8Cfg.Get(itemID)
			if itemCfg == nil {
				logger.ErrorWF("OnBarrierOpenBoxRQ item not found", zap.Error(fmt.Errorf("item: %d not found", itemID)), zap.Any("boxId", req.GetBoxId()))
			} else {
				if itemID == constdef.MazeCommonItemCoin || // 金币
					itemID == constdef.MazeCommonItemDiamond { // 钻石
					bagItems = append(bagItems, &MazeCommon.MazeItem{
						ItemId: proto.Int32(itemID),
						Count:  proto.Int64(count),
					})
				} else {
					res.Awards = append(res.Awards, &MazeCommon.MazeItem{
						ItemId: proto.Int32(itemID),
						Count:  proto.Int64(count),
					})
					// 背包道具
					if itemCfg.Is_bag == 3 {
						bagItems = append(bagItems, &MazeCommon.MazeItem{
							ItemId: proto.Int32(itemID),
							Count:  proto.Int64(count),
						})
					}
				}
			}
		}
	}

	// 处理需要加入背包的道具
	if len(bagItems) > 0 {
		itemList := itemutil.ItemPb2ItemInfo(bagItems)
		errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), userId, itemservice.ItemOpTypeOpenBox, tradeNo, itemList...)
		if errInfo != nil {
			logger.ErrorWF("OnBarrierOpenBoxRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("bagItems", bagItems))
		}
	}

	// 保存到已获取的道具
	if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreReward(context.TODO(), userId, req.GetBarrierId(), equips, items); err != nil {
		logger.ErrorWF("OnBarrierPickItemRQ SaveBarrierScoreRewardItem err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}

	return nil
}
