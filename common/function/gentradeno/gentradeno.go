package gentradeno

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
	"gitlab.ifreetalk.com/maze-plate/io/kafka_interface/item/error_record_kafka"
	"gitlab.ifreetalk.com/maze-plate/protodef/Common"
)

// 扣除物品
func GenSubGoodFailRecord(agent fkserver.UserContext, tradeNo uint64, opType, addType int32, desc string) *error_record_kafka.SendGoodFailRecord {
	svrType := fkconfig.GetServerConfig().ServerTypeID
	return error_record_kafka.NewGoodFailRecord(agent.UserID,
		tradeNo,
		desc, svrType, opType, addType,
		error_record_kafka.TradeTypeForSub)
}

func ItemToFailItem(in *Common.Item) *error_record_kafka.FailInfo {
	recordItem := &error_record_kafka.FailInfo{}
	recordItem.ItemId = in.GetItemId()
	recordItem.Count = in.GetCount()
	//recordItem.LimitType = 0
	return recordItem
}

func ItemToTimeoutItem(in *Common.Item) *error_record_kafka.TimeOutInfo {
	recordItem := &error_record_kafka.TimeOutInfo{}
	recordItem.ItemId = in.GetItemId()
	recordItem.Count = in.GetCount()
	//recordItem.LimitType = 0
	return recordItem
}

func ItemsToFailRecord(in []*Common.Item) (out []*error_record_kafka.FailInfo) {
	for _, item := range in {
		out = append(out, ItemToFailItem(item))
	}
	return
}

func ItemsToTimeoutRecord(in []*Common.Item) (out []*error_record_kafka.TimeOutInfo) {
	for _, item := range in {
		out = append(out, ItemToTimeoutItem(item))
	}
	return
}

func GetTradeNum() (tradeNum uint64) {
	m := uniqueid.NewTradeNoMaker(uint64(fkconfig.GetServerConfig().ServerID))
	if m == nil {
		return
	}
	tradeNum = m.MakeTradeNo()
	return
}
