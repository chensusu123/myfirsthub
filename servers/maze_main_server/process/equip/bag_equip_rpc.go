package equip

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeequipbagrecord"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"context"
)

func OnSvrMazeEquipAssembleRQ(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	agent := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*MazeEquipSvr.SvrMazeEquipAssembleRQ)
	res := rsMsg.(*MazeEquipSvr.SvrMazeEquipAssembleRS)
	res.ErrInfo = errors.NO_ERROR
	agent.InfoWF("OnSvrMazeEquipAssembleRQ with", zap.Any("rq", req))

	defer func() {
		agent.InfoWF("OnSvrMazeEquipAssembleRQ end ", zap.Any("req", req), zap.Any("res", res))
	}()

	lock := globalLock.LockOp(uint64(shardingID))
	defer lock.Unlock()

	equipGuids := make([]int64, 0)
	if req.GetEquipGuid() > 0 {
		equipGuids = append(equipGuids, req.GetEquipGuid())
	}
	if req.GetReplacedEquipGuid() > 0 {
		equipGuids = append(equipGuids, req.GetReplacedEquipGuid())
	}
	equipMap, err := mazebagequipredis.GetBatchEquipInfo(agent, agent.UserID, equipGuids...)
	if err != nil {
		agent.ErrorWF("OnSvrMazeEquipAssembleRQ GetBatchEquipInfo error!", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var equipDelList, equipAddList []*MazeGameEquip.MazeEquipInfo
	equip := equipMap[req.GetEquipGuid()]
	if equip == nil {
		agent.ErrorWF("OnSvrMazeEquipAssembleRQ equip guid not exist", zap.Any("equipMap", equipMap), zap.Any("equipGuid", req.GetEquipGuid()))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	var replacedEquip *MazeEquipCache.MazeEquipInfoDb
	if req.GetReplacedEquipGuid() > 0 {
		replacedEquip = equipMap[req.GetReplacedEquipGuid()]
		if replacedEquip == nil {
			agent.ErrorWF("OnSvrMazeEquipAssembleRQ equip guid not exist", zap.Any("equipMap", equipMap), zap.Any("replacedEquipGuid", req.GetReplacedEquipGuid()))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
	}
	defer func() {
		PushMazeEquipBagLog(agent, agent.UserID, []*MazeEquipCache.MazeEquipInfoDb{replacedEquip}, []*MazeEquipCache.MazeEquipInfoDb{equip}, 0, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_DRESS_EQUIP), mazeequipbagrecord.MazeDressEquip, 0)
	}()
	equipCli := packtopb.EquipInfoToCliPbGuid(equip)
	equipDelList = append(equipDelList, equipCli)
	if replacedEquip != nil {
		if replacedEquip != nil {
			replacedEquipCli, err := packtopb.EquipInfoToCliPB(agent, replacedEquip)
			if err != nil {
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				ctx.ErrorWF("OnSvrDollEquipAssembleRQ EquipInfoToCliPBEx error", zap.Error(err))
				return err
			}
			equipAddList = append(equipAddList, replacedEquipCli)
		}
	}
	SendMazeBagEquipChgIDEx(agent, agent.UserID, equipAddList, equipDelList, nil, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_DRESS_EQUIP))
	return
}
