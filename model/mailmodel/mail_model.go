package mailmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

// 附件结构
type Attachment struct {
	ItemID int32  `json:"item_id,omitempty"` // 物品ID
	Count  int64  `json:"count,omitempty"`   // 数量
	Extra  string `json:"extra,omitempty"`   // 附加信息
}

type MailInfo struct {
	ID          uint64        `json:"mail_id,omitempty"`       // 邮件ID
	Title       string        `json:"title,omitempty"`         // 标题
	Content     string        `json:"content,omitempty"`       // 描述
	Label       int32         `json:"label,omitempty"`         // 标签
	Attachments []*Attachment `json:"attachments,omitempty"`   // 附件
	Sender      string        `json:"sender,omitempty"`        // 发送者
	ReciverID   uint64        `json:"reciver_id,omitempty"`    // 接收者
	IsRead      bool          `json:"is_read,omitempty"`       // 是否已读
	IsGetAttach bool          `json:"is_get_attach,omitempty"` // 是否领取附件
	SendTime    int64         `json:"send_time,omitempty"`     // 发送时间戳(毫秒)
	ExpireTime  int64         `json:"expire_time,omitempty"`   // 过期时间戳(毫秒，0表示永不过期)
}

type MailModel struct {
	MailMap map[int32]map[uint64]*MailInfo `json:"mail_map,omitempty"`
}

func getMailKey(userId uint64) string {
	return fmt.Sprintf("maze:mail:u:%d", userId)
}

func NewMailModel(ctx context.Context, userID uint64) (*MailModel, error) {
	mailModel := &MailModel{
		MailMap: make(map[int32]map[uint64]*MailInfo),
	}
	if err := mailModel.load(ctx, userID); err != nil {
		return nil, err
	}
	return mailModel, nil
}

func (p *MailModel) load(ctx context.Context, userID uint64) (err error) {
	return io.LoadSvrData(context.TODO(), getMailKey(userID), p)
}

func (p *MailModel) Save(ctx context.Context, userID uint64) (err error) {
	return io.SaveSvrData(context.TODO(), getMailKey(userID), p)
}

func (p *MailModel) Del(ctx context.Context, userID uint64, stageId int32) (err error) {
	return io.DeleteSvrData(context.TODO(), getMailKey(userID))
}
