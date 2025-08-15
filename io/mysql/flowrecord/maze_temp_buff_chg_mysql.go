package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeTempBuffChangeRecordTableName = "maze_temp_buff_change_record"

// 保存临时buff变化流水
func SaveTempBuffChgRecord(logger fklog.FKLogI, record *mazetempbuffchgmsg.MazeTempBuffChangeMsg) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeTempBuffChangeRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveTempBuffChgRecord succ", zap.Any("flowrecord", record))

	return
}
