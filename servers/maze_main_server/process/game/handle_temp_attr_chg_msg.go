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
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrieritemservice"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	cur, max, cd, err := barrieritemservice.GbarrierItemsService.CheckBloodAttr(ctx, userId, userInfo.Barrier)
	if err != nil {
		logger.CtxError(ctx, "HandleTempBuffMsg CheckBloodAttr fail", zap.Error(err), zap.Int32("Barrier", userInfo.Barrier))
		return
	}
	bloodDrugUseInfo := &MazeGame.BloodDrugUseInfo{
		CurCount: proto.Int32(int32(cur)),
		MaxCount: proto.Int32(int32(max)),
		Cooldown: proto.Int32(int32(cd)),
	}
	SendMazeBloodChgPack(ctx, userId, bloodDrugUseInfo)
}

// SendMazeBloodChgPack 向客户端推送血瓶信息变化包
func SendMazeBloodChgPack(ctx context.Context, userId uint64, bloodDrugInfo *MazeGame.BloodDrugUseInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	moneyPack := &MazeGame.MazeBloodDrugUseInfoChgID{
		BloodDrugUseInfo: bloodDrugInfo,
	}
	logger.CtxInfo(ctx, "SendMazeBloodChgPack send client with", zap.Uint64("userId", userId), zap.Any("bloodDrugInfo", bloodDrugInfo))
	defer func() {
		if err != nil {
			logger.CtxError(ctx, "SendMazeBloodChgPack send client fail", zap.Error(err), zap.Uint64("userId", userId), zap.Any("bloodDrugInfo", bloodDrugInfo))
		}
	}()
	return online.ClusterPush(ctx, uint64(userId), 10687, moneyPack)
}
