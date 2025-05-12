/*
 * @Author: majian
 * @Date: 2025-03-12 16:30:38
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 17:51:47
 */
package equip

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/demconstdef"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeLevelV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeBuffData"
	"go.uber.org/zap"
)

type MazeUserLevelRecord = mazeuserlevelkafka.MazeUserLevelRecord

func HandleMazeLvChg(logger fklog.FKLogI, pack *MazeUserLevelRecord) {
	// pack := &structsdef.MazeUserLevelRecord{}
	// err = json.Unmarshal(data, pack)
	// if err != nil {
	// 	logger.ErrorWF("HandleMazeLvChg Unmarshal",
	// 		zap.String("Value", string(data)),
	// 		zap.Any("err", err),
	// 	)
	// 	return err
	// }
	if pack.UserId <= 0 {
		return
	}
	logger.SetUid(pack.UserId)
	logger.WarnWF("HandleMazeLvChg recv kafka notify", zap.Any("pack", pack))

	ChkEquipPosUnlock(logger, pack.UserId, UnlockSrcDollLv, true)

	UpdateMazeLvBuff(logger, pack.UserId)
	return
}

func UpdateMazeLvBuff(logger fklog.FKLogI, userId uint64) error {
	lv, err := mazeuserlevelredis.GetUserLevel(logger, userId)
	if err != nil {
		return err
	}
	mazeLvCfg := GMazeLevelV8Cfg.Get(int32(lv))
	if mazeLvCfg == nil {
		logger.ErrorWF("UpdateMazeLvBuff no found cfg", zap.Int64("lv", lv))
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
	err = mazebuffinforedis.SaveMazeLvBuff(logger, userId, otherDb)
	if err != nil {
		logger.ErrorWF("UpdateMazeLvBuff SaveMazeLvBuff fail", zap.Error(err),
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
	e := mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	if e != nil {
		logger.ErrorWF("UpdateMazeLvBuff SendDollAttrCalcNotify fail", zap.Error(e))
	}

	logger.InfoWF("UpdateMazeLvBuff end", zap.Any("otherDb", otherDb),
		zap.Int64("mazeLv", lv))
	return nil
}
