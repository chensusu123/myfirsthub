// @Author: ZhaoXiming 2025/3/24 19:31
// @Desc:

package tradeno

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
	"strconv"
)

func GetTradeNum() (tradeNum uint64) {
	sectionId := appconfig.GlobalConfig().Global.SectionID
	serverId, _ := strconv.ParseUint(sectionId, 10, 64)
	return uniqueid.NewTradeNoMaker(serverId).MakeTradeNo()
}
