/*
 * @Author: majian
 * @Date: 2025-03-12 16:30:38
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 17:51:47
 */
package equip

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/servers/maze_main_server/process/equip/demconstdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type MazeUserLevelRecord = mazeuserlevelkafka.MazeUserLevelRecord

func HandleMazeLvChg(ctx context.Context, pack *MazeUserLevelRecord) {
	logger := fklog.ContextAppLogger(ctx)
	// pack := &structsdef.MazeUserLevelRecord{}
	// err = json.Unmarshal(data, pack)
	// if err != nil {
	// 	logger.CtxError(ctx,"HandleMazeLvChg Unmarshal",
	// 		zap.String("Value", string(data)),
	// 		zap.Any("err", err),
	// 	)
	// 	return err
	// }
	if pack.UserId <= 0 {
		return
	}
	logger.SetUid(pack.UserId)
	logger.CtxWarn(ctx, "HandleMazeLvChg recv kafka notify", zap.Any("pack", pack))

	ChkEquipPosUnlock(ctx, pack.UserId, UnlockSrcDollLv, true)

	UpdateMazeLvBuff(ctx, pack.UserId)
	return
}

func UpdateMazeLvBuff(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	lv, err := mazeuserlevelredis.GetUserLevel(ctx, userId)
	if err != nil {
		return err
	}
	mazeLvCfg := GMazeLevelV8Cfg.GetWithCtx(ctx, int32(lv))
	if mazeLvCfg == nil {
		logger.CtxError(ctx, "UpdateMazeLvBuff no found cfg", zap.Int64("lv", lv))
		return nil
	}
	otherDb := &MazeBuffData.MazeBuffDb{}
	for k, v := range mazeLvCfg.Attr {
		if k <= 0 {
			continue
		}
		otherDb.MazeShowBuffs = append(otherDb.MazeShowBuffs, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(k),
			AttrVal: proto.Int64(v),
		})
	}
	err = mazebuffinforedis.SaveMazeLvBuff(ctx, userId, otherDb)
	if err != nil {
		logger.CtxError(ctx, "UpdateMazeLvBuff SaveMazeLvBuff fail", zap.Error(err),
			zap.Any("otherDb", otherDb),
			zap.Int64("mazeLv", lv))
		return err
	}
	// 通知计算属性
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
	calcAttrNotify.FromServer = demconstdef.MySvr
	calcAttrNotify.UserId = userId
	calcAttrNotify.ChgType = constdef.MazeBuffLvChg
	calcAttrNotify.Session = ""
	calcAttrNotify.BuffSrc = constdef.MazeBuffSrcLv
	e := mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
	if e != nil {
		logger.CtxError(ctx, "UpdateMazeLvBuff SendDollAttrCalcNotify fail", zap.Error(e))
	}

	logger.CtxInfo(ctx, "UpdateMazeLvBuff end", zap.Any("otherDb", otherDb),
		zap.Int64("mazeLv", lv))
	return nil
}
