package equip

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/module/bagmodule"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/errors"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/pb/server/MazeEquipSvr"
)

func OnSvrDollEquipSaleRQ(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	agent := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*MazeEquipSvr.SvrMazeEquipSaleRQ)
	res := rsMsg.(*MazeEquipSvr.SvrMazeEquipSaleRS)
	res.ErrInfo = errors.NO_ERROR
	agent.InfoWF("OnSvrDollEquipSaleRQ with", zap.Any("rq", req))

	defer func() {
		agent.InfoWF("OnSvrDollEquipSaleRQ end ", zap.Any("req", req), zap.Any("res", res))
	}()

	lock := globalLock.LockOp(uint64(shardingID))
	defer lock.Unlock()

	bagEquipMgr := bagmodule.NewBagEquipMgr(agent, agent.UserID)
	err = bagEquipMgr.LoadBagFromRedis()
	if err != nil {
		agent.ErrorWF("OnSvrDollEquipSaleRQ LoadBagFromRedis error!", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	equipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0, len(req.GetEquipGuids()))
	for _, equipGuid := range req.GetEquipGuids() {
		if equipInfo, ok := bagEquipMgr.MainBagEquips.GetMainEquip()[equipGuid]; !ok {
			agent.ErrorWF("OnSvrDollEquipSaleRQ GetBatchEquipInfo error!", zap.Any("equipGuid", equipGuid), zap.Any("equipInfo", equipInfo))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		} else {
			equipList = append(equipList, equipInfo)
			bagEquipMgr.AddBagRem(equipInfo)
		}
	}

	err = bagEquipMgr.SaveBagInfoToRedis()
	if err != nil {
		agent.ErrorWF("OnSvrDollEquipSaleRQ SaveBagInfoToRedis error!", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		PushMazeEquipBagLogEx(agent, agent.UserID, nil, equipList, req.GetTradeNum(), 5, mazeequipbagrecord.MazeDelEquip, 0, 1)
		return
	}
	equipCliList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	// tempEquipCliList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	for _, equipInfo := range equipList {
		equipCliList = append(equipCliList, packtopb.EquipInfoToCliPbGuid(equipInfo))
	}
	PushMazeEquipBagLogEx(agent, agent.UserID, nil, equipList, req.GetTradeNum(), 5, mazeequipbagrecord.MazeDelEquip, 0, 0)
	if len(equipCliList) > 0 {
		SendMazeBagEquipChgIDEx(agent, agent.UserID, nil, equipCliList, nil, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SALE))
	}
	return
}
