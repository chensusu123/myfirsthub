package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/mysql"
)

const MazeUserRebornRecordTableName = "maze_user_reborn_record"

// 保存用户复活流水
func SaveRebornRecord(logger fklog.FKLogI, record *mazerebornkafka.MazeRebornRecord) {
	db, err := mysql.GetMysqlDb(MazeUserRebornRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeUserRebornRecordTableName:", MazeUserRebornRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeUserRebornRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `barrier`, `reborn_count`, `reborn_cost`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.Barrier, record.RebornCount, record.RebornCost, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveRebornRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveRebornRecord succ", zap.Any("flowrecord", record))

	return
}
