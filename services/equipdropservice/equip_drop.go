package equipdropservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/randfuncs"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeEquDropV8Cfg"
	"maze_game_server/config/GMazeEquipPosQuaLvV8Cfg"
	"maze_game_server/model/equipdropmodel"
	"maze_game_server/module/mazeuserinfo"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s service) GetNewEquip(ctx context.Context, userId uint64, level int32, barrierId int32, equipNum int32) (newEquip map[int32]int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "GetNewEquip start", zap.Uint64("userId", userId), zap.Int32("level", level), zap.Int32("barrierId", barrierId), zap.Int32("equipNum", equipNum))
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetNewEquip GetNewEquip fail", zap.Error(err))
		return
	}

	if level <= 0 {
		level = int32(userInfo.Level)
	}

	if barrierId <= 0 {
		barrierId = userInfo.Barrier
	}

	newEquip, err = s.getEquipId(ctx, userId, level, barrierId, equipNum)
	if err != nil {
		return nil, err
	}
	logger.CtxInfo(ctx, "GetNewEquip success", zap.Uint64("userId", userId), zap.Any("newEquip", newEquip))
	return
}

func (s service) GetMazeEquipSpecialDropInfo(ctx context.Context, userId uint64) (*equipdropmodel.EquipSpecialDropModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	info, err := equipdropmodel.NewEquipSpecialDropModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetMazeEquipSpecialDropInfo GetMazeEquipSpecialDropInfo failed", zap.Error(err))
		return nil, errors.New("获取用户buff信息失败")
	}
	return info, nil
}

func (s service) getEquipId(ctx context.Context, userId uint64, mazeLevel int32, barrier int32, addCount int32) (equipMap map[int32]int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	if addCount <= 0 {
		logger.CtxInfo(ctx, "getEquipId addCount <= 0", zap.Any("addCount", addCount))
		return
	}

	info, err := s.GetMazeEquipSpecialDropInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetEquipId GetMazeEquipSpecialDropInfo failed", zap.Error(err), zap.Uint64("userId", userId))
		return nil, err
	}

	cfg, err := getMazeEquDropV8Cfg(ctx, barrier)
	if err != nil {
		logger.CtxInfo(ctx, "GetEquipId GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil, err
	}
	index, ok := info.DropMap[cfg.Order]
	if !ok {
		index = 0
	}

	equipMap = make(map[int32]int32)
	if cfg.Special_drop != nil && index < int32(len(cfg.Special_drop)-1) {
		dropMap, addNum := s.specialEquipDrop(ctx, userId, barrier, mazeLevel, addCount)
		if len(dropMap) > 0 {
			for k, v := range dropMap {
				equipMap[k] += v
			}
		}
		addCount -= addNum
	}

	if addCount <= 0 {
		return
	}
	dropMap := s.regularityEquipDrop(ctx, userId, barrier, mazeLevel, addCount)
	if len(dropMap) > 0 {
		for k, v := range dropMap {
			equipMap[k] += v
		}
	}
	return
}

