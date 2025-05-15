// @Author: ZhaoXiming 2025/3/24 19:31
// @Desc:

package tradeno

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
)

func GetTradeNum() (tradeNum uint64) {
	return uniqueid.NewTradeNoMaker(uint64(fkconfig.GetServerConfig().ServerID)).MakeTradeNo()
}
