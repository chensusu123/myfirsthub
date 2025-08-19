package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeCollectChgRecordTableName = "maze_collect_chg_record"

// 迷宫挂机变化流水
func SaveCollectChgRecord(logger fklog.FKLogI, record *mazecollectrecord.MazeCollectChgRecord) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeCollectChgRecordTableName), ".")

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
		logger.ErrorWF("SaveCollectChgRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveCollectChgRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveCollectChgRecord succ", zap.Any("flowrecord", record))

	return
}
