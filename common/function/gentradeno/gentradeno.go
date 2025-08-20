package gentradeno

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
)

func GetTradeNum() (tradeNum uint64) {
	m := uniqueid.NewTradeNoMaker(uint64(fkconfig.GetServerConfig().ServerID))
	if m == nil {
		return
	}
	tradeNum = m.MakeTradeNo()
	return
}
