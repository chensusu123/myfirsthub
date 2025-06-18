// @Author pangchenyang 2025/6/18 11:05:00
// @Desc: 
package flowrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/mysql"
)

/*
CREATE TABLE maze_alter_profile_record (
    `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
    `new_val` bigint(0) NOT NULL DEFAULT 0 COMMENT '新值',
    `old_val` bigint(0) NOT NULL DEFAULT 0 COMMENT '旧值',
    `create_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '创建时间',
	`chg_desc` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原因描述',
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `user_id`(`user_id`) USING BTREE,
  	INDEX `create_time`(`create_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='修改用户资料流水表';
*/

const AlterProfileTableName = "alter_profile_record"

type AlterProfileRecord struct {
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`         // 用户Id
	NewVal     string `json:"new_val" gorm:"column:new_val"`         // 新值
	OldVal     string `json:"old_val" gorm:"column:old_val"`         // 旧值
	ChgDesc    string `json:"chg_desc" gorm:"column:chg_desc"`       // 变化原因
	CreateTime int64  `json:"create_time" gorm:"column:create_time"` // 时间戳 ms
}

// SaveAlterProfileRecord 保存修改用户资料流水
func SaveAlterProfileRecord(logger fklog.FKLogI, record AlterProfileRecord) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("SaveAlterProfileRecord GetMysqlDb fail",
			zap.String("table", AlterProfileTableName), zap.Error(err))
		return
	}

	res := db.Table(mysql.GetFullyQualifiedTableName(AlterProfileTableName)).Create(record)
	if res.Error != nil {
		logger.ErrorWF("SaveAlterProfileRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	logger.DebugWF("SaveAlterProfileRecord succ", zap.Any("flowrecord", record))
	return
}
