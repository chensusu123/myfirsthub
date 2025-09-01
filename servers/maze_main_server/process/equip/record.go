package equip

import (
	"context"
	"fmt"
	"strings"
	"time"

	"maze_game_server/pb/server/MazeEquipCache"

	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/kafka/mazeequipinstancerecord"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type MazeGameEquipInstanceRecord struct {
	*mazeequipinstancerecord.MazeGameEquipInstanceRecord
}

func PushMazeEquipInstanceLog(ctx context.Context, equipInstanceRecordMap map[int64]*MazeGameEquipInstanceRecord, chgType, isFail int32) error {
	for _, equipRecord := range equipInstanceRecordMap {
		equipRecord.ChgType = chgType
		equipRecord.IsFail = isFail
		mazeequipinstancerecord.PushMazeGameEquipInstanceRecord(context.TODO(), equipRecord.MazeGameEquipInstanceRecord)
	}
	return nil
}

func PushMazeEquipBagLogEx(ctx context.Context, userId uint64, addEquipList, delEquipList []*MazeEquipCache.MazeEquipInfoDb, tradeNum uint64, opType, chgType int32, tempBagTime int64, isFail int32) error {
	addBagEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	delBagEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range addEquipList {
		addBagEquipList = append(addBagEquipList, equipInfo)
	}
	for _, equipInfo := range delEquipList {
		delBagEquipList = append(delBagEquipList, equipInfo)
	}
	if len(addBagEquipList) > 0 || len(delBagEquipList) > 0 {
		PushMazeEquipBagLog(ctx, userId, addBagEquipList, delBagEquipList, tradeNum, opType, chgType, isFail)
	}
	return nil
}

func PushMazeEquipBagLog(ctx context.Context, userId uint64, addEquipList, delEquipList []*MazeEquipCache.MazeEquipInfoDb, tradeNum uint64, opType, chgType, isFail int32) error {
	logger := fklog.ContextAppLogger(ctx)
	addEquipGuidStr := make([]string, 0)
	delEquipGuidStr := make([]string, 0)
	for _, equipInfo := range addEquipList {
		addEquipGuidStr = append(addEquipGuidStr, fmt.Sprintf("%d:%d", equipInfo.GetEquipGuid(), equipInfo.GetEquipId()))
	}
	for _, equipInfo := range delEquipList {
		delEquipGuidStr = append(delEquipGuidStr, fmt.Sprintf("%d:%d", equipInfo.GetEquipGuid(), equipInfo.GetEquipId()))
	}
	record := &mazeequipbagrecord.MazeGameEquipBagRecord{
		UserId:        userId,
		ChgType:       chgType,
		TradeNum:      tradeNum,
		AddEquipGuids: strings.Join(addEquipGuidStr, ","),
		DelEquipGuids: strings.Join(delEquipGuidStr, ","),
		OpType:        opType,
		IsFail:        isFail,
		// GroupID:       fkconfig.EnvVal.GroupID,
		CreateTime: time.Now().UnixNano() / 1000000,
	}
	if err := mazeequipbagrecord.PushMazeGameEquipBagRecord(ctx, record); err != nil {
		logger.CtxError(ctx, "PushMazeEquipBagLog PushMazeGameEquipBagRecord err", zap.Error(err))
	}
	return nil
}
