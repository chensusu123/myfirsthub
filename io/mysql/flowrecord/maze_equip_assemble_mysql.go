package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/mysql"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeEquipAssembleRecordTableName = "maze_equip_assemble_record"

func SaveEquipAssembleRecord(logger fklog.FKLogI, record *dollequipassmeblekakfa.MazeGameEquipAssembleRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipAssembleRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	record.SectionID = appconfig.GlobalConfig().Global.SectionID

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveEquipAssembleRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveEquipAssembleRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveEquipAssembleRecord succ", zap.Any("flowrecord", record))

	return
}
