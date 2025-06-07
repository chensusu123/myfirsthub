package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/mysql"
)

const MazeUserLevelRecordTableName = "maze_user_level_record"

// 保存用户等级变化流水
func SaveUserLevelRecord(logger fklog.FKLogI, record *mazeuserlevelkafka.MazeUserLevelRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeUserLevelRecordTableName:", MazeUserLevelRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeUserLevelRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveUserLevelRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveUserLevelRecord succ", zap.Any("flowrecord", record))

	return
}
