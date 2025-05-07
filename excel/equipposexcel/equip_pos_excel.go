package equipposexcel

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosLvV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/excelutil"
)

func GetPosStrengthCfg(posId, lv int32) (row *GMazeEquipPosLvV8Cfg.MazeEquipPosLvV8ConfigRow) {
	return GetPosStrengthCfgByKey(excelutil.GetEquipPosEnLevelKey(posId, lv))
}

func GetPosStrengthCfgByKey(key int32) (row *GMazeEquipPosLvV8Cfg.MazeEquipPosLvV8ConfigRow) {
	return GMazeEquipPosLvV8Cfg.GetMazeEquipPosLvV8Config(key)
}
