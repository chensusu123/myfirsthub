package codec

import (
	"bytes"
	"fmt"
	"io"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/frame"
	"maze_game_server/lib/nano/serialize"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Message struct {
	PackType  uint16 `json:"pack_type,omitempty"`
	SessionID uint32 `json:"session_id,omitempty"`
	EsRqTime  uint64 `json:"es_rq_time,omitempty"`
	EsRsTime  uint64 `json:"es_rs_time,omitempty"`
	Data      string `json:"data,omitempty"`
}

// TODO 需要补全日志
type JsonPacketCodec struct {
	rts map[uint16]target
	ser serialize.Serializer
	buf *bytes.Buffer
}

func NewJsonPacketCodec(routes *Routes, opts ...CodecOption) (c *JsonPacketCodec) {
	c = &JsonPacketCodec{
		rts: make(map[uint16]target),
	}
	// Set options
	for _, setOpt := range opts {
		setOpt(c)
	}
	// Init route table
	c.initRoutes(routes)
	return
}

// initRoutes 初始化RqID与处理接口的映射关系(路由表)
func (c *JsonPacketCodec) initRoutes(routes *Routes) {
	for _, comp := range routes.comps {
		s := component.NewService(comp, nil)

		// Extract process handler
		if err := s.ExtractHandler(); err != nil {
			panic(err)
		}

		for fnName := range s.Handlers {
			// 在nano内置组件导出功能的基础上，只解析带有RqID/RsID的定义格式：Game.Welcome_1_2
			rqID, rsID, ok := route(fnName)
			if ok {
				c.rts[rqID] = target{
					rsID:    rsID,
					handler: fmt.Sprintf("%s.%s", s.Name, fnName),
				}
			}
		}
	}
}

// NewProcessor implements frame.PacketCodec
func (c *JsonPacketCodec) NewProcessor() (clone frame.PacketProcessor) {
	return &JsonPacketCodec{
		rts: c.rts,
		buf: bytes.NewBuffer(nil),
	}
}

// Serializer implements frame.PacketCodec
func (c *JsonPacketCodec) Serializer() serialize.Serializer {
	return c.ser
}

// Decode implements frame.PacketProcessor.
func (c *JsonPacketCodec) Decode(data []byte) (msgs []*frame.Message, err error) {
	c.buf.Write(data)
	// JSON stream decoder
	stream := jsoniter.NewDecoder(c.buf)

	for stream.More() {
		// Decode JSON object
		stru := Message{}
		if err = stream.Decode(&stru); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		target, found := c.rts[stru.PackType]
		if !found {
			fmt.Printf("packet %d not supported", stru.PackType)
			continue
		}

		msgs = append(msgs, &frame.Message{
			Type:  frame.Request,
			ID:    toMessageID(stru.SessionID, stru.EsRqTime, target.rsID),
			Route: target.handler,
			Data:  []byte(stru.Data),
		})
	}

	return
}

// Encode implements frame.PacketProcessor.
func (c *JsonPacketCodec) Encode(msg *frame.Message) (data []byte, err error) {
	stru := Message{}
	stru.SessionID, stru.EsRqTime, stru.PackType = splitSessionAndPackType(msg.ID)
	stru.EsRsTime = uint64(time.Now().UnixMilli())
	stru.Data = string(msg.Data)
	// Encode
	return json.Marshal(stru)
}
