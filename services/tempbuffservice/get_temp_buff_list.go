package tempbuffservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) GetMazeTempBuffList(ctx context.Context, userId uint64, barrierId int32) ([]*BuffInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "GetMazeTempBuffListRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}
	buffList := s.packShowBuffList(ctx, buffInfo)
	return buffList, nil
}

func (s *service) packShowBuffList(ctx context.Context, buffInfo *tempbuffmodel.TempBuffInfoModel) []*BuffInfo {
	if buffInfo == nil || len(buffInfo.SelectedBuff) == 0 {
		return nil
	}

	buffMap := make(map[int32]int64)
	var buffList []int32
	for _, info := range buffInfo.SelectedBuff {
		buffId := info.BuffId
		buffMap[buffId] += 1
		if buffMap[buffId] != 1 {
			continue
		}

		buffList = append(buffList, buffId)
	}

	showBuffList := make([]*BuffInfo, 0, len(buffMap))
	for _, id := range buffList {
		count := buffMap[id]
		if id == 0 || count == 0 {
			continue
		}

		// 拼接显示buff
		if showBuff := s.packShowBuff(ctx, id, count); showBuff != nil {
			showBuffList = append(showBuffList, showBuff)
		}
	}

	return showBuffList
}

func (s *service) packShowBuff(ctx context.Context, buffId int32, count int64) *BuffInfo {
	logger := fklog.ContextAppLogger(ctx)
	config := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, buffId)
	if config == nil {
		logger.CtxError(ctx, "packShowBuff buff is unknown", zap.Int32("buffId", buffId))
		return nil
	}

	info := &BuffInfo{
		BuffId: buffId,
		Value:  count,
		Name:   config.Affix_name,
	}

	if len(config.Affix_wildcard) == 0 {
		info.Decs = config.Affix_desc
		return info
	}

	strList := make([]interface{}, 0, len(config.Affix_wildcard))
	// 拼接属性
	for _, attrId := range config.Affix_wildcard {
		value, ok := config.Add_attr[attrId]
		if !ok {
			logger.CtxError(ctx, "packShowBuff add attr unknown", zap.Int32("buffId", buffId),
				zap.Int32("attrId", attrId))
			return nil
		}

		attrConfig := GMazeAttributeV8Cfg.GetWithCtx(ctx, attrId)
		if attrConfig == nil {
			logger.CtxError(ctx, "packShowBuff attr is unknown", zap.Int32("buffId", buffId),
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
			logger.CtxError(ctx, "packShowBuff buff Figure unknown", zap.Int32("buffId", buffId),
				zap.Int32("attrId", attrId), zap.Int32("figure", attrConfig.Figure))
			return nil
		}

		strList = append(strList, str)
	}

	info.Decs = fmt.Sprintf(config.Affix_desc, strList...)
	return info
}
