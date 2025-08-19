package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

const MazeUserRebornRecordTableName = "maze_user_reborn_record"

// 保存用户复活流水
func SaveRebornRecord(logger fklog.FKLogI, record *mazerebornkafka.MazeRebornRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserRebornRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveRebornRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveRebornRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveRebornRecord succ", zap.Any("flowrecord", record))

	return
}
