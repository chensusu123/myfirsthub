/*
 * @Author: majian
 * @Date: 2024-08-23 11:30:55
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 11:26:41
 */
package asequipsuittopb

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipSuiteAttrV8Cfg"
)

func PackAsEquipSuitInfo(suitKey int32) *MazeGameEquip.AsEquipSuitInfo {
	o := &MazeGameEquip.AsEquipSuitInfo{}
	row := GMazeEquipSuiteAttrV8Cfg.Get(suitKey)
	if row != nil {
		o.EpSuitId = proto.Int32(row.Suite_id)
	}
	return o
}
