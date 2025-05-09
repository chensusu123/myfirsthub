// @Desc:
// @Author: LiuChongyu 2021/12/15 5:43 下午
// @Update: LiuChongyu 2021/12/15 5:43 下午

package PaipaiUnionIDAuthMysql

import (
	"database/sql"

	"gitlab.ifreetalk.com/plate/freetk/fkutil"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkmysql"
)

func init() {
	fkconfig.RegisterNameNode("PaipaiUnionIDAuthMysql", 2223, a)
}

var a = &fkmysql.MysqlDB{}

func getDBIndex(id uint64) int32 {
	return 0
}

const tableCount = 256

func getTBIndex(id uint64) int32 {
	return int32(id % tableCount)
}

func GetAuth(logger fklog.FKLogI, unionID uint64) (authToken int64, accessToken string, sequence, expireTime int64, err error) {
	dbIndex := getDBIndex(uint64(unionID))
	tbIndex := getTBIndex(uint64(unionID))

	sqlStr := fkutil.K_str("select  auth_token,access_token,sequence, expire_dt from t_paipai_unionid_auth_%d "+
		"where unionid=?  limit 1;", tbIndex)

	db, err := a.GetDB(dbIndex)
	if err != nil {
		a.ErrorWF("PaipaiUnionIDAuthMysql.GetAuth get GetDB failed",
			fklog.Uint64("unionID", unionID),
			fklog.Any("err", err),
		)
		return 0, "", 0, 0, err
	}
	row := db.QueryRow(sqlStr, unionID)
	if row == nil {
		a.InfoWF("PaipaiUnionIDAuthMysql.GetAuth not found",
			fklog.Uint64("unionID", unionID),
			fklog.String("sqlStr", sqlStr),
		)
		return 0, "", 0, 0, nil
	}

	var tmpAuthToken sql.NullInt64
	var tmpAccessToken sql.NullString
	var tmpSequence sql.NullInt64
	var tmpExpireTime sql.NullInt64

	err = row.Scan(&tmpAuthToken, &tmpAccessToken, &tmpSequence, &tmpExpireTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没数据，不算失败
			return 0, "", 0, 0, nil
		}
		a.ErrorWF("PaipaiUnionIDAuthMysql.GetAuth call Scan failed",
			fklog.Uint64("unionID", unionID),
			fklog.String("sqlStr", sqlStr),
			fklog.Any("err", err),
		)
		return 0, "", 0, 0, err
	}

	authToken = tmpAuthToken.Int64
	accessToken = tmpAccessToken.String
	sequence = tmpSequence.Int64
	expireTime = tmpExpireTime.Int64

	a.InfoWF("PaipaiUnionIDAuthMysql.GetAuth dump",
		fklog.Uint64("unionID", unionID),
		fklog.String("accessToken", accessToken),
		fklog.Int64("authToken", authToken),
		fklog.Int64("expireTime", expireTime),
		fklog.Int64("sequence", sequence),
		fklog.String("sqlStr", sqlStr),
	)
	return authToken, accessToken, sequence, expireTime, nil
}
