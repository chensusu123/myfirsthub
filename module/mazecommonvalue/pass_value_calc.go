package mazecommonvalue

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/config/GDollMapPuzzleNewV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeMapEditorConfigIdV8Cfg"
	"maze_game_server/pb/common/MazeGame"
	"strconv"
)

// 计算通关值
func CalcPassValue(logger fklog.FKLogI, barrier, stage int32) (int64, error) {
	passValue, err := CalcInitPassValue(logger, barrier)
	if err != nil {
		return 0, err
	}
	if stage > 0 {
		// 计算本关已通过阶段的通关值
		for _, i := range GDollMapPuzzleNewV8Cfg.GetAll() {
			if i.Level != barrier || i.Stage > stage || i.Stage == 0 {
				continue
			}
			if i.Config_id == "" {
				continue
			}
			for _, j := range GMazeMapEditorConfigIdV8Cfg.GetAll() {
				configId, err := strconv.ParseInt(i.Config_id, 10, 32)
				if err != nil {
					logger.ErrorWF("calcPassValue Parse config id err", zap.String("configId", i.Config_id))
					return 0, err
				}
				if j.Level_id == i.Level && j.Config_id == int32(configId) {
					if j.Add_kungfu != 0 {
						passValue += int64(j.Add_kungfu)
					}
				}
			}
		}
	}

	return passValue, nil
}

func CalcInitPassValue(logger fklog.FKLogI, barrier int32) (int64, error) {
	cfg := GMazeConfigV8Cfg.Get(constdef.PassValueInitCfgId)
	if cfg == nil {
		logger.ErrorWF("calcPassValue PassValueInitCfgId not exist", zap.Int32("cfgId", constdef.PassValueInitCfgId))
		return 0, fmt.Errorf("%d config not exist", constdef.PassValueInitCfgId)
	}
	passValue := cfg.Value_int
	// 计算已经通过关卡的通关值
	for _, row := range GMazeBarriesV8Cfg.GetAll() {
		if row.Order < barrier {
			passValue += int64(row.Barries_add_kongfu)
		}
	}
	return passValue, nil
}

// 发送通关值变化id包
func SendPassValueIdPack(logger fklog.FKLogI, userId uint64, barrierId, stage int32) error {
	passValue, err := CalcPassValue(logger, barrierId, stage)
	if err != nil {

		return err
	}
	commonList := []*CommonValueStruct{
		&CommonValueStruct{
			DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_PASS_VALUE),
			DataValueInt: passValue,
			// ChgReason:    int32(1),
		},
	}
	err = SendCommonValueIdPack(logger, userId, commonList)
	if err != nil {
		logger.ErrorWF("SendCommonValueIdPack err", zap.Error(err))
	}
	return nil
}
