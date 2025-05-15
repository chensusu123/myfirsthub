package error_record_kafka

import "time"

type SendGoodFailRecord struct {
	UserId        uint64         `json:"user_id"`
	ServiceID     uint32         `json:"service_id"`
	TradeNo       uint64         `json:"trade_no"`
	OpType        int32          `json:"op_type"`
	AddType       int32          `json:"add_type"`
	TradeType     int32          `json:"trade_type"` // 1 增加 2 扣除
	CreateDt      int64          `json:"create_dt"`
	ExceptionDesc string         `json:"exception_desc"`
	FailJson      []*FailInfo    `json:"-"`            // 内部使用
	TimeOutJson   []*TimeOutInfo `json:"-"`            // 内部使用
	StrFailJson   string         `json:"fail_json"`    // 添加道具时加失败的，扣道具时扣成功的
	StrTimeOut    string         `json:"timeout_json"` // 超时的
}

var (
	TradeTypeForAdd int32 = 1
	TradeTypeForSub int32 = 2
)

type FailInfo struct {
	ItemId    int32 `json:"item_id"`
	Count     int64 `json:"count"`
	LimitType int32 `json:"limit_type"`
}

type TimeOutInfo struct {
	ItemId    int32 `json:"item_id"`
	Count     int64 `json:"count"`
	LimitType int32 `json:"limit_type"`
}

func NewGoodFailRecord(userId, tradeNo uint64, desc string, serverId uint32, opType, addType, tradeType int32) *SendGoodFailRecord {
	record := new(SendGoodFailRecord)
	record.UserId = userId
	record.TradeNo = tradeNo
	record.ServiceID = serverId
	record.AddType = addType
	record.OpType = opType
	record.TradeType = tradeType
	record.ExceptionDesc = desc
	record.CreateDt = int64(time.Now().UnixMilli())
	return record
}
