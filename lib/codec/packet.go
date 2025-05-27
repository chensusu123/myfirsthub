package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"maze_game_server/lib/codec/raw_pkg"

	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/frame"
	"maze_game_server/lib/nano/serialize"
)

// TODO 需要补全日志
type EsPacketCodec struct {
	rts  map[uint16]target
	ser  serialize.Serializer
	buf  *bytes.Buffer
	size int // last packet length
}

func NewEsPacketCodec(routes *Routes, opts ...CodecOption) (c *EsPacketCodec) {
	c = &EsPacketCodec{
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
func (c *EsPacketCodec) initRoutes(routes *Routes) {
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
func (c *EsPacketCodec) NewProcessor() (clone frame.PacketProcessor) {
	return &EsPacketCodec{
		rts:  c.rts,
		buf:  bytes.NewBuffer(nil),
		size: -1,
	}
}

// Serializer implements frame.PacketCodec
func (c *EsPacketCodec) Serializer() serialize.Serializer {
	return c.ser
}

// Decode implements frame.PacketProcessor.
func (c *EsPacketCodec) Decode(data []byte) (msgs []*frame.Message, err error) {
	c.buf.Write(data)

	twoBytes := [2]byte{}

	for {
		if c.buf.Len() < 2 {
			break
		}

		// Read header
		_, err = c.buf.Read(twoBytes[:])
		if err != nil {
			return
		}

		packetLen := int(binary.LittleEndian.Uint16(twoBytes[:]))
		if packetLen > math.MaxUint16 || packetLen < 4 {
			err = errors.New("invalid packet len err")
			return
		}

		if c.buf.Len() < packetLen-2 {
			break
		}

		var packet = make([]byte, packetLen)
		copy(packet, twoBytes[:])
		copy(packet[2:], c.buf.Next(packetLen-2))

		stru := raw_pkg.StruSvrEsRawBaseHead{}
		// Decode
		err = stru.UnPack(packet)
		if err != nil {
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
			Data:  stru.Data,
		})
	}

	return
}

// Encode implements frame.PacketProcessor.
func (c *EsPacketCodec) Encode(msg *frame.Message) (data []byte, err error) {
	stru := raw_pkg.StruSvrEsRawBaseHead{}
	stru.SessionID, stru.EsRqTime, stru.PackType = splitSessionAndPackType(msg.ID)
	stru.EsRsTime = uint64(time.Now().UnixMilli())
	stru.Data = msg.Data
	stru.SetTeaflag()
	// Encode
	return stru.Pack()
}

func ToMessageID(sessionID uint32, rqTime uint64, rsID uint16) (messageID uint64) {
	return toMessageID(sessionID, rqTime, rsID)
}

// toMessageID 将ES协议头的部分字段转换成MessageID使用
// |63<---- session id ---->32|
// |31<---- RqTime的末尾16个位(如果响应超过1分钟则会计算失真) ---->16|
// |15<---- RsID ---->0|
func toMessageID(sessionID uint32, rqTime uint64, rsID uint16) (messageID uint64) {
	return uint64(sessionID)<<32 | rqTime&0xFFFF<<16 | uint64(rsID)
}

func splitSessionAndPackType(messageID uint64) (sessionID uint32, rqTime uint64, rsID uint16) {
	if messageID > 0 {
		sessionID = uint32(messageID & 0xFFFFFFFF00000000 >> 32)
		rqTime = uint64(time.Now().UnixMilli())&0xFFFFFFFFFFFF0000 | uint64(messageID&0x00000000FFFF0000>>16)
		rsID = uint16(messageID & 0xFFFF)
	}
	return
}
