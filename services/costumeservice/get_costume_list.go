package costumeservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/config/GMazeEquipTypeGroupResV8Cfg"
	"maze_game_server/excel/mazeequiptyperesv8"
	"maze_game_server/module/dollassembleinfo"
)

func (s service) GetUserCostume(logger fklog.FKLogI, userId uint64) (map[int32]int32, error) {
	res := make(map[int32]int32)
	assembleInfo, effect, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
	if err != nil {
		logger.ErrorWF("OnGetMazeAssembleRQ Get Assemble info fail", zap.Error(err))
		return nil, err
	}
	_ = effect
	var groupModel map[int32]int32 = nil
	allRows := GMazeEquipPosRankV8Cfg.GetAllMazeEquipPosRankV8Config()
	for _, posCfg := range allRows {
		var unlock int32
		for _, aEquip := range assembleInfo.MazeEquips {
			if posCfg.Pos_id == aEquip.GetEquipPos().GetPos() {
				equipId := aEquip.EquipInfo.GetEquipId()
				if aEquip.EquipInfo == nil {
					continue
				}
				if aEquip.EquipInfo.GetEquipId() == 0 {
					continue
				}
				equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, 0, 0)
				if equipTypeResCfg == nil {
					logger.ErrorWF("OnGetMazeAssembleRQ GetEquipTypeResCfg failed", zap.Int32("equipId", equipId))
					return nil, fmt.Errorf("装备id找不到")
				}
				if equipTypeResCfg.Maze_model_group_id != 0 {
					groupResCfg := GMazeEquipTypeGroupResV8Cfg.Get(equipTypeResCfg.Maze_model_group_id)
					if groupResCfg == nil {
						logger.ErrorWF("OnGetMazeAssembleRQ GMazeEquipTypeGroupResV8Cfg failed", zap.Int32("Maze_model_group_id", equipTypeResCfg.Maze_model_group_id))
						return nil, fmt.Errorf("配置找不到")
					}
					groupModel = groupResCfg.Maze_model_group
				} else {
					res[posCfg.Pos_id] = equipTypeResCfg.Maze_model
				}
				unlock = 1
				break
			}
		}
		if unlock == 1 {
			continue
		}
		res[posCfg.Pos_id] = 0
	}
	if groupModel != nil {
		for k, v := range groupModel {
			res[k] = v
		}
	}
	return res, nil
}
