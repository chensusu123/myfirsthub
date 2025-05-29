package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/mysql"
)

const MazeBarrierUserRecordTableName = "maze_barrier_user_record"

// 保存用户迷宫闯关流水
func SaveBarrierUserRecord(logger fklog.FKLogI, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) {
	db, err := mysql.GetMysqlDb(MazeBarrierUserRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeBarrierUserRecordTableName:", MazeBarrierUserRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeBarrierUserRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `barrier`, `game_ret`, `awards`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.Barrier, record.GameRet, record.Awards, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveBarrierUserRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveBarrierUserRecord succ", zap.Any("flowrecord", record))

	return
}
