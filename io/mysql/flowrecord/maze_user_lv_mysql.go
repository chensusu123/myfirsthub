package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/mysql"
)

const MazeUserLevelRecordTableName = "maze_user_level_record"

// 保存用户等级变化流水
func SaveUserLevelRecord(logger fklog.FKLogI, record *mazeuserlevelkafka.MazeUserLevelRecord) {
	db, err := mysql.GetMysqlDb(MazeUserLevelRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeUserLevelRecordTableName:", MazeUserLevelRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeUserLevelRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `old_level`, `old_total_exp`, `new_level`, `new_total_exp`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.OldLevel, record.OldTotalExp, record.NewLevel, record.NewTotalExp, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveUserLevelRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveUserLevelRecord succ", zap.Any("flowrecord", record))

	return
}
