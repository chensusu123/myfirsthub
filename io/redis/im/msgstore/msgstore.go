package msgstore

type Message struct {
	MessageID  uint64 `json:"message_id,omitempty"`
	Type       int32  `json:"type,omitempty"`
	Content    []byte `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	HasRead    bool   `json:"has_read,omitempty"`
}

type Session struct {
	SessionID  uint64 `json:"session_id,omitempty"`
	Type       int32  `json:"type,omitempty"`
	PeerID     uint64 `json:"peer_id,omitempty"`
	GroupID    int64  `json:"group_id,omitempty"`
	Top        bool   `json:"top,omitempty"`
	Recent     []Message
	CreateTime int64 `json:"create_time,omitempty"`
}

type SessionList struct {
	Sessions []Session
}
