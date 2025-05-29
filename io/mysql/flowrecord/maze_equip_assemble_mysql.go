package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/mysql"
)

const MazeEquipAssembleRecordTableName = "maze_equip_assemble_record"

func SaveEquipAssembleRecord(logger fklog.FKLogI, record *dollequipassmeblekakfa.MazeGameEquipAssembleRecord) {
	db, err := mysql.GetMysqlDb(MazeEquipAssembleRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipAssembleRecordTableName:", MazeEquipAssembleRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeEquipAssembleRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `equip_pos`, `op_type`, `new_equip_id`, `new_guid`, `old_equip_id`, `old_guid`, `old_f_elem`, `new_f_elem`, `ret_code`, `code_mask`, `trans_id`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.EquipPos, record.OpType, record.NewEquipId, record.NewGuid, record.OldEquipId, record.OldGuid,
		record.OldFElem, record.NewFElem, record.RetCode, record.CodeMask, record.TransID, record.GroupId, record.OpTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveEquipAssembleRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipAssembleRecord succ", zap.Any("flowrecord", record))

	return
}
