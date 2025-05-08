/*
@Author: xiaobo
@Date: 2023/10/27 11:19
@Description:
*/

package cgkargs

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
)

// AlreadyReadMailExpireTime 已读、已领取信封过期时间
var AlreadyReadMailExpireTime int64 = 86400

var (
	BroadUserCount     int64 = 5 // 发送广播信封时，同时处理的用户数量
	GetLeagueUserCount int64 = 5 // 获取联盟下的家族的用户，同时获取的家族的数量
)

func init() {
	param.Int64P(&BroadUserCount, "broad:user:count", 5)
	param.Int64P(&GetLeagueUserCount, "get:league:user:count", 5)
	fkconfig.RegCGKParamChange(OnParamChg)
}

func OnParamChg(logger fklog.FKLogI) (err error) {

	logger.InfoWF("OnParamChg end")
	return
}
