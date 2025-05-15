/*
 * @Author: majian
 * @Date: 2024-08-22 16:58:32
 * @Last Modified by: majian
 * @Last Modified time: 2025-01-08 19:15:56
 */
package calcassembleattr

import (
	"errors"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipSuiteInfoV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"go.uber.org/zap"
)

func CalcEquipSuit(logger fklog.FKLogI, equips []*MazeEquipCache.MazeEquipPosInfo) (suitMgr *EquipSuitMgr, err error) {
	suitMgr = NewEquipSuitMgr()
	for _, equipPos := range equips {
		if !assemble.IsAssembleEquip(equipPos) {
			continue
		}
		if equipPos.GetEquipInfo().GetSuitId() > 0 {
			suitId := equipPos.GetEquipInfo().GetSuitId()
			suitInfo := suitMgr.GetSuitInfo(suitId)
			if suitInfo == nil {
				desV8Row := GMazeEquipSuiteInfoV8Cfg.GetMazeEquipSuiteInfoV8Config(suitId)
				if desV8Row == nil {
					logger.ErrorWF("CalcEquipSuit cannot find cfg", zap.String("sheet", GMazeEquipSuiteInfoV8Cfg.GetConfigDesc()))
					return nil, errors.New("cannot find suit cfg")
				}
				suitInfo = suitMgr.InitSuit(suitId, desV8Row.Max_num)
			}
			suitInfo.AddActivePos(equipPos.GetEquipPos().GetPos(), equipPos.GetEquipLoadInfo().GetEquipGuid())
		}
	}
	return suitMgr, nil
}
