package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/mysql"
)

const MazeUserRebornRecordTableName = "maze_user_reborn_record"

// 保存用户复活流水
func SaveRebornRecord(logger fklog.FKLogI, record *mazerebornkafka.MazeRebornRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeUserRebornRecordTableName:", MazeUserRebornRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeUserRebornRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveRebornRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveRebornRecord succ", zap.Any("flowrecord", record))

	return
}
