package buff

import (
	"fmt"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuff"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuffSvr"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeattributeconfig"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixlvv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"
	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 15:15
 * @Description: 查询迷宫buff列表
 */

func GetMazeTempBuffListRQ(logger fknet.TCPContext, shardingID uint64, request, response proto.Message) error {
	defer fkprometheus.InfoPMT("GetMazeTempBuffListRQ")()
	start := time.Now()
	req := request.(*MazeTempBuff.GetMazeTempBuffListRQ)
	res := response.(*MazeTempBuff.GetMazeTempBuffListRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	defer func() {
		logger.InfoWF("GetMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId := shardingID, req.GetStageId()
	if userId == 0 || stageId == 0 {
		logger.WarnWF("GetMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}

	// todo 检查用户是不是小程序用户

	buffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("GetMazeTempBuffListRQ GetMazeTempBuff", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return nil
	}
	res.BuffList = packShowBuffList(logger, buffInfo)
	return nil
}

func packShowBuffList(logger fklog.FKLogI, buffInfo *MazeTempBuffSvr.TempBuffInfo) []*MazeTempBuff.MazeBuffInfo {
	if buffInfo == nil || len(buffInfo.SelectedBuff) == 0 {
		return nil
	}

	buffMap := make(map[int32]int64)
	var buffList []int32
	for _, info := range buffInfo.GetSelectedBuff() {
		buffId := info.GetBuffId()
		buffMap[buffId] += 1
		if buffMap[buffId] != 1 {
			continue
		}

		buffList = append(buffList, buffId)
	}

	showBuffList := make([]*MazeTempBuff.MazeBuffInfo, 0, len(buffMap))
	for _, id := range buffList {
		count := buffMap[id]
		if id == 0 || count == 0 {
			continue
		}

		// 拼接显示buff
		if showBuff := packShowBuff(logger, id, count); showBuff != nil {
			showBuffList = append(showBuffList, showBuff)
		}
	}

	return showBuffList
}

func packShowBuff(logger fklog.FKLogI, buffId int32, count int64) *MazeTempBuff.MazeBuffInfo {
	config := mazeenergyaffixlvv8config.GetAffixConfig(buffId)
	if config == nil {
		logger.WarnWF("packShowBuff buff is unknown", zap.Int32("buffId", buffId))
		return nil
	}

	info := &MazeTempBuff.MazeBuffInfo{
		BuffId: proto.Int32(buffId),
		Value:  proto.Int64(count),
		Name:   proto.String(config.Affix_name),
	}

	if len(config.Affix_wildcard) == 0 {
		info.Decs = proto.String(config.Affix_desc)
		return info
	}

	strList := make([]interface{}, 0, len(config.Affix_wildcard))
	// 拼接属性
	for _, attrId := range config.Affix_wildcard {
		value, ok := config.Add_attr[attrId]
		if !ok {
			logger.WarnWF("packShowBuff add attr unknown", zap.Int32("buffId", buffId),
				zap.Int32("attrId", attrId))
			return nil
		}

		attrConfig := mazeattributeconfig.GetMazeAttributeConfig(attrId)
		if attrConfig == nil {
			logger.WarnWF("packShowBuff attr is unknown", zap.Int32("buffId", buffId),
				zap.Int32("attrId", attrId))
			return nil
		}

		var str string
		totalValue := value * count
		switch attrConfig.Figure {
		case 1: // 固定值
			str = fmt.Sprintf("%d", totalValue)
		case 2: // 万分比
			str = fmt.Sprintf("%g%%", float64(totalValue)/100.0)
		case 3: // 白万分比
			str = fmt.Sprintf("%g%%", float64(totalValue)/10000.0)
		default:
			logger.WarnWF("packShowBuff buff Figure unknown", zap.Int32("buffId", buffId),
				zap.Int32("attrId", attrId), zap.Int32("figure", attrConfig.Figure))
			return nil
		}

		strList = append(strList, str)
	}

	info.Decs = proto.String(fmt.Sprintf(config.Affix_desc, strList...))
	return info
}
