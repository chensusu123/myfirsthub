package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/mysql"
)

const MazeEquipAssembleRecordTableName = "maze_equip_assemble_record"

func SaveEquipAssembleRecord(logger fklog.FKLogI, record *dollequipassmeblekakfa.MazeGameEquipAssembleRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeEquipAssembleRecordTableName:", MazeEquipAssembleRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeEquipAssembleRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveEquipAssembleRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveEquipAssembleRecord succ", zap.Any("flowrecord", record))

	return
}
