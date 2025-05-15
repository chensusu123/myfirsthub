/**
 * Created by GoLand.
 * User: majiange
 * Date: 2018/7/26
 * Time: 下午2:13
 */
package errors

import (
	"errors"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
)

func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// deprecated: xxx
func NewErrorInfo(errCode int64, errMsg string) *MessageType.ErrorInfo {
	info := &MessageType.ErrorInfo{}
	info.ErrCode = proto.Int64(errCode)
	info.ErrMsg = []byte(errMsg)
	return info
}

type ErrCode interface {
	Error() string
	Wrap(string) *MessageType.ErrorInfo
	ToInfo() *MessageType.ErrorInfo
}

type CodeError struct {
	Code int64
	Msg  string
}

func (er *CodeError) Error() string {
	return fmt.Sprintf("error code with %d  info %s", er.Code, er.Msg)
}

func (er *CodeError) Wrap(msg string) *MessageType.ErrorInfo {
	// er.Msg = msg
	return &MessageType.ErrorInfo{ErrCode: proto.Int64(er.Code), ErrMsg: []byte(msg)}
}

func (er *CodeError) ToInfo() *MessageType.ErrorInfo {
	return &MessageType.ErrorInfo{ErrCode: proto.Int64(er.Code), ErrMsg: []byte(er.Msg)}
}

func (er *CodeError) WrapMsg(msg string) *CodeError {
	// er.Msg = msg
	return &CodeError{Code: er.Code, Msg: msg}
}

