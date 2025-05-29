package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollmazefoekafka"
	"maze_game_server/io/mysql"
)

const MazeFoeRecordTableName = "maze_foe_record"

// 保存用户打怪流水
func SaveFoeRecord(logger fklog.FKLogI, record *dollmazefoekafka.DollMazeFoeRecord) {
	db, err := mysql.GetMysqlDb(MazeFoeRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeFoeRecordTableName:", MazeFoeRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeFoeRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `barrier`, `area`, `level`, `master_id`, `award_list`, `equips`, `equip_points`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.Barrier, record.Area, record.Level, record.MasterId, record.AwardList, record.Equips, record.EquipPoints,
		record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveFoeRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveFoeRecord succ", zap.Any("flowrecord", record))

	return
}
