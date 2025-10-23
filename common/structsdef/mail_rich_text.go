package structsdef

import (
	"maze_game_server/pb/common/MazeCommon"
)

// 邮件面板信息 (MailPanelInfo)
type MailPanelInfo struct {
	Content  string              `json:"content,omitempty"`  // 富文本内容
	Contents []*MailPanelContent `json:"contents,omitempty"` // 邮件内容，包含多条
	User     []*MailUserInfo     `json:"user,omitempty"`     // 用户信息
	Goods    []*MailGoodsInfo    `json:"goods,omitempty"`    // 物品信息
	Value    []*MailValueInfo    `json:"value,omitempty"`    // 数值信息
}

// 邮件面板内容 (MailPanelContent)
type MailPanelContent struct {
	Type    int32  `json:"type,omitempty"`    // 内容类型
	Content string `json:"content,omitempty"` // 内容文本
}

// 简单用户信息结构
type SimpleUserInfo struct {
	UserId    uint64 `json:"user_id,omitempty"`    // 用户ID
	Nickname  string `json:"nickname,omitempty"`   // 昵称
	IconToken int32  `json:"icon_token,omitempty"` // 头像令牌
	UserSex   int32  `json:"user_sex,omitempty"`   // 性别
}

// 邮件用户信息 (MailUserInfo)
type MailUserInfo struct {
	Key   string          `json:"key,omitempty"`   // 通配符键
	Type  int32           `json:"type,omitempty"`  // 用户类型
	User  *SimpleUserInfo `json:"user,omitempty"`  // 用户信息
	Style *MailStyleInfo  `json:"style,omitempty"` // 样式信息
}

// 邮件物品信息 (MailGoodsInfo)
type MailGoodsInfo struct {
	Key   string               `json:"key,omitempty"`   // 通配符键
	Item  *MazeCommon.MazeItem `json:"item,omitempty"`  // 物品信息
	Style *MailStyleInfo       `json:"style,omitempty"` // 样式信息
}

// 邮件数值信息 (MailValueInfo)
type MailValueInfo struct {
	Key   string         `json:"key,omitempty"`   // 通配符键
	Value int64          `json:"value,omitempty"` // 数值
	Style *MailStyleInfo `json:"style,omitempty"` // 样式信息
}

// 邮件样式信息 (MailStyleInfo)
type MailStyleInfo struct {
	Color    string `json:"color,omitempty"`     // 颜色
	Font     string `json:"font,omitempty"`      // 字体
	FontS    uint32 `json:"font_s,omitempty"`    // 字号
	FontW    uint32 `json:"font_w,omitempty"`    // 字重
	IconW    uint32 `json:"icon_w,omitempty"`    // 头像宽度
	IconH    uint32 `json:"icon_h,omitempty"`    // 头像高度
	ShowName bool   `json:"show_name,omitempty"` // 是否显示名称
	ShowIcon bool   `json:"show_icon,omitempty"` // 是否显示图标
}

// 邮件面板内容类型
const (
	MAIL_PANEL_CONTENT_MAIL_RECEIVER = 1 // 邮件接收者
	MAIL_PANEL_CONTENT_MAIL_CONTENT  = 2 // 邮件内容
	MAIL_PANEL_CONTENT_MAIL_FROM     = 3 // 邮件发送者
)

// 特殊奖励配置 (SpecialAwardConfig)
type SpecialAwardConfig struct {
	ConfigId int32  `json:"config_id,omitempty"` // 配置ID
	Desc     string `json:"desc,omitempty"`      // 描述
}

// 随机宝箱信息 (RandomBoxInfo)
type RandomBoxInfo struct {
	BoxId    int32                  `json:"box_id,omitempty"`    // 宝箱ID
	ItemList []*MazeCommon.MazeItem `json:"item_list,omitempty"` // 物品列表
}

// 贸易宝箱信息 (TradeBoxInfo)
type TradeBoxInfo struct {
	BoxId    int32                  `json:"box_id,omitempty"`    // 宝箱ID
	ItemList []*MazeCommon.MazeItem `json:"item_list,omitempty"` // 物品列表
}

// 投票信息 (VoteInfo)
type VoteInfo struct {
	VoteId     int32    `json:"vote_id,omitempty"`     // 投票ID
	Title      string   `json:"title,omitempty"`       // 投票标题
	Options    []string `json:"options,omitempty"`     // 选项列表
	EndTime    int64    `json:"end_time,omitempty"`    // 结束时间
	MaxChoices int32    `json:"max_choices,omitempty"` // 最大选择数
}

// 链接参数 (LinkParam)
type LinkParam struct {
	LinkType int32  `json:"link_type,omitempty"` // 链接类型
	Param    string `json:"param,omitempty"`     // 参数
}

// 仇敌奖励 (EnemyAward)
type EnemyAward struct {
	EnemyId uint64               `json:"enemy_id,omitempty"` // 仇敌ID
	Award   *MazeCommon.MazeItem `json:"award,omitempty"`    // 奖励
}

// 友谊之舟返航奖励信息 (FriendShipAwardInfo)
type FriendShipAwardInfo struct {
	ShipId   int32                  `json:"ship_id,omitempty"`   // 船只ID
	ItemList []*MazeCommon.MazeItem `json:"item_list,omitempty"` // 物品列表
}

// 提现卡地图红包信息 (MailCashCardPackInfo)
type MailCashCardPackInfo struct {
	PackId     int32 `json:"pack_id,omitempty"`     // 红包ID
	Amount     int64 `json:"amount,omitempty"`      // 金额
	ExpireTime int64 `json:"expire_time,omitempty"` // 过期时间
}
