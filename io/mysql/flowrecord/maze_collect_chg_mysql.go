package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/mysql"
)

const MazeCollectChgRecordTableName = "maze_collect_chg_record"

// 迷宫挂机变化流水
func SaveCollectChgRecord(logger fklog.FKLogI, record *mazecollectrecord.MazeCollectChgRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeCollectChgRecordTableName:", MazeCollectChgRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeCollectChgRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveCollectChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveCollectChgRecord succ", zap.Any("flowrecord", record))

	return
}
