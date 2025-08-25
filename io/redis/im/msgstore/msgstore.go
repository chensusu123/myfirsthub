package msgstore

type Message struct {
	MessageID  uint64 `json:"message_id,omitempty"`
	Type       int32  `json:"type,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	UserID     uint64 `json:"user_id,omitempty"`
}
