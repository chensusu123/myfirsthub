package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollmazefoekafka"
	"maze_game_server/io/mysql"
)

const MazeFoeRecordTableName = "maze_foe_record"

// 保存用户打怪流水
func SaveFoeRecord(logger fklog.FKLogI, record *dollmazefoekafka.DollMazeFoeRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeFoeRecordTableName:", MazeFoeRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeFoeRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveFoeRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveFoeRecord succ", zap.Any("flowrecord", record))

	return
}
