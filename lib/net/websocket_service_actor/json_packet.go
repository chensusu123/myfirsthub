package websocket_service_actor

type NoramlJsonMsg struct {
	MsgType int         `json:"msg_type"`
	Data    interface{} `json:"data"`
}
