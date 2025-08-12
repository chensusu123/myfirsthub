package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/io/mysql"
)

const MazeUserEnergyRecordTableName = "maze_user_energy_record"

// 保存用户体力变化流水
func SaveUserEnergyRecord(logger fklog.FKLogI, record *mazeenergyrecord.MazeEnergyChgRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeUserEnergyRecordTableName:", MazeUserEnergyRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeUserEnergyRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveUserEnergyRecord fail", zap.Error(err), zap.Any("record", record))
		return
	}
	logger.InfoWF("SaveUserEnergyRecord success", zap.Any("record", record))

	return
}
