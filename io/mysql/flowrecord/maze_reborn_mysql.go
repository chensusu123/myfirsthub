package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeUserRebornRecordTableName = "maze_user_reborn_record"

// 保存用户复活流水
func SaveRebornRecord(logger fklog.FKLogI, record *mazerebornkafka.MazeRebornRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserRebornRecordTableName), ".")

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
		logger.ErrorWF("SaveRebornRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveRebornRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveRebornRecord succ", zap.Any("flowrecord", record))

	return
}
