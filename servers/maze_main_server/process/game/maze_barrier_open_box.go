package game

import (
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/config/GMazeBoxV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"google.golang.org/protobuf/proto"

	"go.uber.org/zap"
)

func OnBarrierOpenBoxRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnBarrierOpenBoxRQ")()

	req := rqMsg.(*MazeGame.BarrierOpenBoxRQ)
	res := rsMsg.(*MazeGame.BarrierOpenBoxRS)

	logger.InfoWF("OnBarrierOpenBoxRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnBarrierOpenBoxRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.BoxId = req.BoxId
	res.OpData = req.OpData

	userId := shardingID

	boxCfg := GMazeBoxV8Cfg.Get(int32(req.GetBoxId()))
	if boxCfg == nil {
		logger.ErrorWF("OnBarrierOpenBoxRQ get box cfg fail", zap.Any("boxId", req.GetBoxId()))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}

	// 更新box表格 读表校验宝箱对应的关卡id
	if fkutil.ToInt64(boxCfg.Level_id) != int64(req.GetBarrierId()) {
		logger.ErrorWF("OnBarrierOpenBoxRQ barrier and box not match", zap.Any("boxId", req.GetBoxId()), zap.Any("barrierId", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("宝箱关卡信息不匹配")
		return
	}

	// 宝箱掉落装备
	tradeNo := gentradeno.GetTradeNum()
	equip := make(map[int32]int32)
	for _, v := range boxCfg.Award_equip {
		if v > 0 {
			equip[v] += 1
		}
	}
	if len(equip) > 0 {
		_, err = addequip.AddEquipToBagWithOpdata(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD), req.GetOpData(), tradeNo, equip)
		if err != nil {
			logger.ErrorWF("OnBarrierOpenBoxRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equip))
		}
	}

	var bagItems []*MazeCommon.MazeItem
	// 增加掉落物品返回
	for itemID, count := range boxCfg.Drop_item {
		if itemID > 0 {
			itemCfg := GMazeItemsV8Cfg.Get(itemID)
			if itemCfg == nil {
				logger.ErrorWF("OnBarrierOpenBoxRQ item not found", zap.Error(fmt.Errorf("item: %d not found", itemID)), zap.Any("boxId", req.GetBoxId()))
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
	// 处理需要加入背包的道具
	if len(bagItems) > 0 {
		errInfo := gentradeno.AddItemEx(logger, userId, 697, tradeNo, req.GetHeader(), bagItems...)
		if errInfo != nil {
			logger.ErrorWF("OnMazeBarrierPassRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("bagItems", bagItems))
		}
	}

	// 通关值
	res.Kongfu = proto.Int32(boxCfg.Add_kongfu)

	return nil
}
