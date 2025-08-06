package equipdropservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/randfuncs"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeEquDropV8Cfg"
	"maze_game_server/config/GMazeEquipPosQuaLvV8Cfg"
	"maze_game_server/model/equipdropmodel"
	"maze_game_server/module/mazeuserinfo"
)

func (s service) GetNewEquip(logger fklog.FKLogI, userId uint64, barrierId int32, level int32, equipNum int32) (newEquip map[int32]int32, err error) {
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("GetNewEquip GetNewEquip fail", zap.Error(err))
		return
	}

	if level <= 0 {
		level = int32(userInfo.Level)
	}

	if barrierId <= 0 {
		barrierId = userInfo.Barrier
	}

	newEquip, err = s.getEquipId(logger, userId, level, barrierId, equipNum)
	if err != nil {
		return nil, err
	}
	return
}

func (s service) GetMazeEquipSpecialDropInfo(logger fklog.FKLogI, userId uint64) (*equipdropmodel.EquipSpecialDropModel, error) {
	info, err := equipdropmodel.NewEquipSpecialDropModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMazeEquipSpecialDropInfo GetMazeEquipSpecialDropInfo failed", zap.Error(err))
		return nil, errors.New("获取用户buff信息失败")
	}
	return info, nil
}

func (s service) getEquipId(logger fklog.FKLogI, userId uint64, mazeLevel int32, barrier int32, addCount int32) (equipMap map[int32]int32, err error) {
	if addCount <= 0 {
		logger.InfoWF("getEquipId addCount <= 0", zap.Any("addCount", addCount))
		return
	}

	info, err := s.GetMazeEquipSpecialDropInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("GetEquipId GetMazeEquipSpecialDropInfo failed", zap.Error(err), zap.Uint64("userId", userId))
		return nil, err
	}

	cfg, err := getMazeEquDropV8Cfg(barrier)
	if err != nil {
		logger.ErrorWF("GetEquipId GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil, err
	}
	equipMap = make(map[int32]int32)
	if cfg.Special_drop != nil && info.SpecialDropIndex < int32(len(cfg.Special_drop)-2) {
		dropMap, addNum := s.specialEquipDrop(logger, userId, barrier, mazeLevel, addCount)
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
	for i := int32(0); i < addCount; i++ {
		dropMap := s.regularityEquipDrop(logger, userId, barrier, mazeLevel, addCount)
		if len(dropMap) > 0 {
			for k, v := range dropMap {
				equipMap[k] += v
			}
		}
	}
	return
}

// 特殊掉落
func (s service) specialEquipDrop(logger fklog.FKLogI, userId uint64, barrier, mazeLevel, addCount int32) (equipIdMap map[int32]int32, add int32) {

	cfg, err := getMazeEquDropV8Cfg(barrier)
	if err != nil {
		logger.ErrorWF("specialEquipDrop GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil, 0
	}

	if cfg.Special_drop == nil || len(cfg.Special_drop) == 0 {
		logger.ErrorWF("specialEquipDrop GetEquipId special_drop empty")
		return nil, 0
	}

	info, err := s.GetMazeEquipSpecialDropInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("specialEquipDrop GetMazeEquipSpecialDropInfo failed", zap.Error(err), zap.Uint64("userId", userId))
		return nil, 0
	}

	if info.SpecialDropIndex >= int32(len(cfg.Special_drop)-2) {
		logger.InfoWF("specialEquipDrop special_drop max limit")
		return nil, 0
	}

	equipIdMap = make(map[int32]int32)
	newLevel := s.GetMazeBarrierLv(mazeLevel, barrier)
	add = int32(0)
	for i := 0; i < len(cfg.Special_drop); i += 2 {
		if addCount == add {
			break
		}
		if i < int(info.SpecialDropIndex) {
			continue
		}

		quality := cfg.Special_drop[i]
		pos := cfg.Special_drop[i+1]
		equipId := getEquipId(newLevel, quality, pos)
		equipIdMap[equipId] += 1

		info.SpecialDropIndex = int32(i)
		add++
	}

	if add > 0 {
		err = info.Save(logger, userId)
		if err != nil {
			logger.ErrorWF("specialEquipDrop EquipSpecialDropModel save failed", zap.Uint64("userId", userId), zap.Error(err))
			return equipIdMap, add
		}
	}

	return
}

// 常规掉落组
func (s service) regularityEquipDrop(logger fklog.FKLogI, userId uint64, barrier, mazeLevel, addCount int32) (equipIdMap map[int32]int32) {

	cfg, err := getMazeEquDropV8Cfg(barrier)
	if err != nil {
		logger.ErrorWF("regularityEquipDrop GetMazeEquDropV8Cfg failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("barrier", barrier))
		return nil
	}

	if cfg.Regularity_drop == nil || len(cfg.Regularity_drop) == 0 {
		logger.ErrorWF("regularityEquipDrop GetEquipId regularity_drop empty")
		return nil
	}

	//权重
	seqWeight := make(map[int32]int32)
	for quality, weight := range cfg.Regularity_drop {
		seqWeight[quality] = weight
	}
	equipIdMap = make(map[int32]int32)
	newLevel := s.GetMazeBarrierLv(mazeLevel, barrier)
	for i := int32(0); i < addCount; i++ {
		quality := randfuncs.RandByWeightV2(logger, seqWeight, false)
		if quality <= 0 {
			logger.ErrorWF("regularityEquipDrop MazeEquDropV8Cfg err", zap.Any("seqWeight", seqWeight))
			return nil
		}

		pos := fkutil.RandInt32(1, constdef.EquipPosNum+1) //部位等概率随机
		logger.InfoWF("regularityEquipDrop random quality, pos", zap.Int32("level", newLevel), zap.Int32("quality", quality), zap.Int("pos", pos))

		equipId := getEquipId(newLevel, quality, int32(pos))
		equipIdMap[equipId] += 1
	}
	return
}

// 编码特殊掉落记录
//func EnCodeRecordId(quality, pos int32) int32 {
//	return quality*equipSpecialDropRate + pos
//}

// 解码特殊掉落记录
//func DecodeRecordId(recordId int32) (quality, pos int32) {
//	quality = recordId / equipSpecialDropRate
//	pos = recordId % equipSpecialDropRate
//	return
//}

func (s service) GetMazeBarrierLv(level int32, barrier int32) int32 {
	cfg := GMazeBarriesV8Cfg.Get(barrier)
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

func getMazeEquDropV8Cfg(barrier int32) (*GMazeEquDropV8Cfg.MazeEquDropV8ConfigRow, error) {
	cfg := GMazeBarriesV8Cfg.Get(barrier)
	if cfg == nil {
		return nil, errors.New("GMazeBarriesV8Cfg config error")
	}
	eCfg := GMazeEquDropV8Cfg.Get(cfg.Equ_drop)
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
