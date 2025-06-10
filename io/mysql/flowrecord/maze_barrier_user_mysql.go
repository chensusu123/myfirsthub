package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/mysql"
)

const MazeBarrierUserRecordTableName = "maze_barrier_user_record"

// 保存用户迷宫闯关流水
func SaveBarrierUserRecord(logger fklog.FKLogI, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeBarrierUserRecordTableName:", MazeBarrierUserRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeBarrierUserRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveBarrierUserRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveBarrierUserRecord succ", zap.Any("flowrecord", record))

	return
}
