package structsdef

// 人偶属性计算通知消息
type MazeCalcAttrNotifyMsg struct {
	UserId     uint64 `json:"user_id"`     // 用户Id
	FromServer string `json:"from_server"` // 服务来源(服务类型+名字)
	BuffSrc    int32  `json:"buff_src"`    // buff来源
	ChgType    int32  `json:"chg_type"`    // 变化类型
	ChgDesc    string `json:"chg_desc"`    // 变化原因描述
	Session    string `json:"session"`     // session
	Stamp      int64  `json:"stamp"`       // 消息时间戳 ms
	RetryFlag  int32  `json:"retry_flag"`  // 失败重试用,内部用不用设置
}
