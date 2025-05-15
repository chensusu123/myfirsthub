/*
 * @Author: majian
 * @Date: 2024-12-16 14:54:29
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-16 16:36:42
 */
package copyequipgm

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm/copyinterface"
)

func RunCopyTask(logger fklog.FKLogI, srcUserId uint64, dstUserId []uint64) bool {
	var runParam copyinterface.CopyParam
	copyinterface.RangeAiCopy(logger, srcUserId, dstUserId, runParam)
	return true
}

func init() {
	copyinterface.RegistHandler("bagEquipData", CopyBagEquipData, 1)
	// copyinterface.RegistHandler("equipMagicData", CopyMagicData, 2)
	copyinterface.RegistHandler("assembleData", CopyAssembleData, 3)
}
