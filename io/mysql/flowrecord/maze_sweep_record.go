package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeSweepRecordTableName = "maze_sweep_record"

func SaveSweepRecord(logger fklog.FKLogI, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) {
	if record.GameRet != mazebarrieruserkafka.GameRetSweep {
		// 只要扫荡的记录，其他不处理
		return
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeSweepRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveSweepRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveSweepRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveSweepRecord succ", zap.Any("flowrecord", record))

	return
}
