package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/mysql"
)

const MazeAttrChgRecordTableName = "maze_attr_chg_record"

// 保存迷宫属性变化流水
func SaveAttrChgRecord(logger fklog.FKLogI, record *structsdef.MazeGameAttrChgRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeAttrChgRecordTableName:", MazeAttrChgRecordTableName))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(MazeAttrChgRecordTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveAttrChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveAttrChgRecord succ", zap.Any("flowrecord", record))

	return
}
