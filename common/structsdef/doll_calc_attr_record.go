/*
 * @Author: majian
 * @Date: 2024-07-11 11:29:58
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 11:27:30
 */
package structsdef

import "maze_game_server/io/kafka/kafkacommonstruct"

type KafkaCommon = kafkacommonstruct.KafkaCommon

type MazeGameAttrChgRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`           //用户Id
	GroupId    uint32 `json:"group_id" gorm:"column:group_id"`         //分组Id
	AttrId     int32  `json:"attr_id" gorm:"column:attr_id"`           //属性ID
	AttrType   int32  `json:"attr_type" gorm:"column:attr_type"`       //属性类型
	NewVal     int64  `json:"new_val" gorm:"column:new_val"`           //新值
	OldVal     int64  `json:"old_val" gorm:"column:old_val"`           //旧值
	ChgType    int32  `json:"chg_type" gorm:"column:chg_type"`         //变化类型
	ChgSubType int32  `json:"chg_sub_type" gorm:"column:chg_sub_type"` //变化子类型
	ChgDesc    string `json:"chg_desc" gorm:"column:chg_desc"`         //原因描述
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`   //时间戳 ms
	Extra      string `json:"extra" gorm:"column:extra"`               // 扩展信息
}

// 武力属性变化流水
type DollForceAttrChgRecord struct {
	UserId     uint64 `json:"user_id"`     //用户Id
	GroupId    uint32 `json:"group_id"`    //分组Id
	NewVal     string `json:"new_val"`     //新值
	OldVal     string `json:"old_val"`     //旧值
	SrcType    int32  `json:"src_type"`    //来源类型
	ChgReason  int32  `json:"chg_reason"`  //变化原因
	ChgDesc    string `json:"chg_desc"`    //原因描述
	CreateTime int64  `json:"create_time"` //时间戳 ms
}

// 非武力属性变化流水
type MazeGameBuffAttrChgRecord struct {
	UserId     uint64 `json:"user_id"`     //用户Id
	GroupId    uint32 `json:"group_id"`    //分组Id
	NewVal     string `json:"new_val"`     //新值
	OldVal     string `json:"old_val"`     //旧值
	SrcType    int32  `json:"src_type"`    //来源类型
	ChgReason  int32  `json:"chg_reason"`  //变化原因
	ChgDesc    string `json:"chg_desc"`    //原因描述
	CreateTime int64  `json:"create_time"` //时间戳 ms
}

type AttrChgInfo struct {
	AttrId int32 `json:"attr_id"`
	OldVal int64 `json:"old_val"`
	CurVal int64 `json:"cur_val"`
}

// 人偶属性变化通知
type DollAttrChgNotify struct {
	UserId     uint64         `json:"user_id"`      //用户Id
	GroupId    uint32         `json:"group_id"`     //分组Id
	ChgAttrs   []*AttrChgInfo `json:"chg_attrs"`    //变化的属性
	Session    string         `json:"session"`      // session
	ChgType    int32          `json:"chg_type"`     //变化类型
	ChgSubType int32          `json:"chg_sub_type"` //变化子类型
	ChgDesc    string         `json:"chg_desc"`     //原因描述
	CreateTime int64          `json:"create_time"`  //时间戳 ms
}

// 迷宫临时buff变化通知
type MazeTempBuffChangeMsg struct {
	UserId     uint64         `json:"user_id"`     // 用户Id
	GroupId    uint32         `json:"group_id"`    // 分组Id
	StageId    int32          `json:"stage_id"`    // 关卡id
	ChgAttrs   []*AttrChgInfo `json:"chg_attrs"`   // 变化的属性
	ChgType    int32          `json:"chg_type"`    // 变化类型
	ChgDesc    string         `json:"chg_desc"`    // 原因描述
	CreateTime int64          `json:"create_time"` // 时间戳 ms
}
