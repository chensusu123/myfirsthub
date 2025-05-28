package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func OnBarrierMonsterDeathRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierMonsterDeathRQ")()

	req := rqMsg.(*MazeGame.BarrierMonsterDeathRQ)
	res := rsMsg.(*MazeGame.BarrierMonsterDeathRS)

	logger.InfoWF("OnBarrierMonsterDeathRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnBarrierMonsterDeathRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.MonsterId = req.MonsterId
	res.OpData = req.OpData

	userId := shardingID

	foeCfg := GMazeFoeV8Cfg.Get(int32(req.GetMonsterId()))
	if foeCfg == nil {
		logger.ErrorWF("OnBarrierMonsterDeathRQ get foe cfg fail", zap.Any("MonsterId", req.GetMonsterId()))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}

	// 怪物掉落装备
	tradeNo := gentradeno.GetTradeNum()
	equip := make(map[int32]int32)
	for _, v := range foeCfg.Drop_equip {
		if v > 0 {
			equip[v] += 1
		}
	}
	if len(equip) > 0 {
		_, err = addequip.AddEquipToBagWithOpdata(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_MONSTER_DEATH_AWARD), req.GetOpData(), tradeNo, equip)
		if err != nil {
			logger.ErrorWF("OnBarrierMonsterDeathRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equip))
		}
	}

	// 增加掉落物品返回
	for itemID, count := range foeCfg.Drop_item {
		if itemID > 0 {
			res.Awards = append(res.Awards, &MazeCommon.MazeItem{ItemId: proto.Int32(itemID), Count: proto.Int64(count)})
		}
	}
	// 通关值
	res.Kongfu = proto.Int32(foeCfg.Kongfu)

	return
}
