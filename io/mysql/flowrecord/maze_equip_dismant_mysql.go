package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeEquipDismantRecordTableName = "maze_equip_dismant_record"

// 保存装备分解流水
func SaveEquipDismantRecord(logger fklog.FKLogI, record *dollequipdismantlekafka.MazeGameEquipDismantleRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipDismantRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	groupID, err := strconv.Atoi(appconfig.GlobalConfig().Global.SectionID)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord Atoi fail", zap.Error(err))
		return
	}
	record.GroupID = uint32(groupID)

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveEquipDismantRecord succ", zap.Any("flowrecord", record))

	return
}
