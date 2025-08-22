package flowrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

const MazeTempBuffChangeRecordTableName = "maze_temp_buff_change_record"

// 保存临时buff变化流水
func SaveTempBuffChgRecord(logger fklog.FKLogI, record *mazetempbuffchgmsg.MazeTempBuffChangeMsg) {
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeTempBuffChangeRecordTableName), ".")

	record.DataBase = nowDbTable[0]
	record.Table = nowDbTable[1]

	groupID, err := strconv.Atoi(appconfig.GlobalConfig().Global.SectionID)
	if err != nil {
		logger.ErrorWF("SaveEquipDismantRecord Atoi fail", zap.Error(err))
		return
	}
	record.GroupId = uint32(groupID)

	record.ChgAttrsStr = attr2String(record.ChgAttrs)

	// 防止该字段多余序列化 导致流水打点误插入数据库
	//record.ChgAttrs = nil

	// 打到kafka 中
	data, err := json.Marshal(record)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord Marshal Fail",
			zap.Any("record", record))
		return
	}
	err = kafka.GflowKafka.SendMsg(context.TODO(), fmt.Sprintf("%v", time.Now().UnixNano()), data)
	if err != nil {
		logger.ErrorWF("SaveTempBuffChgRecord SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		return
	}
	logger.InfoWF("SaveTempBuffChgRecord succ", zap.Any("flowrecord", record))

	return
}

func attr2String(attrs []*mazetempbuffchgmsg.AttrChgInfo) string {
	res := ""
	for _, val := range attrs {
		if res == "" {
			res = fmt.Sprintf("%d:%d:%d", val.AttrId, val.OldVal, val.CurVal)
		}
		res = fmt.Sprintf("%s_%d:%d:%d", res, val.AttrId, val.OldVal, val.CurVal)
	}
	return res
}
