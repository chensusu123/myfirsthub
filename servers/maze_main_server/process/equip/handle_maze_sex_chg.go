package equip

import (
	"encoding/json"

	"context"
	"maze_game_server/common/structsdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

// 处理人偶性别变化消息

// func init() {
// RegTcpMsgCallBackFunc(constdef.KafkaMDTSex, HandleDollSexChg)
// }

// func HandleDollSexChg(logger fklog.FKLogI, userId uint64, msg []byte) error {
func HandleDollSexChg(ctx context.Context, logger fklog.FKLogI, index int, key, msg []byte) (err error) {
	info := &structsdef.SexChangeInfo{}
	err = json.Unmarshal(msg, info)
	if err != nil {
		logger.ErrorWF("HandleDollSexChg unmarshal comsume info fail", zap.Error(err))
		return err
	}
	userId := fkutil.ToUint64(info.UserId)

	// 过滤用户创建和非性别变化
	if info.CreateChg != 2 || info.FromType != 1 {
		return nil
	}

	msgUid := fkutil.ToUint64(info.UserId)
	if msgUid == 0 {
		logger.ErrorWF("HandleDollSexChg conv userid to int fail", zap.Any("info", msg))
		return nil
	}

	if userId == 0 || msgUid == 0 || userId != msgUid {
		logger.ErrorWF("HandleDollSexChg userId invalid", zap.Any("info", msg), zap.Uint64("userId", userId))
		return nil
	}
	// 检查装备位解锁
	ChkEquipPosUnlock(logger, userId, UnlockSrcSexChg, true)

	// 初始装备套检查
	InitDollEquipSuitSeq(logger, userId)

	// 处理初始化装备
	HandleDollEquipInit(ctx, userId, false)

	// 人偶属性初始化
	HandleDollAttrInit(ctx, userId, "")
	return nil
}
