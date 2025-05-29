package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/mysql"
)

const MazeEquipDismantRecordTableName = "maze_equip_dismant_record"

// 保存装备分解流水
func SaveEquipDismantRecord(logger fklog.FKLogI, record *dollequipdismantlekafka.MazeGameEquipDismantleRecord) {
	db, err := mysql.GetMysqlDb(MazeEquipDismantRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipDismantRecordTableName:", MazeEquipDismantRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeEquipDismantRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `equip_guids`, `trade_num`, `award`, `op_type`, `is_fail`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.EquipGuids, record.TradeNum, record.Award, record.OpType, record.IsFail, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipDismantRecord succ", zap.Any("flowrecord", record))

	return
}
