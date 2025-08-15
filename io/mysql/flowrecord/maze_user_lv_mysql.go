package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeUserLevelRecordTableName = "maze_user_level_record"

// 保存用户等级变化流水
func SaveUserLevelRecord(logger fklog.FKLogI, record *mazeuserlevelkafka.MazeUserLevelRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserLevelRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveUserLevelRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveUserLevelRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveUserLevelRecord succ", zap.Any("flowrecord", record))

	return
}
