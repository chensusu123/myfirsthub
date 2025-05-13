package websocket_service

type NoramlJsonMsg struct {
	MsgType int         `json:"msg_type"`
	Data    interface{} `json:"data"`
}
