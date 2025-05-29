package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/mysql"
)

const MazeMoneyRecordTableName = "maze_money_record"

// 保存用户货币变化流水
func SaveMoneyRecord(logger fklog.FKLogI, record *mazemoneykafka.MazeMoneyRecord) {
	db, err := mysql.GetMysqlDb(MazeMoneyRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeMoneyRecordTableName:", MazeMoneyRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeMoneyRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `old_money_id`, `old_money_count`, `new_money_id`, `new_money_count`, `trade_no`, `chg_reason`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.OldMoneyId, record.OldMoneyCount, record.NewMoneyId, record.NewMoneyCount, record.TradeNo, record.ChgReason,
		record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveMoneyRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveMoneyRecord succ", zap.Any("flowrecord", record))

	return
}
