package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

const MazeMoneyRecordTableName = "maze_money_record"

// 保存用户货币变化流水
func SaveMoneyRecord(logger fklog.FKLogI, record *mazemoneykafka.MazeMoneyRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeMoneyRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveMoneyRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveMoneyRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveMoneyRecord succ", zap.Any("flowrecord", record))

	return
}
