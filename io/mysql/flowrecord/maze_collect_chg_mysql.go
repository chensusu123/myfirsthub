package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/mysql"
)

const MazeCollectChgRecordTableName = "maze_collect_chg_record"

// 迷宫挂机变化流水
func SaveCollectChgRecord(logger fklog.FKLogI, record *mazecollectrecord.MazeCollectChgRecord) {
	db, err := mysql.GetMysqlDb(MazeCollectChgRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeCollectChgRecordTableName:", MazeCollectChgRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeCollectChgRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `op_type`, `start_time`, `last_time`, `new_last_time`, `available_time`, `end_time`, `period_time`, `collect_times`, `barrier_id`, `trade_no`, `add_items`, `remain_items`, `ret_code`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.OpType, record.StartTime, record.LastTime, record.NewLastTime, record.AvailableTime, record.EndTime,
		record.PeriodTime, record.CollectTimes, record.BarrierId, record.TradeNo, record.AddItems, record.RemainItems,
		record.RetCode, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveCollectChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveCollectChgRecord succ", zap.Any("flowrecord", record))

	return
}
