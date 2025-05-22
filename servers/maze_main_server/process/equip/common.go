package equip

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"gitlab.ifreetalk.com/maze-plate/protodef/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/grouplock"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipSuiteInfoV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeequipaffixrandpoolv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/mustarrive"
	"go.uber.org/zap"
)

var globalLock = grouplock.NewGroupLock(10240)

var (
	EQUIP_BAG_FULL_ERROR = &MessageType.ErrorInfo{ErrCode: proto.Int64(81001), ErrMsg: []byte("装备背包已满")}
)

func GetEquipSuitRemGroupMap(suitId int32) (map[int32]struct{}, error) {
	attrMap := make(map[int32]struct{})
	if suitId <= 0 {
		return attrMap, nil
	}
	suitCfg := GMazeEquipSuiteInfoV8Cfg.Get(suitId)
	if suitCfg == nil {
		return nil, errors.New("cannot find suit cfg")
	}
	suitRemGroupMap := mazeequipaffixrandpoolv8.GetPoolLimitCfg(suitCfg.Exclude_attr)
	return suitRemGroupMap, nil
}

// 毫秒时间戳做token
func GetToken() int64 {
	return time.Now().UnixNano() / 1000000
}

// func GetEquipBagLimit(logger fklog.FKLogI, userId uint64) int {
//	equipConfig := GMazeEquipConfigV8Cfg.Get(201)
//	if equipConfig != nil {
//		return int(equipConfig.Value_int)
//	}
//	return 50
// }

func SendMazeBagEquipChgIDEx(logger fklog.FKLogI, userId uint64, addList, delList, chgList []*MazeGameEquip.MazeEquipInfo, opType int32) {
	var mask int32
	if len(addList) > 0 {
		mask |= constdef.EquipChgTypeAdd
	}
	if len(delList) > 0 {
		mask |= constdef.EquipChgTypeDel
	}

	if len(chgList) > 0 {
		mask |= constdef.EquipChgTypeMod
	}
	if mask == 0 {
		return
	}
	req := &MazeGameEquip.MazeBagEquipID{
		EquipList:        addList,
		DelEquips:        delList,
		ChgEquips:        chgList,
		OpType:           proto.Int32(opType),
		Token:            proto.Int64(GetToken()),
		ChgType:          proto.Int32(mask),
		NeedRefreshForce: proto.Int32(1),
	}
	if opType == int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_DRESS_EQUIP) {
		req.NeedRefreshForce = proto.Int32(0)
	}
	logger.InfoWF("SendMazeBagEquipChgIDEx send client with", zap.Any("res", req))
	err := mustarrive.SendArrivePacket(logger, int64(userId), 10409, req)
	if err != nil {
		logger.ErrorWF("SendMazeBagEquipChgIDEx SendArrivePacket error", zap.Error(err))
	} else {
		logger.InfoWF("SendMazeBagEquipChgIDEx SendArrivePacket success")
	}
}
