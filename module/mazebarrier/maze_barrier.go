package mazebarrier

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeBoxV8Cfg"
	"go.uber.org/zap"
)

func GetBarrierPassAward(logger fklog.FKLogI, barrierId int32) (awardMap map[int32]int64, equipMap map[int32]int32, err error) {
	cfg := GMazeBarriesV8Cfg.Get(barrierId)
	if cfg == nil {
		logger.ErrorWF("GetBarrierPassAward get barrier cfg fail", zap.Any("barrier", barrierId))
		err = errors.New("barrier cfg nil")
		return
	}

	boxAward := GMazeBoxV8Cfg.Get(cfg.Box_id)
	if cfg == nil {
		logger.ErrorWF("GetBarrierPassAward get box cfg fail", zap.Any("box", cfg.Box_id))
		err = errors.New("barrier cfg nil")
		return
	}

	awardMap = make(map[int32]int64)
	equipMap = make(map[int32]int32)
	for _, e := range boxAward.Award_equip {
		if e > 0 {
			equipMap[e] += 1
		}
	}

	for k, v := range boxAward.Drop_item {
		if k > 0 && v > 0 {
			awardMap[k] += v
		}
	}

	return
}

func GetBarrierPassAwardWithFirst(logger fklog.FKLogI, barrierId int32) (awardMap map[int32]int64, equipMap map[int32]int32, err error) {
	cfg := GMazeBarriesV8Cfg.Get(barrierId)
	if cfg == nil {
		logger.ErrorWF("GetBarrierPassAwardWithFirst get barrier cfg fail", zap.Any("barrier", barrierId))
		err = errors.New("barrier cfg nil")
		return
	}

	boxAward := GMazeBoxV8Cfg.Get(cfg.Box_id_first)
	if cfg == nil {
		logger.ErrorWF("GetBarrierPassAwardWithFirst get box cfg fail", zap.Any("box", cfg.Box_id_first))
		err = errors.New("barrier cfg nil")
		return
	}

	awardMap, equipMap, err = GetBarrierPassAward(logger, barrierId)
	if err != nil {
		return
	}

	for _, e := range boxAward.Award_equip {
		if e > 0 {
			equipMap[e] += 1
		}
	}

	for k, v := range boxAward.Drop_item {
		if k > 0 && v > 0 {
			awardMap[k] += v
		}
	}

	return
}
