// @Desc:
// @Author: LiuChongyu 2021/11/17 5:13 下午
// @Update: LiuChongyu 2021/11/17 5:13 下午

package UnionIDBindMysql

import (
	"database/sql"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkmysql"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/zap"
)

var gMysqlCli = &fkmysql.MysqlDB{}

func init() {
	fkconfig.RegisterNameNode("UnionIDBindMysql", 2220, gMysqlCli)
}

func getDBIndex(id uint64) int32 {
	return int32(id % 16)
}

func getTBIndex(id uint64) int32 {
	return int32((id >> 4) % 16)
}

// GetUserID2UnionID 获取user_id对应的union_id列表
func GetUserID2UnionID(logger fklog.FKLogI, userID uint64) (uint64, error) {
	dbIndex := getDBIndex(userID)
	tbIndex := getTBIndex(userID)
	sqlStr := fkutil.K_str("select union_id from t_user_id_to_union_id_%d where user_id=?", tbIndex)

	db, err := gMysqlCli.GetDB(dbIndex)
	if err != nil {
		gMysqlCli.ErrorWF("get DB error:",
			zap.Any("err", err),
			zap.Uint64("userID:", userID),
		)
		return 0, err
	}

	var unionID uint64
	err = db.QueryRow(sqlStr, userID).Scan(&unionID)
	if err == sql.ErrNoRows {
		logger.WarnWF("user_id not union_id",
			zap.Uint64("userID", userID),
		)
		return 0, nil
	} else if err != nil {
		logger.ErrorWF("failed to query sql, error:",
			zap.Any("err", err),
			zap.Uint64("userID:", userID),
		)
		return 0, err
	}

	return unionID, nil
}
