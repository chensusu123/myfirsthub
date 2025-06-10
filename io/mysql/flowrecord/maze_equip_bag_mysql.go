package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/mysql"
)

const MazeEquipBagRecordTableName = "maze_equip_bag_record"

// 保存背包流水
func SaveEquipBagRecord(logger fklog.FKLogI, record *mazeequipbagrecord.MazeGameEquipBagRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipBagRecordTableName:", MazeEquipBagRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeEquipBagRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveEquipBagRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipBagRecord succ", zap.Any("flowrecord", record))

	return
}
