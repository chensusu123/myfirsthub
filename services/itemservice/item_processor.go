package itemservice

import (
	"context"
	"maze_game_server/pb/common/MessageType"
)

var GloRegIns = NewRegister()

type AddItemRes struct {
	SucItem  []*ItemInfo // 操作成功的道具
	FailItem []*ItemInfo // 操作失败的道具(业务逻辑错误)
}

func (a *AddItemRes) Merge(b *AddItemRes) {
	if a == nil {
		return
	}
	if b == nil {
		return
	}
	a.SucItem = append(a.SucItem, b.SucItem...)
	a.FailItem = append(a.FailItem, b.FailItem...)
}

type ItemProcessor interface {
	// GatherItem 加道具
	GatherItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items ...*ItemInfo) (*AddItemRes, error)
	// DeductItem 扣道具
	DeductItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*ItemInfo) (*AddItemRes, *MessageType.ErrorInfo)
	// DeductItemCheck 扣道具检查
	DeductItemCheck(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*ItemInfo) (*AddItemRes, *MessageType.ErrorInfo)
	// GetItem 查询数据
	GetItem(ctx context.Context, userId uint64, items []*ItemInfo) *MessageType.ErrorInfo
}