// 特殊掉落
func (s service) specialEquipDrop(ctx context.Context, userId uint64, barrier, mazeLevel, addCount int32) (equipIdMap map[int32]int32, add int32) {
	logger := fklog.ContextAppLogger(ctx)
	cfg, err := getMazeEquDropV8Cfg(ctx, barrier)
	if err != nil {
		logger.CtxInfo(ctx, "specialEquipDrop GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil, 0
	}

	if cfg.Special_drop == nil || len(cfg.Special_drop) == 0 {
		logger.CtxInfo(ctx, "specialEquipDrop GetEquipId special_drop empty")
		return nil, 0
	}

	info, err := s.GetMazeEquipSpecialDropInfo(ctx, userId)
	if err != nil {
		logger.CtxInfo(ctx, "specialEquipDrop GetMazeEquipSpecialDropInfo failed", zap.Error(err), zap.Uint64("userId", userId))
		return nil, 0
	}

	index, ok := info.DropMap[cfg.Order]
	if ok {
		if index >= int32(len(cfg.Special_drop)-1) {
			logger.CtxInfo(ctx, "specialEquipDrop special_drop max limit")
			return nil, 0
		}
	} else {
		index = 0
	}

	equipIdMap = make(map[int32]int32)
	newLevel := s.GetMazeBarrierLv(ctx, mazeLevel, barrier)
	add = int32(0)
	for i := 0; i < len(cfg.Special_drop); i += 2 {
		if addCount == add {
			break
		}
		if i < int(index) {
			continue
		}

		quality := cfg.Special_drop[i]
		pos := cfg.Special_drop[i+1]
		equipId := getEquipId(newLevel, quality, pos)
		equipIdMap[equipId] += 1

		index = int32(i + 1)
		info.DropMap[cfg.Order] = int32(i + 1)
		add++
		logger.CtxInfo(ctx, "specialEquipDrop getEquipId success", zap.Int32("level", newLevel), zap.Int32("quality", quality), zap.Int32("pos", pos), zap.Int32("equipId", equipId), zap.Any("equipIdMap", equipIdMap), zap.Any("info.DropMap", info.DropMap))

	}

	if add > 0 {
		err = info.Save(ctx, userId)
		if err != nil {
			logger.CtxError(ctx, "specialEquipDrop EquipSpecialDropModel save failed", zap.Uint64("userId", userId), zap.Error(err))
			return equipIdMap, add
		}
	}

	return
}

// 常规掉落组
func (s service) regularityEquipDrop(ctx context.Context, userId uint64, barrier, mazeLevel, addCount int32) (equipIdMap map[int32]int32) {
	logger := fklog.ContextAppLogger(ctx)
	cfg, err := getMazeEquDropV8Cfg(ctx, barrier)
	if err != nil {
		logger.CtxInfo(ctx, "regularityEquipDrop GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil
	}

	if cfg.Regularity_drop == nil || len(cfg.Regularity_drop) == 0 {
		logger.CtxInfo(ctx, "regularityEquipDrop GetEquipId regularity_drop empty")
		return nil
	}

	//权重
	seqWeight := make(map[int32]int32)
	for quality, weight := range cfg.Regularity_drop {
		seqWeight[quality] = weight
	}
	equipIdMap = make(map[int32]int32)
	newLevel := s.GetMazeBarrierLv(ctx, mazeLevel, barrier)
	for i := int32(0); i < addCount; i++ {
		quality := randfuncs.RandByWeightV2(logger, seqWeight, false)
		if quality <= 0 {
			logger.CtxInfo(ctx, "regularityEquipDrop MazeEquDropV8Cfg err", zap.Any("seqWeight", seqWeight))
			return nil
		}

		pos := fkutil.RandInt32(1, constdef.EquipPosNum+1) //部位等概率随机

		equipId := getEquipId(newLevel, quality, int32(pos))
		logger.CtxInfo(ctx, "regularityEquipDrop getEquipId success", zap.Int32("level", newLevel), zap.Int32("quality", quality), zap.Int32("pos", int32(pos)), zap.Int32("equipId", equipId))
		equipIdMap[equipId] += 1
	}
	return
}

func (s service) GetMazeBarrierLv(ctx context.Context, level int32, barrier int32) int32 {
	cfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrier)
	if cfg == nil {
		return level
	}
	if level < cfg.Drop_equip_lv_min {
		return cfg.Drop_equip_lv_min
	}
	if level > cfg.Drop_equip_lv_max {
		return cfg.Drop_equip_lv_max
	}
	return level
}

func getMazeEquDropV8Cfg(ctx context.Context, barrier int32) (*GMazeEquDropV8Cfg.MazeEquDropV8ConfigRow, error) {
	cfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrier)
	if cfg == nil {
		return nil, errors.New("GMazeBarriesV8Cfg config error")
	}
	eCfg := GMazeEquDropV8Cfg.GetWithCtx(ctx, cfg.Equ_drop)
	if eCfg == nil {
		return nil, errors.New("GMazeEquDropV8Cfg config error")
	}
	return eCfg, nil
}

func getEquipId(level, quality, pos int32) int32 {
	dataList := GMazeEquipPosQuaLvV8Cfg.GetAll()
	for _, row := range dataList {
		if row.Level == level && row.Quality == quality {
			switch pos {
			case 1:
				return row.Pos1
			case 2:
				return row.Pos2
			case 3:
				return row.Pos3
			case 4:
				return row.Pos4
			case 5:
				return row.Pos5
			case 6:
				return row.Pos6
			case 7:
				return row.Pos7
			case 8:
				return row.Pos8
			}
			break
		}
	}
	return 0
}

func (s service) GmDelete(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	info := &equipdropmodel.EquipSpecialDropModel{}
	err := info.Del(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GlobalEquipDropService GmDelete err", zap.Error(err), zap.Uint64("userId", userId))
		return err
	}
	logger.CtxInfo(ctx, "GlobalEquipDropService GmDelete success", zap.Error(err), zap.Uint64("userId", userId))
	return nil
}
