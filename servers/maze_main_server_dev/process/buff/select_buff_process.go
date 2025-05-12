package buff

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixlvv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazetempbuffchgmsg"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeTempBuff"
	"gitlab.ifreetalk.com/plate/protodef/MazeTempBuffSvr"
	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 15:48
 * @Description: 选择迷宫buff
 */

func SelectMazeTempBuffRQ(logger fknet.TCPContext, shardingID uint64, request, response proto.Message) error {
	defer fkprometheus.InfoPMT("SelectMazeTempBuffRQ")()
	start := time.Now()
	req := request.(*MazeTempBuff.SelectMazeTempBuffRQ)
	res := response.(*MazeTempBuff.SelectMazeTempBuffRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	defer func() {
		logger.InfoWF("SelectMazeTempBuffRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, buffId := shardingID, req.GetStageId(), req.GetLevel(), req.GetBuffId()
	if userId == 0 || stageId == 0 || level == 0 || buffId == 0 {
		logger.WarnWF("SelectMazeTempBuffRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
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

	err = updateBuffInfo(logger, userId, stageId, level, buffId, buffInfo)
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

func updateBuffInfo(logger fklog.FKLogI, userId uint64, stageId, level, buffId int32,
	buffInfo *MazeTempBuffSvr.TempBuffInfo) error {
	buffInfo.BuffSequence = &MazeTempBuffSvr.BuffSequence{
		Index: proto.Int32(level),
	}

	buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &MazeTempBuffSvr.SelectedBuffInfo{
		BuffId: proto.Int32(buffId),
		Level:  proto.Int32(level),
	})

	var totalMap map[int32]int64
	totalMap, buffInfo.TotalBuff = getTotalBuff(logger, buffInfo.GetSelectedBuff())
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
	return nil
}

func getTotalBuff(logger fklog.FKLogI, buffList []*MazeTempBuffSvr.SelectedBuffInfo) (
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
