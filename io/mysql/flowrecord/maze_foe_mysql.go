package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/dollmazefoekafka"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeFoeRecordTableName = "maze_foe_record"

// 保存用户打怪流水
func SaveFoeRecord(logger fklog.FKLogI, record *dollmazefoekafka.DollMazeFoeRecord) {
	if record.AwardList == "null" {
		record.AwardList = ""
	}
	if record.CreateTime == 0 {
		record.CreateTime = time.Now().UnixNano() / 1000000
	}

	record.FoeList = ""

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeFoeRecordTableName), ".")
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
		logger.ErrorWF("SaveFoeRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveFoeRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveFoeRecord succ", zap.Any("flowrecord", record))

	return
}
