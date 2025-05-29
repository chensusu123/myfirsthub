package flowrecord

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/mysql"
)

const MazeAttrChgRecordTableName = "maze_attr_chg_record"

// 保存迷宫属性变化流水
func SaveAttrChgRecord(logger fklog.FKLogI, record *structsdef.MazeGameAttrChgRecord) {
	db, err := mysql.GetMysqlDb(MazeAttrChgRecordTableName)
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err), zap.Any("MazeAttrChgRecordTableName:", MazeAttrChgRecordTableName))
		return
	}
	tableName := mysql.GetShardingTableName(MazeAttrChgRecordTableName)
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `attr_id`, `attr_type`, `new_val`, `old_val`, `chg_type`, `chg_sub_type`, `chg_desc`, `extra`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		tableName,
	)
	_, err = db.Exec(sqlStr,
		record.UserId, record.AttrId, record.AttrType, record.NewVal, record.OldVal, record.ChgType, record.ChgSubType,
		record.ChgDesc, record.Extra, record.GroupId, record.CreateTime, record.ServerId)
	if err != nil {
		logger.ErrorWF("SaveAttrChgRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.InfoWF("SaveAttrChgRecord succ", zap.Any("flowrecord", record))

	return
}
