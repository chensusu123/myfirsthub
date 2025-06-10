package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/mysql"
)

const MazeEquipPosLevelUpRecordTableName = "maze_equip_pos_level_up_record"

// 保存装备位强化流水
func SaveEquipPosStrengRecord(logger fklog.FKLogI, record *equipposstrengrecordkafka.EquipPosLevelUpRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipPosLevelUpRecordTableName:", MazeEquipPosLevelUpRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeEquipPosLevelUpRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveEquipPosStrengRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipPosStrengRecord succ", zap.Any("flowrecord", record))

	return
}
