package flowrecord

import (
	"encoding/json"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/mysql"
)

const MazeTempBuffChangeRecordTableName = "maze_temp_buff_change_record"

// 保存临时buff变化流水
func SaveTempBuffChgRecord(logger fklog.FKLogI, record *mazetempbuffchgmsg.MazeTempBuffChangeMsg) {
	db, err := mysql.GetMysqlDb(MazeTempBuffChangeRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeTempBuffChangeRecordTableName:", MazeTempBuffChangeRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeTempBuffChangeRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `stage_id`, `chg_attrs`, `chg_type`, `chg_desc`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?)",
		tableName,
	)
	chgAttrs, err := json.Marshal(record.ChgAttrs)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord marshal failed", zap.Any("record", record), zap.Error(err))
		return
	}
	_, err = db.Exec(sqlStr,
		record.UserId, record.StageId, chgAttrs, record.ChgType, record.ChgDesc, record.GroupId, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveTempBuffChgRecord succ", zap.Any("flowrecord", record))

	return
}
