// @Author pangchenyang 2025/5/6 19:07:00
// @Desc: 
package constdef

// 跟客户端约定的错误码
const (
	DE_ERR_SWITCH_SUIT_CD                   = 70001 // 切换装备CD中
	DE_ERR_EQUIP_SUIT_NOT_MATCH             = 70002 // 套装不匹配
	DE_ERR_EQUIP_SUIT_NOT_INIT              = 70003 // 未设置当前装备套
	DE_ERR_EQUIP_DISMANTLE_QUANLITY_WRONG   = 70004 // 装备分解功能选定的品质和选定的装备不匹配
	DE_ERR_RQ_RATE_LIMITER                  = 70005 // 操作太频繁
	DE_ERR_DRESS_EQUIP_DATA_NO_MATCH        = 70008 // 穿戴装备数据不匹配
	DE_ERR_FORCE_PREVIEW_FAIL               = 70009 // 预览武力值失败
	DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS = 70014 // 分解的装备不存在
	// MAZE_ERR_ENERGY_LESS                    = 80000 // 体力不足
	// MAZE_ERR_ENERGY_FULL                    = 80001 // 体力已满
)
