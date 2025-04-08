package process

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBoxV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/function/addequip"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/function/gentradeno"

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

	tradeNo := gentradeno.GetTradeNum()
	equip := make(map[int32]int32)
	for _, v := range boxCfg.Award_equip {
		equip[v] += 1
	}

	// DOLL_EQUIP_MAZE_FOE = 19;//迷宫怪物掉落
	// DOLL_EQUIP_MAZE_BOX_AWARD = 20;//迷宫宝箱掉落
	_, err = addequip.AddEquipToBag(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD), tradeNo, equip)
	if err != nil {
		logger.ErrorWF("OnBarrierOpenBoxRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
			zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equip))
	}
	return nil
}
