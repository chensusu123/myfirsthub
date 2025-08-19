package flowrecord

import (
	"maze_game_server/io/kafka/mazeenergyrecord"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

const MazeUserEnergyRecordTableName = "maze_user_energy_record"

// 保存用户体力变化流水
func SaveUserEnergyRecord(logger fklog.FKLogI, record *mazeenergyrecord.MazeEnergyChgRecord) {
	// nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserEnergyRecordTableName), ".")

	// record.DataBase = nowDbTable[0]
	// record.Table = nowDbTable[1]
	// record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// // 打到kafka 中
	// data, err := json.Marshal(record)
	// if err != nil {
	// 	logger.ErrorWF("SaveUserEnergyRecord Marshal Fail",
	// 		zap.Any("record", record))
	// 	return
	// }
	// err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	// if err != nil {
	// 	logger.ErrorWF("SaveUserEnergyRecord SendMsg",
	// 		zap.Any("record", record),
	// 		zap.Error(err),
	// 	)
	// 	return
	// }
	// logger.InfoWF("SaveUserEnergyRecord succ", zap.Any("flowrecord", record))

	return
}