func NewCodeError(code int64, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func NewCommonCodeError(msg string) *MessageType.ErrorInfo {
	return &MessageType.ErrorInfo{ErrCode: proto.Int64(100001), ErrMsg: []byte(msg)}
}

// !!!警告!!! 错误码的范围50000~110000，否则会和C++的错误码冲突
// !!!警告!!! 错误码的范围50000~110000，否则会和C++的错误码冲突
var (
	NO_ERROR_CODE int64 = 0x80000000
	NO_ERROR            = &MessageType.ErrorInfo{ErrCode: proto.Int64(NO_ERROR_CODE), ErrMsg: []byte("")} // 没有错误
	// 通用错误
	COMMON_ERROR_TIPS = NewCodeError(100001, "")

	// 内部错误
	UNMARSHAL_ERROR      = NewCodeError(100002, "") // 数据解析错误
	MODULE_ERROR         = NewCodeError(100003, "获取数据失败")
	DB_SAVE_ERROR        = NewCodeError(100004, "保存数据错误")
	RANCH_PROTO_ERROR    = NewCodeError(100005, "牧场协议版本号不一致")
	RANCH_ACTIVITY_ERROR = NewCodeError(100006, "牧场协议活动号不一致")

	// 重复请求 token
	REAPEATED_REQUEST_TOKEN = NewCodeError(62091, "请求信息没变化")

	// 材料不够,跳转错误
	GOLD_NOT_ENOUGH = NewCodeError(50001, "金币不够")
	EXP_NOT_ENOUGH  = NewCodeError(50002, "探险经验不够")
	// SILVER_NOT_ENOUGH      = NewCodeError(50003, "钱不足")
	PROSPERITY_NOT_ENOUGH  = NewCodeError(50004, "繁荣度不够")
	BAG_ITEM_NOT_ENOUGH    = NewCodeError(50005, "材料不够")
	DIAMOND_NOT_ENOUGH     = NewCodeError(50006, "钻石不够")
	LEVEL_NOT_ENOUGH       = NewCodeError(50007, "探索等级不够")
	ADD_SILVER_COIN        = NewCodeError(50008, "添加钱失败")
	SILVER_COIN_NOT_ENOUGH = NewCodeError(50009, "绿钞不足")
	QUERY_SILVER_COIN      = NewCodeError(50010, "获取钱信息失败")

	INFLUENCE_NOT_ENOUGH            = NewCodeError(50011, "影响力不足")
	CONFIDENCE_NOT_ENOUGH           = NewCodeError(50012, "信心不足")
	ADD_GOLD_COIN                   = NewCodeError(50013, "添加金币失败")
	NOT_FOUND_BUILD                 = NewCodeError(50014, "找不到对应建筑")
	ENERGY_NOT_ENOUGH               = NewCodeError(50015, "精力不够")
	ERR_VALET_NOT_ENOUGH            = NewCodeError(20054029, "体力不足")
	SYSTEM_HELP_LOAD_LIMIT          = NewCodeError(50016, "帮装总次数达到上限")
	DIAMOND_NOT_SAND                = NewCodeError(50017, "元宝不够")
	SILVER_REACH_MAX                = NewCodeError(50018, "绿钞已达上限")
	SILVER_LIMIT_LESS               = NewCodeError(50019, "绿钞上限不足")
	HEARSAY_NOT_ENOUGH              = NewCodeError(50020, "道听途说不足")
	BEER_NOT_ENOUGH                 = NewCodeError(50021, "啤酒不足")
	LUCKTICK_NOT_ENOUGH             = NewCodeError(50022, "抽奖券不足")
	MONEY_CODE_NOT_REGIST           = NewCodeError(50023, "该业务未注册钻石或者金币或者背包类型")
	BUILD_EXTEND_PERMIT_NOT_ENOUGH  = NewCodeError(50024, "扩地许可证不足")
	NotProcessItemType              = NewCodeError(50025, "未注册处理的道具类型")
	ItemCheckFailure                = NewCodeError(50026, "检查添加上限") /// 改成上限错误码
	BUILD_GEM_NOT_ENOUGH            = NewCodeError(50027, "铭文原盘不足")
	ResourceLayerNotMatch           = NewCodeError(50028, "请升级船等级")
	DARK_STEEL_NOT_ENOUGH           = NewCodeError(50029, "强化数量不足")
	ALTAR_EXP_NOT_ENOUGH            = NewCodeError(50030, "祭坛经验数量不足")
	DECORATION_BUILD_POINT_ENOUGH   = NewCodeError(50031, "装饰物建设点数量不足")
	HIGH_WASH_STONE_NOT_ENOUGH      = NewCodeError(50032, "高级洗练石数量不足")
	MAGIC_WASH_STONE_NOT_ENOUGH     = NewCodeError(50033, "魔法洗练石数量不足")
	NEW_STEEL_ONE_NOT_ENOUGH        = NewCodeError(50034, "优质钢数量不足")
	REFINED_STEEL_NOT_ENOUGH        = NewCodeError(50035, "百炼玄铁数量不足")
	HIGH_STEEL_NOT_ENOUGH           = NewCodeError(50036, "高级玄铁数量不足")
	FETTER_POINT_NOT_ENOUGH         = NewCodeError(50037, "羁绊点数量不足")
	EACH_ADD_ITEM_ERROR             = NewCodeError(50038, "添加物品达到单笔奖励上限")
	ISLAND_EXP_ITEM_NOT_ENOUGH      = NewCodeError(50039, "小岛经验数量不足")
	EQUIP_UPGRADE_ITEM_ERROR        = NewCodeError(50040, "装备升级材料数量不足")
	BUILDING_BUILD_POINT_ITEM_ERROR = NewCodeError(50041, "建筑建设点数量不足")
	GOLD_INGOT_ITEM_ERROR           = NewCodeError(50042, "元宝数量不足")
	ITEM_CHECK_TIME_OUT             = NewCodeError(50043, "检查超时")
	ITEM_ADD_TIME_OUT               = NewCodeError(50044, "道具添加超时")
	ITEM_SUB_TIME_OUT               = NewCodeError(50045, "道具扣除超时")
	ITEM_CHECK_ERROR                = NewCodeError(50046, "检查道具逻辑失败")
	PEARL_NOT_ENOUGH                = NewCodeError(50047, "珍珠道具不足")
	HEX_STONE_NOT_ENOUGH            = NewCodeError(50048, "神石道具不足")
	ITEM_CHECK_ITEM_NOT_ENOUGH      = NewCodeError(50049, "检查道具不足")
	ITEM_EXPLOIT_NOT_ENOUGH         = NewCodeError(50050, "功勋数量不足")
	HONOR_POINTS_NOT_ENOUGH         = NewCodeError(50051, "荣誉点数量不足")
	NotOwnBoxError                  = NewCodeError(50054, "用户没有宝箱了")
	ITEM_ADD_ERROR                  = NewCodeError(50055, "道具添加失败")
	ITEM_ADD_LIMIT                  = NewCodeError(50056, "道具添加添加上限")
	BAG_NOT_EXIST_ERROR             = NewCodeError(50056, "仓库没有初始化")
	ITEM_ADD_LIMIT_PART             = NewCodeError(50058, "因为道具添加达到上限,成功添加一部分")

	// 协议错误
	COST_NOT_MATCH      = NewCodeError(61101, "花费信息不匹配") // 前后端不一致,直接返回正确的
	DATA_NOT_MATCH      = NewCodeError(61102, "配置信息不匹配")
	CONFIG_NOT_FOUND    = NewCodeError(61103, "配置没有找到")
	ARGS_NOT_MATCH      = NewCodeError(61104, "无效参数")
	SPEED_UP            = NewCodeError(61105, "加速失败")
	SPEED_UP_TIME_OVER  = NewCodeError(61106, "时间已到")
	SPEED_UP_COUNT_ZERO = NewCodeError(61107, "加速道具个数0") // 扣0个道具

	// 逻辑错误
	STOREHOUSE_FULL_ERROR       = NewCodeError(71101, "仓库已满")
	STOREHOUSE_ADD_ERROR        = NewCodeError(71102, "添加仓库物品失败")
	STOREHOUSE_ITEM_LIMIT_ERROR = NewCodeError(71103, "添加仓库单个物品数量已达上限")
	BAG_FULL_ERROR              = NewCodeError(71104, "背包已满")     // 老业务的背包
	BAG_ITEM_LIMIT_MAX_ERROR    = NewCodeError(71105, "背包道具已达上限") // 老业务的背包
	GAME_PROTO_SMALL_ERROR      = NewCodeError(71106, "user protocol small 200")
	SPEEDUP_ITEM_NUM_SHORT      = NewCodeError(72101, "加速道具不足")
	ERR_BAG_DATA_NOT_MATCH      = NewCodeError(72102, "背包数据不匹配")
	ERR_BUILDING_DATA_ERROR     = NewCodeError(72103, "建筑未初始化")
	NOT_AUTHORITY               = NewCodeError(61951, "没有操作权限")
	COMMON_BUY_COST_ERR         = NewCodeError(62301, "花费值错误") // 前后端不一致,按后端扣并返回正确的

	SLOT_NOT_UNLOCK           = NewCodeError(61995, "栏位没有解锁")
	SLOT_NIL_ERROR            = NewCodeError(61994, "栏位为空")
	CAN_NOT_USE_POISON_ERROR  = NewCodeError(61992, "保护期内不能使用")
	NOT_NEED_USE_POISON_ERROR = NewCodeError(61993, "没有中毒，不需解毒")

	FORMULA_IN_COOKING = NewCodeError(61923, "料理制作中")
	FORMULA_LOCKED     = NewCodeError(61925, "配方未解锁")
	FORMULA_UNFINISH   = NewCodeError(61927, "制作未完成")
	FORMULA_STAR_MAXED = NewCodeError(61928, "已满星")
	KITCKEN_LEVEL      = NewCodeError(61936, "厨房等级不够")

	CUISINE_ERROR_FULL_HP  = NewCodeError(61938, "你的血量已满，不需要使用")
	CUISINE_ERROR_HAS_DEAD = NewCodeError(61939, "你已阵亡，无法使用")
	CUISINE_ERROR_NO_DEAD  = NewCodeError(61940, "尚未阵亡，无法使用")

	FINDING_AWARD_RECEIVED = NewCodeError(61992, "已领取奖励")
	NO_GRAPEVINE           = NewCodeError(61980, "没有新传闻")
	CLUE_RECEIVED          = NewCodeError(61987, "已获得该线索")
	CLUE_UNCOLLECTED       = NewCodeError(61988, "未集齐线索")
	NPC_IN_WORKING         = NewCodeError(61989, "Npc工作中")
	FINDING_OWNED          = NewCodeError(61990, "发现物已拥有")
	FINDING_NOT_OWNED      = NewCodeError(61991, "发现物未获得")

	TASK_ERROR_STATUS = NewCodeError(62001, "错误的任务状态")
	TASK_ERROR_TYPE   = NewCodeError(62002, "错误的任务类型")

	TASK_NOT_FOUND   = NewCodeError(62010, "任务不存在")
	ADD_SEED_FAILED  = NewCodeError(62011, "加种子碎片失败")
	ADD_MONEY_FAILED = NewCodeError(62012, "加现金红包失败")

	HOPE_SEEK_HELP_COUNT_LIMIT = NewCodeError(62014, "已无求助次数")
	HOPE_ORDER_SAIL            = NewCodeError(62015, "希望号已发出")
	HOPE_ORDER_INVALID         = NewCodeError(62016, "希望号已失效")
	HOPE_ORDER_ITEM_FIN        = NewCodeError(62017, "此货舱已被装完")
	HOPE_ITEM_NO_HELP          = NewCodeError(62018, "求助状态已失效")

	EXTEND_ERROR_MAX_CONCURRENT_NUM = NewCodeError(62019, "已达到同时扩建最大数量")
	EXTEND_ERROR_POPULATION_LESS    = NewCodeError(62020, "人口数量不足")
	EXTEND_ERROR_ISLAND_AREA_LESS   = NewCodeError(62021, "小岛面积不足")

	BUILDING_ERROR_NOT_NEAR_MAIN_ROAD = NewCodeError(62022, "未连接主干道")

	// 希望号新(62023-62050)
	ERR_HOPE_NO_HAS_TIME_AWARD                   = NewCodeError(62023, "no has forward awards")
	ERR_HOPE_CANNOT_SPEEDUP                      = NewCodeError(62024, "cannot speedup")
	ERR_HOPE_CANNOT_SEND                         = NewCodeError(62025, "cannot send")
	ERR_HOPE_NOT_UNLOCK                          = NewCodeError(62026, "希望号暂未解锁")
	ERR_HOPE_CREATE_ORDER                        = NewCodeError(62027, "生成希望号订单失败")
	ERR_HOPE_LEAVE                               = NewCodeError(62028, "希望号已经发出啦")
	ERR_HOPE_OTHER_LOAD                          = NewCodeError(62029, "此货箱已经有人帮装啦")
	ORDER_DAY_HELP_COUNT_LIMIT                   = NewCodeError(62030, "帮助次数已达每日上限")
	ORDER_DAY_HELP_ONE_USER_COUNT_LIMIT          = NewCodeError(62031, "帮助同一好友次数已达每日上限")
	ORDER_DAY_HELP_ONE_USER_ONE_SHIP_COUNT_LIMIT = NewCodeError(62032, "帮助同一好友同一订单已达每日上限")

	// 友谊之舟(62051-62090)
	ERR_SHIP_MAX_COUNT_LIMIT      = NewCodeError(62051, "友谊之舟数量已达上限")
	ERR_SHIP_CREATE_ORDER         = NewCodeError(62052, "生成友谊之舟订单失败")
	ERR_SHIP_LEAVE                = NewCodeError(62053, "友谊之舟已经发出啦")
	ERR_SHIP_BOX_FIN              = NewCodeError(62054, "此货舱已被装完")
	ERR_SHIP_BOX_NOT_IN_SEEK_HELP = NewCodeError(62055, "未在求助中")
	ERR_NOT_UNLOCK                = NewCodeError(62056, "未解锁")

	// 建筑(62091-62100)
	ERR_BUILD_NO_CHG  = NewCodeError(62091, " building list no chg")
	ERR_BUILD_MAX_CNT = NewCodeError(62092, "建造已达上限")

	// 船只(62100-62150)
	ERR_SHIP_NOT_EXIST                     = NewCodeError(62100, "船只不存在")
	ERR_SHIP_FULL_LEVEL                    = NewCodeError(62101, "船已满级")
	ERR_SHIP_FULL_STAR                     = NewCodeError(62102, "船已满星")
	ERR_SHIP_FULL_DUR                      = NewCodeError(62103, "耐久已满")
	ERR_SHIP_STATE_ERR                     = NewCodeError(62104, "船只正在使用中")
	ERR_SHIP_DAMAGE                        = NewCodeError(62105, "船只已损坏")
	ERR_SHIP_EXIST                         = NewCodeError(62106, "船只已存在")
	ERR_SHIP_HOUSE_FULL                    = NewCodeError(62107, "船坞已满")
	ERR_SHIP_USING                         = NewCodeError(62111, "船只正在使用中")
	ERR_SHIP_TYPE_LESS                     = NewCodeError(62112, "此类型船只没有装配")
	ERR_SHIP_LOAD_LESS                     = NewCodeError(62113, "此船只没有装配")
	ERR_SHIP_WHARF_LESS                    = NewCodeError(62114, "同时出航船只数量已达上限,请升级科技")
	ERR_SHIP_CNNOT_SPEED_STATE             = NewCodeError(62115, "当前状态不能加速")
	ERR_SHIP_LOAD_LESS_Captain             = NewCodeError(62116, "此船只没有装配船长")
	ERR_SHIP_LOAD_LESS_Helmsman            = NewCodeError(62117, "此船只没有装配舵手")
	ERR_SHIP_LOAD_LESS_Gunner              = NewCodeError(62118, "此船只没有装配炮手")
	ERR_SHIP_STRENGTH_LESS                 = NewCodeError(62119, "船只体力不足")
	ERR_SHIP_Captain_STRENGTH_LESS         = NewCodeError(62120, "船长体力不足")
	ERR_SHIP_Helmsman_STRENGTH_LESS        = NewCodeError(62121, "舵手体力不足")
	ERR_SHIP_Gunner_STRENGTH_LESS          = NewCodeError(62122, "炮手体力不足")
	ERR_SHIP_SAILOR_STRENGTH_LESS          = NewCodeError(62123, "航海士体力不足")
	ERR_SHIP_SAILOR_STRENGTH_COST          = NewCodeError(62124, "扣除航海士体力失败")
	ERR_SHIP_SPEED_ITEM_TYPE               = NewCodeError(62125, "加速道具类型不对")
	ERR_SHIP_NO_HAS_CABIN                  = NewCodeError(62126, "船只没有这个仓位")
	ERR_SHIP_CABIN_LOCK                    = NewCodeError(62127, "仓位没有解锁")
	ERR_SHIP_NO_NEED_FIX                   = NewCodeError(62128, "没有需要维修的船只")
	ERR_SHIP_REPAIRING                     = NewCodeError(62129, "船只正在修理中")
	ERR_SHIP_GATHER_MAX_LIMIT              = NewCodeError(62130, "采矿船只达到上限")
	ERR_SHIP_REFORM_NOT_EXIST              = NewCodeError(62131, "还没有改造信息")
	ERR_SHIP_GATHER_FAMILY_MAX_LIMIT       = NewCodeError(62132, "采家族资源船只达到上限")
	ERR_SHIP_COST_ACTION_POWER             = NewCodeError(62133, "扣除行动力失败")
	ERR_SHIP_REFORM_ATTR_FULL              = NewCodeError(62134, "所有属性已改造满值")
	ERR_SHIP_CABIN_CANNOT_INSTALL_GEMSTONE = NewCodeError(62135, "此舱室不能镶嵌宝石")
	ERR_SHIP_TECH_LV_LESS                  = NewCodeError(62136, "科技等级不足")
	ERR_SHIP_EXPIRE                        = NewCodeError(62137, "船只已过期")
	ERR_SHIP_ACTION_FULL                   = NewCodeError(62138, "行动力已满,不需要补充")
	ERR_SHIP_ACTION_BUY_COUNT_LIMIT        = NewCodeError(62140, "每日购买次数达到上限")
	ERR_SHIP_ATTR_FULL_LV                  = NewCodeError(63141, "船属性强化已满级")
	ERR_SHIP_LV_LESS                       = NewCodeError(63142, "船等级不足")
	ERR_SHIP_HOUSE_CITY_LV_LESS            = NewCodeError(63143, "岛主府等级不足")
	ERR_SHIP_HAS_FLAG_SHIP                 = NewCodeError(63144, "此船已设置为旗舰船")
	ERR_SHIP_FIT_TYPE_NOT_MATCH            = NewCodeError(63145, "装备类型不匹配")
	ERR_SHIP_FULL_STAGE                    = NewCodeError(63146, "已经满阶")
	ERR_SHIP_FULL_EQUIP_BAG                = NewCodeError(63147, "船装备空间已满")
	ERR_SHIP_LESS_EQUIP_BAG                = NewCodeError(63148, "船装备空间不足")
	ERR_SHIP_UPGRADE_GRP_BIG               = NewCodeError(63149, "旗舰与其他船只等级差距大于上限")
	ERR_SHIP_LV_MORE_THAN_FLAG_SHIP        = NewCodeError(63150, "当前船只等级大于旗舰船只等级")
	ERR_SHIP_EQUIP_STRENGTHEN_MAX_LV       = NewCodeError(63151, "装备已强化满级")
	ERR_SHIP_EQUIP_STRENGTHEN_NO_CONFIG    = NewCodeError(63152, "此装备不能强化")
	ERR_GEM_RECAST_COST_LESS               = NewCodeError(63153, "重铸材料不足")

	// 出航  64100~64200
	ERR_SHIP_SAIL_WAIT_RESULT = NewCodeError(64100, "等待出航结果")
	ErrShipCloseSail          = NewCodeError(64101, "海上风浪很大,现在无法出航")
	ErrNoBackBussServer       = NewCodeError(64102, "未找到后端业务")

	// 航海士(62150-62200)
	ERR_SAILOR_NOT_UNLOCK                = NewCodeError(62150, "sailor not unlock")
	ERR_SAILOR_FULL_LEVEL                = NewCodeError(62151, "航海士已满级")
	ERR_SAILOR_UPGRADE_EXP_NOT_ENOUGH    = NewCodeError(62152, "升级经验不足")
	ERR_SAILOR_FULL_STAR                 = NewCodeError(62153, "航海士已满星")
	ERR_SAILOR_ALREADY_EXIST             = NewCodeError(62154, "航海士已合成")
	ERR_SAILOR_ALREADY_LOAD              = NewCodeError(62155, "航海士已装配")
	ERR_SAILOR_UPGRADE_LEVEL_NOT_ENOUGH  = NewCodeError(62156, "升级航海士等级不足")
	ERR_SAILOR_SKILL_FULL_LEVEL          = NewCodeError(62157, "航海士技能已满级")
	ERR_SAILOR_NOT_FIRE                  = NewCodeError(62158, "航海士不能解雇")
	ERR_SAILOR_NOT_RECRUIT               = NewCodeError(62159, "该航海士不能招募")
	ERR_SAILOR_LEVEL_TECH_LIMIT          = NewCodeError(62160, "航海士等级已达科技上限")
	ERR_SAILOR_SKILL_LEVEL_TECH_LIMIT    = NewCodeError(62161, "技能等级已达科技上限")
	ERR_SAILOR_UPGRADE_FAST              = NewCodeError(62162, "")
	ERR_SAILOR_UPGRADE_STAR_NOT_ENOUGH   = NewCodeError(62163, "升星卡片不足")
	ERR_SAILOR_STAR_NOT_MAIN_CARD        = NewCodeError(62164, "该卡不是主卡")
	ERR_SAILOR_ALREADY_HIRE              = NewCodeError(62165, "该航海士已被雇用")
	ERR_SAILOR_UPGRADE_LIKING_NOT_ENOUGH = NewCodeError(62166, "升级好感度不足")
	ERR_SAILOR_UPGRADE_TYPE_NOT_CORRECT  = NewCodeError(62167, "航海士升级类型错误")

	// 资源工厂（矿场，林场）
	ERR_RES_FACTORY_NO_MATCH_WORKER = NewCodeError(62180, "没有匹配的工人")
	ERR_RES_FACTORY_REPEATED_HIRE   = NewCodeError(62181, "重复雇佣同一个工人")
	ERR_RES_FACTORY_WORKER_FULL     = NewCodeError(62182, "工人数量已达上限")
	ERR_RES_FACTORY_GET_IDLE_FAIL   = NewCodeError(62183, "获取材料信息失败")
	ERR_RES_FACTORY_HOUSE_FULL      = NewCodeError(62184, "物品上限 科技限制")
	ERR_SHIP_EXP_FULL               = NewCodeError(62185, "航海术已达上限")
	ERR_MERITORIOUS_SERVICE_FULL    = NewCodeError(62186, "功勋已达上限")

	// 章节任务
	ERR_CHAPTER_TASK_ALL_FINISH = NewCodeError(63002, "已经完成所有任务")
	// 酒馆 63003-63153
	ERR_PUB_DAILY_REFRESH_UPTO_LIMIT = NewCodeError(63003, "单日刷新达到上限")
	// 邮件
	ERR_ANNEX_HAS_OPEN            = NewCodeError(63004, "信封奖励已领取")
	ERR_MAIL_NOT_FOUND            = NewCodeError(63005, "获取不到对应的信封")
	ERR_MAIL_EXPIRE               = NewCodeError(63006, "信封已经过期")
	ERR_MAIL_NOT_FOUND_ANNEX      = NewCodeError(63007, "该信封没有奖励")
	ERR_MAIL_GET_ANNEX_ERROR      = NewCodeError(63008, "获取奖励信息失败")
	ERR_MAIL_ANNEX_ADD_ITEM_ERROR = NewCodeError(63009, "领取附件时，加物品失败")
	ERR_MAIL_ALREADY_EXIST        = NewCodeError(63010, "信封已存在")
	ERR_MAIL_SEND                 = NewCodeError(63011, "信封发送失败")
	ERR_AWAD_MAIL_EMPTY           = NewCodeError(63012, "信封已全部领取")

	// 科技63154-科技63204
	ERR_TECH_NODE_NOT_UNLOCK = NewCodeError(63154, "科技节点尚未解锁")
	ERR_TECH_NODE_FULL_LV    = NewCodeError(63155, "科技节点已经满级")
	ERR_TECH_NODE_DO_UPGRADE = NewCodeError(63156, "升级中")
	ERR_TECH_NODE_PRE_COND   = NewCodeError(63157, "有前置条件未满足")

	// 地图
	ERR_WORLD_MAP_FAMILY = NewCodeError(63158, "用户未加入家族")

	// npc贸易
	// 63161-63180
	ERR_NPC_TRADE_DAY_COUNT_LESS     = NewCodeError(63161, "今天贸易次数已用尽")
	ERR_NPC_TRADE_BOX_STATE          = NewCodeError(63162, "当前不可进行此操作")
	ERR_NPC_TRADE_ORDER_INVALID      = NewCodeError(63163, "订单已经过期")
	ERR_NPC_TRADE_NO_SEEK_HELP_COUNT = NewCodeError(63164, "求助次数已用尽")

	ERR_TECH_NODE_NOT_UPGRADE_DOING = NewCodeError(63170, "科技节点不在升级状态")
	ERR_TECH_NODE_LV_NOT_MATCH      = NewCodeError(63171, "科技节点等级不匹配")

	// 战令
	// 63181-63199
	ERR_TOKEN_COMPOSE_CANNOT_COMPOSE         = NewCodeError(63181, "此战令无法合成")
	ERR_TOKEN_COMPOSE_COSTS_TOKEN_NOT_ENOUGH = NewCodeError(63182, "需要消耗的战令不足")
	ERR_TOKEN_COMPOSE_TOKEN_COMPOSE_FAIL     = NewCodeError(63183, "合成战令失败")
	ERR_TOKEN_COMPOSE_TOKEN_NUMBER_INVALID   = NewCodeError(63184, "合成数量无效")

	// 盟战 buff
	// 63200-65219
	ERR_LEAGUE_BATTLE_BUFF_NOT_START         = NewCodeError(63200, "活动未开启")
	ERR_LEAGUE_BATTLE_BUFF_CFG_NOT_FOUND     = NewCodeError(63201, "盟战配置未找到")
	ERR_LEAGUE_BATTLE_BUFF_PERMISSION_DENIED = NewCodeError(63202, "权限不足")
	ERR_LEAGUE_BATTLE_BUFF_USER_MSG_ERR      = NewCodeError(63203, "用户信息错误")

	ERR_COMMAND_POWER_NOT_MATCH     = NewCodeError(63204, "战力输入错误")
	ERR_COMMAND_FAMILY_ALREADY_HAS  = NewCodeError(63205, "家族已有指挥命令")
	ERR_COMMAND_MATCH_USERS_ZERO    = NewCodeError(63206, "没有用户匹配")
	ERR_COMMAND_POINT_ERR           = NewCodeError(63207, "投影点数据错误")
	ERR_COMMAND_BRICK_ADD_NOT_START = NewCodeError(63208, "加成活动未开启")

	// 宝石合成失败
	ERR_GEM_MERGE_FAIL = NewCodeError(65220, "宝石合成失败")
	// 前置区域杂草未除完
	ERR_EXPAND_NOT_FINISH       = NewCodeError(65221, "前置区域杂草未除完")
	ERR_RECEIVE_EXPAND_RED_FAIL = NewCodeError(65222, "扩地红包领取失败")

	ERR_GEM_COUNT_LESS    = NewCodeError(65223, "宝石数量不足")
	ERR_ITEM_COUNT_LESS   = NewCodeError(65224, "物品数量不足")
	ERR_BUILD_ORDER_STATE = NewCodeError(65225, "建造中不能提交订单")

	// 组队服务
	ERR_SESSINON_ALREADY_EXIST = NewCodeError(66110, "当前session已经存在")

	// 初始化服务 66210-66310
	ErrorNotifyServerName = NewCodeError(66210, "server_name不合法")
	ErrorInit             = NewCodeError(66211, "初始化失败")
	ErrorInitIsRunning    = NewCodeError(66212, "正在初始化")
	ErrorInited           = NewCodeError(66213, "已经初始化过")
)

// deprecated: xxx
var ERROR_GET_BAG_GOODS = NewErrorInfo(100, "获取背包数据失败")

// deprecated: xxx
var ERROR_UNMARSHAL = NewErrorInfo(101, "解包失败")

// deprecated: xxx
var ERROR_SERVER = NewErrorInfo(102, "服务器内部错误")

// deprecated: xxx
var ERROR_GOLD_COIN_NOT_ENOUGH = NewErrorInfo(107, "金币不足")

// deprecated: xxx
var ERROR_DIAMOND_NOT_EQUAL = NewErrorInfo(112, "钻石消耗不一致")

// 个人探索
var (
	// deprecated: xxx
	EXPLORE_ORDIER_FINISH = NewErrorInfo(104, "订单已完成")
	// deprecated: xxx
	EXPLORE_TEAM_RUNNLING = NewErrorInfo(105, "组队探索中")
	// deprecated: xxx
	EXPLORE_EVENT_NO_CLUE = NewErrorInfo(106, "该事件没有对应线索")
	// deprecated: xxx
	EXPLORE_SELT_TEAM_QUIT = NewErrorInfo(108, "您已退出狩猎")
	// deprecated: xxx
	EXPLORE_TARGET_TEAM_QUIT = NewErrorInfo(109, "对方退出狩猎，请换个队伍")
	// deprecated: xxx
	EXPLORE_PVP_BATTLE_CD = NewErrorInfo(110, "对方3分钟后可继续战斗，请换个队伍试试吧")
)

// 组队探索
var (
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_ARGS = NewErrorInfo(61950, "参数错误")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_AUTHORITY = NewErrorInfo(61951, "无权限")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_GET_TEAM_INFO = NewErrorInfo(61952, "队伍不存在")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_GET_USER_INFO = NewErrorInfo(61952, "获取用户信息失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_GET_USER_TEAM_INFO = NewErrorInfo(61953, "获取用户组队信息失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_OVER_TEAM_EXPLORE_LIMIT = NewErrorInfo(61955, "超过组队探索次数限制")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_ALREADY_IN_ONE_TEAM = NewErrorInfo(61956, "已经在一个队伍中，不能进行此操作")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_GENERATE_TEAM_ID = NewErrorInfo(61957, "生成新的队伍ID失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_ADD_TEAM = NewErrorInfo(61959, "添加新的队伍失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_ADD_TEAM_MEMBER = NewErrorInfo(61960, "添加队伍成员失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_START_TEAM_EXPLORE = NewErrorInfo(61961, "队伍已经开始探索")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_TEAM_IS_FULL = NewErrorInfo(61962, "队伍人员已满")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_GET_CAR_INFO = NewErrorInfo(61963, "获取用户车辆信息失败")
	// deprecated: xxx
	ERROR_TEAM_EXPLORE_SET_PVP_ONCE = NewErrorInfo(61964, "PVP选项只能设置一次")
)

// 车辆相关
var (
	// deprecated: xxx
	ERROR_CAR_INFO = NewErrorInfo(61970, "获取车辆信息失败")
	// deprecated: xxx
	ERROR_CAR_RQ_PARAM = NewErrorInfo(61971, "请求参数错误")
	// deprecated: xxx
	ERROR_CAR_NO_FOUND_EQUIPMENT = NewErrorInfo(61972, "未找到装备位信息")
	// deprecated: xxx
	ERROR_CAR_EQUIPMENT_LV = NewErrorInfo(61973, "当前等级与请求参数不符")
	// deprecated: xxx
	ERROR_CAR_EQUIPMENT_FULL = NewCodeError(61974, "当前装备位已经满级")
	// deprecated: xxx
	ERROR_CAR_EQUIPMENT_CONFIG = NewErrorInfo(61975, "装备位配表错误")
	// deprecated: xxx
	ERROR_CAR_UPGRADE_COST = NewErrorInfo(61976, "扣除消耗失败")
	// deprecated: xxx
	ERROR_CAR_UPGRADE = NewErrorInfo(61977, "升级装备位失败")
	// deprecated: xxx
	ERROR_CAR_SKIN_NOT_EXIST = NewErrorInfo(61978, "皮肤不存在")
	// deprecated: xxx
	ERROR_CAR_SKIN_UNLOCK = NewErrorInfo(61979, "皮肤未解锁")
	// deprecated: xxx
	ERROR_CAR_REPLACE_SKIN = NewErrorInfo(61980, "更换皮肤失败")
	// deprecated: xxx
	ERROR_CAR_SKIN_CLEAR_NEW_FLAG = NewErrorInfo(61981, "清除皮肤NEW失败")
	// deprecated: xxx
	ERROR_CAR_SKIN_OVER_MAX_DIFF = NewErrorInfo(61982, "升级超过最大等级差")
	// deprecated: xxx
	ERROR_CAR_STAR_FULL = NewCodeError(61983, "当前星级已经满级")
)

// 炼金
var (
	// deprecated: xxx
	ERROR_GET_ALCHEMY_BUILDING = NewErrorInfo(62050, "未获取到building")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_MANAGER = NewErrorInfo(62051, "获取炼金模块失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_GET_BAG_INFO = NewErrorInfo(62055, "获取背包信息失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_COST_BAG_GOODS = NewErrorInfo(62056, "扣除背包物品失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_GENERATE_SPAR = NewErrorInfo(62057, "生成晶石失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_SAVE_SPAR = NewErrorInfo(62058, "存储晶石失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_DRAW_ITEM_ERR = NewErrorInfo(62059, "扣除炼金房物品失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_ADD_BAG_ITEM_ERR = NewErrorInfo(62060, "背包添加晶石失败")
	// deprecated: xxx
	ERROR_GET_ALCHEMY_RANGE_CFG_ERR = NewErrorInfo(62061, "获取炼金生成晶石范围出错")
)

// 私聊
var (
	// deprecated: xxx
	ERROR_GET_CHAT_AREA_NIL = NewErrorInfo(62080, "聊天类型错误")
	// deprecated: xxx
	ERROR_GET_CHAT_MSG_NIL = NewErrorInfo(62081, "聊天消息为空")
)

// 好友列表
var (
	// deprecated: xxx
	ERROR_GET_FRIEND_LIST_MGR = NewErrorInfo(62080, "获取好友列表模块失败")
)

// Npc Trade
var (
	ERROR_NPC_TRADE_RECEIVE_FAIL = NewErrorInfo(62090, "领取奖励失败")
)

// 扩地 家园海盗
var (
	ERROR_AREA_NOT_MATCH            = NewErrorInfo(62110, "地块不匹配")
	ERROR_AREA_PIRATE_NOT_DEAD      = NewErrorInfo(62111, "海盗未打完")
	ERROR_AREA_AWARD_NOT_OPEN       = NewCodeError(62112, "扩地奖励未领取完")
	ERROR_AREA_GET_AWARD_FAIL       = NewCodeError(62113, "扩地奖励领取失败")
	ERROR_AREA_SHIP_ITEM_CHECK_FAIL = NewCodeError(62113, "船装备已满")
)

// 订单 62120 - 62130
var (
	ERROR_ORDER_LIMIT                 = NewCodeError(62120, "今日订单数已达上限")
	ERROR_AREA_AWARD_NIL              = NewCodeError(62113, "无未领取的奖励")
	ERROR_AREA_POS_NOT_FIND           = NewCodeError(62114, "未生成扩地工期")
	ERROR_NEED_REFRESH_ORNAMENT_ORDER = NewCodeError(62115, "需要刷新装饰物订单")
)

// 扩地 62210 - 62220
var (
	ERROR_EXTEND_COUNT_LIMIT = NewCodeError(62211, "扩地次数不足")
)

var (
	CheckShipBuildingError  = errors.New("检查港口失败")
	ShipBuildingUnlockError = errors.New("港口尚未解锁")
)
