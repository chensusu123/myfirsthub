/*
* @Author: majian
* @Date: 2022-08-24 11:39
 */
package uniqueid

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
)

// 生成有符号64位唯一ID
func GenUniqueIdInt64() int64 {
	uniqueid.NewTextMsgIdMaker(int64(fkconfig.GetServerConfig().ServerID))
	return uniqueid.GTextMsgIdMaker.MakeTextMsgId()
}

// 生成无符号64位唯一ID
func GenUniqueIdUInt64() uint64 {
	uniqueid.NewTradeNoMaker(uint64(fkconfig.GetServerConfig().ServerID))
	return uniqueid.GTradeNoMaker.MakeTradeNo()
}
