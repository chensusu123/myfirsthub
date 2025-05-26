package card

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazecardlistgroupredis"
	"maze_game_server/io/redis/userriddlemonthlyredis"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/pb/common/MazeCard"
	"maze_game_server/pb/server/MazeBuffData"
)

// 月卡变化通知
type MazeCardNotifyMsg struct {
	UserId       uint64 `json:"user_id"`
	NewValue     int64  `json:"new_value"`
	OldValue     int64  `json:"old_value"`
	ChgValue     int    `json:"chg_value"`
	Type         int    `json:"type"`
	CreateTime   int64  `json:"create_time"`
	CreateMsTime int64  `json:"create_ms_time"`
	OrderId      int64  `json:"order_id"`
	Desc         string `json:"desc"`
	Session      string `json:"session"`
	CustomData   string `json:"custom_data"`
	GroupId      int    `json:"group_id"`
	MapId        int    `json:"map_id"`
}

func OnMazeCardChangeProcess(c context.Context, logger fklog.FKLogI, index int, key, data []byte) error {
	defer fkprometheus.DebugPMT("OnMazeCardChangeProcess")()
	msg := new(MazeCardNotifyMsg)
	err := json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("OnMazeCardChangeProcess unmarshal err", zap.String("data", string(data)), zap.Error(err))
		return err
	}

	userId := msg.UserId
	logger.SetUid(userId)
	logger.InfoWF("OnMazeCardChangeProcess start", zap.Any("msg", msg))
	// 检查月卡状态
	expirationTime, err := userriddlemonthlyredis.GetMazeCardExpirationTime(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeCardChangeProcess GetMazeCardExpirationTime failed", zap.Error(err))
		return err
	}

	if expirationTime == 0 || expirationTime < time.Now().Unix() {
		logger.WarnWF("OnMazeCardChangeProcess is expiration", zap.Int64("expirationTime", expirationTime))
		return nil
	}

	return AddMazeCard(logger, userId, expirationTime)
}

func AddMazeCard(logger fklog.FKLogI, userId uint64, expirationTime int64) error {
	logger.DebugWF("AddMazeCard start", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	realBuff, showBuff := mazeconfigv8config.GetMazeConfig(100), mazeconfigv8config.GetMazeConfig(101)
	if len(realBuff) != 0 || len(showBuff) != 0 {
		// 添加月卡buff
		attrDb := &MazeBuffData.MazeBuffDb{
			MazeRealBuffs: PackMazeBuff(realBuff),
			MazeShowBuffs: PackMazeBuff(showBuff),
		}

		err := mazebuffinforedis.SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcMonthCard, attrDb)
		if err != nil {
			logger.ErrorWF("AddMazeCard SaveMazeBuffInfo failed", zap.Uint64("userId", userId), zap.Error(err))
			return err
		}

		// 推送属性变化通知
		msg := &structsdef.MazeCalcAttrNotifyMsg{
			UserId:     userId,
			FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
			BuffSrc:    constdef.MazeBuffSrcMonthCard,
			ChgType:    constdef.MazeBuffChgTypeCardOpen,
		}

		err = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, msg)
		if err != nil {
			logger.ErrorWF("AddMazeCard SendMazeAttrCalcNotify failed", zap.Any("msg", msg), zap.Error(err))
			return err
		}

		// 添加检查队列
		err = mazecardlistgroupredis.SetMazeCard(logger, userId, expirationTime)
		if err != nil {
			logger.ErrorWF("AddMazeCard SetMazeCard failed", zap.Uint64("userId", userId),
				zap.Int64("time", expirationTime), zap.Error(err))
			return err
		}
	}

	// 向客户端推送
	SendMazeCardMsg(logger, userId, constdef.MazeCardChgOpenType, expirationTime)
	logger.InfoWF("AddMazeCard end", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	return nil
}

func DeleteMazeCard(logger fklog.FKLogI, userId uint64) error {
	logger.DebugWF("DeleteMazeCard start", zap.Uint64("userId", userId))
	// 删除月卡buff
	err := mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcMonthCard)
	if err != nil {
		logger.ErrorWF("DeleteMazeCard DelMazeBuffBySrc failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	// 推送属性变化通知
	msg := &structsdef.MazeCalcAttrNotifyMsg{
		UserId:     userId,
		FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		BuffSrc:    constdef.MazeBuffSrcMonthCard,
		ChgType:    constdef.MazeBuffChgTypeCardExpiration,
	}

	// 删除检查队列
	err = mazecardlistgroupredis.DelMazeCard(logger, userId)
	if err != nil {
		logger.ErrorWF("DeleteMazeCard BatchDelMazeCard failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	err = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, msg)
	if err != nil {
		logger.ErrorWF("DeleteMazeCard SendMazeAttrCalcNotify failed", zap.Any("msg", msg), zap.Error(err))
		return err
	}

	// 向客户端推送
	SendMazeCardMsg(logger, userId, constdef.MazeCardChgExpirationType, 0)
	logger.InfoWF("DeleteMazeCard end", zap.Uint64("userId", userId))
	return nil
}

func SendMazeCardMsg(logger fklog.FKLogI, userId uint64, state int32, expirationTime int64) {
	msg := &MazeCard.MazeCardChangeID{
		State:          proto.Int32(state),
		ExpirationTime: proto.Int64(expirationTime),
	}
	_ = msg
	// TODO 为什么注释掉？？
	//_ = commonmustarriveredis.SendArrivePacketWithLogFix(logger, userId, 16200, msg)
}

func PackMazeBuff(buffMap map[int32]int64) []*MazeBuffData.MazeBuffAttr {
	if len(buffMap) == 0 {
		return nil
	}

	buffList := make([]*MazeBuffData.MazeBuffAttr, 0, len(buffMap))
	for id, value := range buffMap {
		buffList = append(buffList, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(value),
		})
	}

	return buffList
}
