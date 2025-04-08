package structsdef

type KvPair struct {
	K int32
	V int64
}

type KvPairs []*KvPair

// 人偶武力结构
type DollForce struct {
	ShowForce int64 // 展示武力
	RealForce int64 // 真实武力
	PvpForce  int64 // pvp武力
}

// 性别变化消息
type SexChangeInfo struct {
	UserId     string `json:"user_id"`
	Sex        int32  `json:"sex"`
	CreateChg  int32  `json:"create_chg"`
	SourceType int32  `json:"source_type"`
	FromType   int32  `json:"from_type"`
}

// 迷宫用户等级变化流水
type MazeUserLevelRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	OldLevel   int32  `json:"old_level"`   // 旧等级
	NewLevel   int64  `json:"new_level"`   // 新等级
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}
