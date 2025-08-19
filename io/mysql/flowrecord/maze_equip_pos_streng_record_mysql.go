package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeEquipPosLevelUpRecordTableName = "maze_equip_pos_level_up_record"

// 保存装备位强化流水
func SaveEquipPosStrengRecord(logger fklog.FKLogI, record *equipposstrengrecordkafka.EquipPosLevelUpRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipPosLevelUpRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]
	groupID, err := strconv.Atoi(appconfig.GlobalConfig().Global.SectionID)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord Atoi fail", zap.Error(err))
		return
	}
	record.GroupId = uint32(groupID)

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveEquipPosStrengRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveEquipPosStrengRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveEquipPosStrengRecord succ", zap.Any("flowrecord", record))

	return
}
