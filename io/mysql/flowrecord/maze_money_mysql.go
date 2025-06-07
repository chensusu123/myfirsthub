package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/mysql"
)

const MazeMoneyRecordTableName = "maze_money_record"

// 保存用户货币变化流水
func SaveMoneyRecord(logger fklog.FKLogI, record *mazemoneykafka.MazeMoneyRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeMoneyRecordTableName:", MazeMoneyRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeMoneyRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveMoneyRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveMoneyRecord succ", zap.Any("flowrecord", record))

	return
}
