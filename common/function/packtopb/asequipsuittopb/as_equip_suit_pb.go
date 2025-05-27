/*
 * @Author: majian
 * @Date: 2024-08-23 11:30:55
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 11:26:41
 */
package asequipsuittopb

import (
	"maze_game_server/config/GMazeEquipSuiteAttrV8Cfg"
	"maze_game_server/pb/common/MazeGameEquip"

	"google.golang.org/protobuf/proto"
)

func PackAsEquipSuitInfo(suitKey int32) *MazeGameEquip.AsEquipSuitInfo {
	o := &MazeGameEquip.AsEquipSuitInfo{}
	row := GMazeEquipSuiteAttrV8Cfg.Get(suitKey)
	if row != nil {
		o.EpSuitId = proto.Int32(row.Suite_id)
	}
	return o
}
