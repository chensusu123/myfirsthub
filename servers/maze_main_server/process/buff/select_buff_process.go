package buff

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/excel/mazeenergyaffixlvv8config"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazetempbuffredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 15:48
 * @Description: 选择迷宫buff
 */

func (b *Buff) SelectMazeTempBuffRQ_10437_10438(s *session.Session, req *MazeTempBuff.SelectMazeTempBuffRQ) (err error) {
	defer fkprometheus.InfoPMT("SelectMazeTempBuffRQ")()

	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	res := &MazeTempBuff.SelectMazeTempBuffRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.Type = req.Type
	defer func() {
		err = s.Response(res)
		logger.InfoWF("SelectMazeTempBuffRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, buffId, buffType := uint64(s.UID()), req.GetStageId(), req.GetLevel(), req.GetBuffId(), int32(req.GetType())
	if userId == 0 || stageId == 0 || level == 0 || buffId == 0 {
		logger.WarnWF("SelectMazeTempBuffRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}
	if buffType != int32(MazeTempBuff.Type_UP_LEVEL) && buffType != int32(MazeTempBuff.Type_USE_ITEM) {
		logger.ErrorWF("SelectMazeTempBuffRQ buffType args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff类型参数错误")
		return nil
	}

	// todo 检查用户是不是小程序用户

	buffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("SelectMazeTempBuffRQ GetMazeTempBuff", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return nil
	}

	if buffInfo == nil {
		logger.WarnWF("SelectMazeTempBuffRQ buff is nil", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return nil
	}

	// 检查是否可选，选的buff是否为当次可选buff
	if !checkSelectBuff(logger, level, buffId, buffInfo) {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff信息异常")
		return nil
	}

	err = updateBuffInfo(logger, userId, stageId, level, buffId, buffType, buffInfo)
	if err != nil {
		logger.ErrorWF("SelectMazeTempBuffRQ updateBuffInfo failed", zap.Error(err))
		return err
	}

	res.BuffList = packShowBuffList(logger, buffInfo)
	return nil
}

func checkSelectBuff(logger fklog.FKLogI, level, buffId int32, buffInfo *MazeTempBuffSvr.TempBuffInfo) bool {
	// 检查是否可选，选的buff是否为当次可选buff
	if level != buffInfo.GetBuffSequence().GetIndex() {
		logger.WarnWF("checkSelectBuff level unknown", zap.Int32("level", level),
			zap.Int32("needLevel", buffInfo.GetBuffSequence().GetIndex()))
		return false
	}

	for _, id := range buffInfo.GetBuffSequence().GetSelectBuffList() {
		if buffId == id {
			return true
		}
	}

	logger.WarnWF("checkSelectBuff buffId unknown", zap.Int32("buffId", buffId),
		zap.Int32s("buffList", buffInfo.GetBuffSequence().GetSelectBuffList()))
	return false
}

func updateBuffInfo(logger fklog.FKLogI, userId uint64, stageId, level, buffId, buffType int32,
	buffInfo *MazeTempBuffSvr.TempBuffInfo) error {
	areaId := buffInfo.BuffSequence.GetAreaId()
	areaIndex := buffInfo.BuffSequence.GetAreaIndex()
	buffInfo.BuffSequence = &MazeTempBuffSvr.BuffSequence{
		Index: proto.Int32(level),
	}

	buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &MazeTempBuffSvr.SelectedBuffInfo{
		BuffId:    proto.Int32(buffId),
		Level:     proto.Int32(level),
		Type:      proto.Int32(buffType),
		AreaId:    proto.Int32(areaId),
		AreaIndex: proto.Int32(areaIndex),
	})

	var totalMap map[int32]int64
	totalMap, buffInfo.TotalBuff = GetTotalBuff(logger, buffInfo.GetSelectedBuff())
	// 更新buff信息
	err := mazetempbuffredis.SetMazeTempBuff(logger, userId, stageId, buffInfo)
	if err != nil {
		logger.ErrorWF("updateBuffInfo SetMazeTempBuff failed", zap.Any("info", buffInfo), zap.Error(err))
		return err
	}

	// 推送buff变化信息
	msg := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
		UserId:  userId,
		StageId: stageId,
		ChgType: 1,
		ChgDesc: "选择buff",
	}

	// 计算buff变化
	config := mazeenergyaffixlvv8config.GetAffixConfig(buffId)
	if config != nil {
		chgAttrs := make([]*mazetempbuffchgmsg.AttrChgInfo, 0, len(config.Add_attr))
		for id, value := range config.Add_attr {
			chgAttrs = append(chgAttrs, &mazetempbuffchgmsg.AttrChgInfo{
				AttrId: id,
				OldVal: totalMap[id] - value,
				CurVal: totalMap[id],
			})
		}

		msg.ChgAttrs = chgAttrs
	}

	_ = mazetempbuffchgmsg.PushTempBuffChangeMsg(logger, msg)

	// 同步到buff中心
	TempBuffChangeSync(logger, userId, buffInfo)

	return nil
}

// 同步到buff中心
func TempBuffChangeSync(logger fklog.FKLogI, userId uint64, buffInfo *MazeTempBuffSvr.TempBuffInfo) error {
	forceAttr, err := GetSelectBuffForceAttr(buffInfo.TotalBuff)
	if err != nil {
		logger.ErrorWF("updateBuffInfo GetSelectBuffForceAttr failed", zap.Error(err))
		return err
	}
	attrDb := &MazeBuffData.MazeBuffDb{
		MazeRealBuffs: PackMazeBuff(forceAttr),
		// MazeShowBuffs: PackMazeBuff(showBuff), // todo 现在暂时没有展示武力值
	}

	err = mazebuffinforedis.SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcSelectBuffForce, attrDb)
	if err != nil {
		logger.ErrorWF("AddMazeCard SaveMazeBuffInfo failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	// 推送属性计算消息
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		UserId: userId,
		// FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		ChgType: constdef.MazeBuffChgForceValue,
		Session: "buff",
		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	}
	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	return nil
}

func GetTotalBuff(logger fklog.FKLogI, buffList []*MazeTempBuffSvr.SelectedBuffInfo) (
	map[int32]int64, []*MazeTempBuffSvr.TotalBuffInfo) {
	totalMap := make(map[int32]int64)
	for _, info := range buffList {
		// 获取buff实际加成
		config := mazeenergyaffixlvv8config.GetAffixConfig(info.GetBuffId())
		if config != nil {
			for id, value := range config.Add_attr {
				totalMap[id] += value
			}
		}
	}

	totalList := make([]*MazeTempBuffSvr.TotalBuffInfo, 0, len(totalMap))
	for id, value := range totalMap {
		totalList = append(totalList, &MazeTempBuffSvr.TotalBuffInfo{
			BuffId:    proto.Int32(id),
			BuffValue: proto.Int64(value),
		})
	}

	return totalMap, totalList
}

func GetSelectBuffForceAttr(buffInfo []*MazeTempBuffSvr.TotalBuffInfo) (forceAttr map[int32]int64, err error) {
	if buffInfo == nil {
		return
	}
	forceAttr = make(map[int32]int64)
	for _, v := range buffInfo {
		// 读属性表
		cfg := GMazeAttributeV8Cfg.Get(v.GetBuffId())
		if cfg == nil {
			err = fmt.Errorf("GetSelectBuffForceAttr GetMazeAttributeFormulaV8Cfg nil, attrID: %d", v.GetBuffId())
			return
		}
		if cfg.Type == forceAttrType {
			forceAttr[cfg.Id] += v.GetBuffValue()
		}
	}
	return
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

const forceAttrType = 5
