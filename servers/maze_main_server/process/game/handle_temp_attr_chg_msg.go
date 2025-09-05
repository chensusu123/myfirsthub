/*
 * @Author: majian
 * @Date: 2025-03-07 22:13:29
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 13:47:21
 */
package game

import (
	"context"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/module/mazeuserinfo"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type MazeTempBuffChangeMsg = mazetempbuffchgmsg.MazeTempBuffChangeMsg

func HandleTempBuffMsg(ctx context.Context, msg *MazeTempBuffChangeMsg) {
	logger := fklog.ContextAppLogger(ctx)
	// msg := &structsdef.MazeTempBuffChangeMsg{}
	// err = json.Unmarshal(data, msg)
	// if err != nil {
	// 	logger.CtxError(ctx,"HandleTempBuffMsg Unmarshal", zap.Error(err),
	// 		zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
	// 	return
	// }
	userId := msg.UserId
	if msg.ChgType == 2 {
		// 由于检查buff是在进入关卡时检查，那不需要重新计算buff的技能，因为进入关卡本身会计算
		logger.CtxInfo(ctx, "HandleTempBuffMsg chgType == 2", zap.Any("msg", msg),
			zap.Uint64("userId", userId))
		return
	}

	logger.CtxInfo(ctx, "HandleTempBuffMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 || len(msg.ChgAttrs) == 0 {
		return
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "HandleTempBuffMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	mazeBattleInfo, err := GetMazeBattleData(ctx, userId, userInfo.Barrier)
	if err != nil {
		logger.CtxError(ctx, "HandleTempBuffMsg GetMazeBattleData fail", zap.Error(err))
		return
	}
	SendMazeBarrierChgPack(ctx, userId, mazeBattleInfo)

	// TODO 处理战斗数据变化(武力 血量 技能属性等)
	return
}
