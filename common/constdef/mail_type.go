package constdef

// 邮件大类 (MailClassType)
type MailClassType int32

const (
	SYSTEM_MAIL      MailClassType = 0 // 系统类信封
	INTERACTION_MAIL MailClassType = 1 // 互动类信封
)

// 邮件类型 (MailType)
type MailType int32

const (
	SYSTEM                 MailType = 1   // 系统邮件
	BATTLE                 MailType = 2   // 战斗邮件
	TEAM                   MailType = 3   // 组队邮件
	SWITCH_FAMILY          MailType = 4   // 家族成员位置
	FAMILY_EXPAND          MailType = 5   // 家族扩地
	FAMILY_RECRUIT         MailType = 6   // 家族招募
	SHOPPING_PACK          MailType = 7   // 商城礼包
	MASTER_TRANSFER        MailType = 8   // 族长转让
	MASTER_BUILD           MailType = 9   // 建筑类型
	LEAGUE_TRADE           MailType = 10  // 盟运
	TEAM_MINING            MailType = 11  // 组队采集
	INQUIRE_SAILOR         MailType = 12  // 打探航海士
	LEAGUE_BRICK           MailType = 13  // 联盟搬砖
	TRAIN_ORDER            MailType = 14  // 小火车订单
	INTERCEPT_NOTIFY       MailType = 15  // 拦截通知
	INTERCEPT_THANKS       MailType = 16  // 拦截感谢
	HELP_THANKS            MailType = 17  // 帮助感谢信
	RESOURCE_WAR           MailType = 18  // 资源战
	ALLIANCESIGNIN         MailType = 19  // 微信小程序联盟签到
	RANK_ACTIVITY          MailType = 20  // 排行榜奖励
	SEASON_FINALS_CHAMPION MailType = 21  // 赛季决赛冠军联盟奖励
	SEASON_LEADER          MailType = 22  // 赛季领军者通知
	MAIL_COMMON            MailType = 100 // 通用邮件类型
)

// 邮件子类型 (MailSubType)
type MailSubType int32

const (
	// 奖励类
	AWARD MailSubType = 1

	// 组队相关
	TEAM_INVITAE MailSubType = 2
	TEAM_JOIN    MailSubType = 3
	TEAM_QUIT    MailSubType = 4

	// 家族相关
	SWITCH_FAMILY_POS   MailSubType = 8
	FAMILY_EXPAND_CALL  MailSubType = 9
	FAMILY_RECRUIT_CALL MailSubType = 10

	// 贸易相关
	TRANS_GRAB      MailSubType = 21
	TRANS_INTERCEPT MailSubType = 22

	// 战斗相关
	DOLL_BATTLE_ROB_SUCCESS_MAIL MailSubType = 84
)

// 邮件内容类型 (MailContentType)
type MailContentType int32

const (
	CONTENT_TEXT                    MailContentType = 1  // 文本
	CONTENT_AWARD                   MailContentType = 2  // 奖励类
	CONTENT_OPERATA                 MailContentType = 3  // 操作类
	CONTENT_WAR_REPORT              MailContentType = 4  // 战报类
	CONTENT_FAMILY_TEAM             MailContentType = 5  // 宝箱奖励类
	CONTENT_GIFT_PACK               MailContentType = 6  // 礼包奖励类型
	CONTENT_CLIENT_OPERATA          MailContentType = 7  // 客户端操作类
	CONTENT_AWARD_AND_INVITE_FRIENT MailContentType = 8  // 领取奖励，并要求用户加好友
	CONTENT_TRAIN_BOX               MailContentType = 9  // 小火车宝箱
	CONTENT_TRADE_BOX               MailContentType = 10 // 贸易宝箱类型
	CONTENT_NEW_WAR_REPORT          MailContentType = 11 // 新类型战报
	CONTENT_VOTE                    MailContentType = 12 // 投票功能
	CONTENT_TEXT_AND_AWARD          MailContentType = 13 // 文件带显示奖励
	FRIEND_SHIP_RECEIVE_AWARD       MailContentType = 14 // 友谊之舟返航领奖
	CONTENT_DOLL_BATTLE_REPORT      MailContentType = 15 // 人偶战报
)

// 邮件广播标志 (MailBroadcastFlag)
type MailBroadcastFlag int32

const (
	MAIL_NOT_BROADCAST    MailBroadcastFlag = 0 // 不广播
	MAIL_FAMILY_BROADCAST MailBroadcastFlag = 1 // 家族广播
	MAIL_LEAGUE_BROADCAST MailBroadcastFlag = 2 // 联盟广播
	MAIL_MAP_BROADCAST    MailBroadcastFlag = 3 // 地图广播
)

// 邮件操作类型 (MailOperateType)
type MailOperateType int32

const (
	CommandJump MailOperateType = 1 // 指令跳转
)

// 邮件点击类型 (MailClickType)
type MailClickType int32

const (
	ACTION_JS_BRIDGE    MailClickType = 0 // js-bridge
	ACTION_REQUEST_HTTP MailClickType = 1 // 访问某个http链接
	ACTION_SERVER       MailClickType = 2 // 访问固定服务
	ACTION_NOT_CLICK    MailClickType = 3 // 按钮不让点击
	ACTION_COCOS        MailClickType = 4 // cocos处理
)

// 邮件操作服务ID (MailOpServerId)
type MailOpServerId int32

const (
	MAIL_OP_SERVER_TRAIN       MailOpServerId = 1 // 火车
	MAIL_OP_SERVER_FAMILYTEAM  MailOpServerId = 2 // 组队寻宝
	MAIL_OP_SERVER_PHONE       MailOpServerId = 3 // 电话
	MAIL_OP_SERVER_TRAIN_ORDER MailOpServerId = 4 // 小火车订单服务
	MAIL_OP_SERVER_OFFICIAL    MailOpServerId = 5 // 官员服务
	MAIL_OP_WORLD_GIFT_PACK    MailOpServerId = 6 // 世界礼包赠送
	MAIL_OP_SERVER_MAIL        MailOpServerId = 7 // 通知信封服务发送感谢信
	MAIL_OP_SERVER_VOTE        MailOpServerId = 8 // 投票功能
)

// 附件类型 (AnnexType)
type AnnexType int32

const (
	AnnexType_normal                   AnnexType = 1 // 普通奖励
	AnnexType_high                     AnnexType = 2 // 高级奖励
	AnnexType_high_v2                  AnnexType = 3 // 高级奖励类型2(家族战奖励)
	AnnexType_sys_league_trade_offense AnnexType = 4 // 盟运系统发船进攻方奖励
	AnnexType_sys_league_trade_defense AnnexType = 5 // 盟运系统发船防守方奖励
)

// 邮件标签 (MailLabel)
type MailLabel int32

const (
	MailLabelSystem   MailLabel = 0 // 系统
	MailLabelFamily   MailLabel = 1 // 家族
	MailLabelAlliance MailLabel = 2 // 联盟
)

// 邮件面板内容类型
const (
	MAIL_PANEL_CONTENT_MAIL_RECEIVER = 1 // 邮件接收者
	MAIL_PANEL_CONTENT_MAIL_CONTENT  = 2 // 邮件内容
	MAIL_PANEL_CONTENT_MAIL_FROM     = 3 // 邮件发送者
)
