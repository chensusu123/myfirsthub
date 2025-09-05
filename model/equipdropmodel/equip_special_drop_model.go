package equipdropmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func getRedisKey(userId uint64) string {
	return fmt.Sprintf("maze:equip:special:drop:u:%d", userId)
}

type EquipSpecialDropModel struct {
	EquipPoints int32           `json:"equip_points"` // 装备分
	DropMap     map[int32]int32 `json:"drop_map"`     // 特殊掉落记录：key=MazeEquDropV8ConfigRow表序列id，value=已经掉落的特殊掉落下标
}

func NewEquipSpecialDropModel(ctx context.Context, userID uint64) (*EquipSpecialDropModel, error) {
	info := &EquipSpecialDropModel{
		DropMap: make(map[int32]int32),
	}
	if err := info.load(ctx, userID); err != nil {
		return nil, err
	}
	return info, nil
}

func (info *EquipSpecialDropModel) load(ctx context.Context, userID uint64) (err error) {
	//bytes, err := mazeequipspecialdropredis.GetMazeEquipSpecialDropInfo(ctx, userID)
	//if err != nil {
	//	return err
	//}
	//if bytes == nil {
	//	return nil
	//}
	//err = serialize.Unmarshal(bytes, info)
	//if err != nil {
	//	logger.CtxError(ctx,"EquipSpecialDropModel load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID))
	//	return err
	//}

	logger := fklog.ContextAppLogger(ctx)
	err = io.LoadSvrData(context.TODO(), getRedisKey(userID), info)
	if err != nil {
		logger.CtxError(ctx, "EquipSpecialDropModel load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID))
		return err
	}
	return
}

func (info *EquipSpecialDropModel) Save(ctx context.Context, userID uint64) (err error) {
	//bytes, err := serialize.Marshal(info)
	//if err != nil {
	//	logger.CtxError(ctx,"EquipSpecialDropModel save Marshal failed", zap.Error(err), zap.Uint64("userID", userID))
	//	return err
	//}
	//return mazeequipspecialdropredis.SetMazeEquipSpecialDropInfo(logger, userID, bytes)
	logger := fklog.ContextAppLogger(ctx)
	err = io.SaveSvrData(context.TODO(), getRedisKey(userID), info)
	if err != nil {
		logger.CtxError(ctx, "EquipSpecialDropModel load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID))
		return err
	}
	return nil
}

func (info *EquipSpecialDropModel) Del(ctx context.Context, userID uint64) (err error) {
	//return mazeequipspecialdropredis.GMDel(logger, userID)
	return io.DeleteSvrData(ctx, getRedisKey(userID))
}
