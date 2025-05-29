package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/mysql"
)

const MazeEquipBagRecordTableName = "maze_equip_bag_record"

// 保存背包流水
func SaveEquipBagRecord(logger fklog.FKLogI, record *mazeequipbagrecord.MazeGameEquipBagRecord) {
	db, err := mysql.GetMysqlDb(MazeEquipBagRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipBagRecordTableName:", MazeEquipBagRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeEquipBagRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `chg_type`, `trade_num`, `add_equip_guids`, `del_equip_guids`, `op_type`, `is_fail`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.ChgType, record.TradeNum, record.AddEquipGuids, record.DelEquipGuids, record.OpType,
		record.IsFail, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveEquipBagRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipBagRecord succ", zap.Any("flowrecord", record))

	return
}
