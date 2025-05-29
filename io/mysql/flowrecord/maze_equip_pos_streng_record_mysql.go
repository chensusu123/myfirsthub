package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/mysql"
)

const MazeEquipPosLevelUpRecordTableName = "maze_equip_pos_level_up_record"

// 保存装备位强化流水
func SaveEquipPosStrengRecord(logger fklog.FKLogI, record *equipposstrengrecordkafka.EquipPosLevelUpRecord) {
	db, err := mysql.GetMysqlDb(MazeEquipPosLevelUpRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipPosLevelUpRecordTableName:", MazeEquipPosLevelUpRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeEquipPosLevelUpRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `pos_id`, `old_pos_lv`, `new_pos_lv`, `old_pos_suit_id`, `new_pos_suit_id`, `trade_no`, `cost_items`, `result`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.PosId, record.OldPosLv, record.NewPosLv, record.OldPosSuitId, record.NewPosSuitId, record.TradeNo,
		record.CostItems, record.Result, record.GroupId, record.OpTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveEquipPosStrengRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipPosStrengRecord succ", zap.Any("flowrecord", record))

	return
}
