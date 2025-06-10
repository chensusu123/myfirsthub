package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/mysql"
)

const MazeEquipDismantRecordTableName = "maze_equip_dismant_record"

// 保存装备分解流水
func SaveEquipDismantRecord(logger fklog.FKLogI, record *dollequipdismantlekafka.MazeGameEquipDismantleRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipDismantRecordTableName:", MazeEquipDismantRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeEquipDismantRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveEquipDismantRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipDismantRecord succ", zap.Any("flowrecord", record))

	return
}
