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

const MazeBarrierUserRecordTableName = "maze_barrier_user_record"

// 保存用户迷宫闯关流水
func SaveBarrierUserRecord(logger fklog.FKLogI, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeBarrierUserRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveBarrierUserRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveBarrierUserRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveBarrierUserRecord succ", zap.Any("flowrecord", record))

	return
}
