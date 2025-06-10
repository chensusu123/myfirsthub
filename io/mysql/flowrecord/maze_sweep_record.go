package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/mysql"
)

const MazeSweepRecordTableName = "maze_sweep_record"

func SaveSweepRecord(logger fklog.FKLogI, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) {
	if record.GameRet != mazebarrieruserkafka.GameRetSweep {
		// 只要扫荡的记录，其他不处理
		return
	}
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeSweepRecordTableName:", MazeSweepRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeSweepRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveSweepRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveSweepRecord succ", zap.Any("flowrecord", record))

	return
}
