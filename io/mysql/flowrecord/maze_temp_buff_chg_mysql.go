package flowrecord

import (
	"encoding/json"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/mysql"
)

const MazeTempBuffChangeRecordTableName = "maze_temp_buff_change_record"

// 保存临时buff变化流水
func SaveTempBuffChgRecord(logger fklog.FKLogI, record *mazetempbuffchgmsg.MazeTempBuffChangeMsg) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeTempBuffChangeRecordTableName:", MazeTempBuffChangeRecordTableName))
		return
	}

	chgAttrs, err := json.Marshal(record.ChgAttrs)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord marshal failed", zap.Any("record", record), zap.Error(err))
		return
	}
	record.ChgAttrsStr = string(chgAttrs)
	res := db.Table(mysql.GetFullyQualifiedTableName(MazeTempBuffChangeRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveTempBuffChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveTempBuffChgRecord succ", zap.Any("flowrecord", record))

	return
}
