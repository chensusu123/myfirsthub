package mailmodel

import (
	"fmt"
	"maze_game_server/io"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
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

	// 新增字段 - 支持完整的适配押镖邮件服务功能
	MailType          int32  `json:"mail_type,omitempty"`            // 邮件类型
	MailSubType       int32  `json:"mail_sub_type,omitempty"`        // 邮件子类型
	MailContentType   int32  `json:"mail_content_type,omitempty"`    // 邮件内容类型
	ExtType           int32  `json:"ext_type,omitempty"`             // 发奖类型信封使用该字段区分业务流水
	SourceType        int32  `json:"source_type,omitempty"`          // 标示信封或者电话
	PhoneType         int32  `json:"phone_type,omitempty"`           // 电话类型
	BroadcastFlag     int32  `json:"broadcast_flag,omitempty"`       // 广播标志
	BroadcastId       uint64 `json:"broadcast_id,omitempty"`         // 和broadcast_flag配合使用
	ClassType         int32  `json:"class_type,omitempty"`           // 信封类别
	NotTotalRecv      bool   `json:"not_total_recv,omitempty"`       // 是否可以一键领取
	PagePopUp         bool   `json:"page_pop_up,omitempty"`          // 客户端是否需要弹框展示
	AllSupportPop     bool   `json:"all_support_pop,omitempty"`      // 新的支持所有信封样式的控制自动弹出
	UserLimitType     int32  `json:"user_limit_type,omitempty"`      // 限制发送用户类型
	UserLimitValueMin int32  `json:"user_limit_value_min,omitempty"` // 限制发送用户最小数值
	UserLimitValueMax int32  `json:"user_limit_value_max,omitempty"` // 限制发送用户最大数值

	// 富文本相关
	RichContent string `json:"rich_content,omitempty"` // 富文本内容
	PanelInfo   string `json:"panel_info,omitempty"`   // 面板信息(JSON)

	// 高级功能
	VoteInfo         string `json:"vote_info,omitempty"`          // 投票信息(JSON)
	LinkUrl          string `json:"link_url,omitempty"`           // 链接URL
	LinkParam        string `json:"link_param,omitempty"`         // 跳转参数(JSON)
	OfflinePushText  string `json:"offline_push_text,omitempty"`  // 离线推送消息内容
	ThanksNote       string `json:"thanks_note,omitempty"`        // 感谢语
	BattleRecordBody string `json:"battle_record_body,omitempty"` // 战报内容
	RewardTimes      int32  `json:"reward_times,omitempty"`       // 多次奖励合并发放设置
}

type MailModel struct {
	MailMap map[int32]map[uint64]*MailInfo `json:"mail_map,omitempty"`
}

func getMailKey(userId uint64) string {
	return fmt.Sprintf("maze:mail:u:%d", userId)
}

func NewMailModel(logger fklog.FKLogI, userID uint64) (*MailModel, error) {
	mailModel := &MailModel{
		MailMap: make(map[int32]map[uint64]*MailInfo),
	}
	if err := mailModel.load(logger, userID); err != nil {
		return nil, err
	}
	return mailModel, nil
}

func (p *MailModel) load(logger fklog.FKLogI, userID uint64) (err error) {
	return io.LoadSvrData(logger, getMailKey(userID), p)
}

func (p *MailModel) Save(logger fklog.FKLogI, userID uint64) (err error) {
	return io.SaveSvrData(logger, getMailKey(userID), p)
}

func (p *MailModel) Del(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	return io.DeleteSvrData(logger, getMailKey(userID))
}
