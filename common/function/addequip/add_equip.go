package addequip

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/io/rpc/dollequipbagrpc"
	"maze_game_server/pb/server/MazeEquipSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 商城购买装备
func AddEquipToBag(ctx context.Context, userId uint64, opType int32, tradeNo uint64, equipMap map[int32]int32) (rsAdd *MazeEquipSvr.SvrAddMazeEquipRS, err error) {
	logger := fklog.ContextAppLogger(ctx)
	rqAdd := &MazeEquipSvr.SvrAddMazeEquipRQ{
		UserId:      proto.Uint64(userId),
		OpType:      proto.Int32(opType),
		TradeNumber: proto.Uint64(tradeNo),
	}
	for equipId, count := range equipMap {
		for index := int32(0); index < count; index++ {
			equipCond := &MazeEquipSvr.SvrEquipInfo{EquipId: proto.Int32(equipId)}
			rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
		}
	}
	// for k := range equips {
	//	equipCond := &DollEquipSvr.SvrEquipInfo{EquipId: proto.Int32(k)}
	//	// 初始化武器子类型固定是1 刀
	//	equipCfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(k)
	//	if equipCfg != nil && equipCfg.Pos == 1 {
	//		equipCond.Conditions = append(equipCond.Conditions, &DollEquipSvr.ConditionInfo{
	//			Id: proto.Int32(int32(DollEquipSvr.EQUIP_ADD_CONDITION_ASSIGN_EQUIP_SUB_TYPE)), Value: proto.Int64(1)})
	//	}
	//	rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
	// }
	// rqAdd.NotNotify = proto.Bool(true)

	rsAdd = &MazeEquipSvr.SvrAddMazeEquipRS{}
	err = dollequipbagrpc.MazeBagAddRQ(ctx, rqAdd, rsAdd)
	if err != nil {
		logger.CtxError(ctx, "addEquipToBag fail", zap.Error(err), zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
	} else {
		if rsAdd.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
			logger.CtxError(ctx, "addEquipToBag rs fail", zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
			err = errors.New("装备加背包失败")
		}
	}
	return
}

func AddEquipToBagWithOpdata(ctx context.Context, userId uint64, opType int32, opData string, tradeNo uint64, equipMap map[int32]int32) (rsAdd *MazeEquipSvr.SvrAddMazeEquipRS, err error) {
	logger := fklog.ContextAppLogger(ctx)
	rqAdd := &MazeEquipSvr.SvrAddMazeEquipRQ{
		UserId:      proto.Uint64(userId),
		OpType:      proto.Int32(opType),
		TradeNumber: proto.Uint64(tradeNo),
	}
	for equipId, count := range equipMap {
		for index := int32(0); index < count; index++ {
			equipCond := &MazeEquipSvr.SvrEquipInfo{EquipId: proto.Int32(equipId)}
			rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
		}
	}
	// for k := range equips {
	//	equipCond := &DollEquipSvr.SvrEquipInfo{EquipId: proto.Int32(k)}
	//	// 初始化武器子类型固定是1 刀
	//	equipCfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(k)
	//	if equipCfg != nil && equipCfg.Pos == 1 {
	//		equipCond.Conditions = append(equipCond.Conditions, &DollEquipSvr.ConditionInfo{
	//			Id: proto.Int32(int32(DollEquipSvr.EQUIP_ADD_CONDITION_ASSIGN_EQUIP_SUB_TYPE)), Value: proto.Int64(1)})
	//	}
	//	rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
	// }
	// rqAdd.NotNotify = proto.Bool(true)

	rsAdd = &MazeEquipSvr.SvrAddMazeEquipRS{}
	err = dollequipbagrpc.MazeBagAddRQWithOpData(ctx, rqAdd, rsAdd, opData)
	if err != nil {
		logger.CtxError(ctx, "AddEquipToBagWithOpdata fail", zap.Error(err), zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
	} else {
		if rsAdd.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
			logger.CtxError(ctx, "AddEquipToBagWithOpdata rs fail", zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
			err = errors.New("装备加背包失败")
		}
	}
	return
}

// // 实例化装备
// func InstanceEquip(ctx context.Context, userId uint64, opType int32, tradeNo uint64, equipNumPerCycle int32, equipMap map[int32]int32) (rsAdd *MazeEquipSvr.SvrMazeEquipInstanceRS, err error) {
//	rqAdd := &MazeEquipSvr.SvrMazeEquipInstanceRQ{
//		UserId:      proto.Uint64(userId),
//		EquipList:   make([]*MazeEquipSvr.SvrInstanceEquipInfo, 0),
//		OpType:      proto.Int32(opType),
//		TradeNumber: proto.Uint64(tradeNo),
//		Desc:        []byte{},
//	}
//	for equipId, count := range equipMap {
//		for index := int32(0); index < count; index++ {
//			equipCond := &MazeEquipSvr.SvrInstanceEquipInfo{EquipId: proto.Int32(equipId), PopupEquipNum: proto.Int32(equipNumPerCycle)}
//			rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
//		}
//	}
//
//	rsAdd = &MazeEquipSvr.SvrMazeEquipInstanceRS{}
//	err = dollequipbagrpc.MazeBagInstanceRQ(logger, rqAdd, rsAdd)
//	if err != nil {
//		logger.CtxError(ctx,"InstanceEquip fail", zap.Error(err), zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
//	} else {
//		if rsAdd.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
//			logger.CtxError(ctx,"InstanceEquip rs fail", zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
//			err = errors.New("装备实例化失败")
//		}
//	}
//	return
// }
