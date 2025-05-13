// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package websocket_service

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet/fkpkg"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 64 * 1024
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
	fklog.FKLogI
	sessionId uint64

	readBuff      []byte
	readLen       int
	readOff       int
	readBuffInit  bool
	waitPacketEnd bool
	readPacketLen int
	packetLen     int
	isClose       atomic.Bool
	*fknet.FkTags
}

func (c *Client) resetReadState() {
	c.waitPacketEnd = false
	c.readPacketLen = 0
	c.packetLen = 0
}

const (
	PackageBytesLen = 2
)

func packetLenIsValid(packetLen int) bool {
	return packetLen > math.MaxUint16 || packetLen < 4
}

func (c *Client) ProcsssBytes(data []byte) (int, error) {
	if !c.readBuffInit {
		c.readBuffInit = true
		c.readBuff = make([]byte, 65535)
	}

	messageLen := len(data)
	if c.waitPacketEnd {
		// 等待包结束
		if c.readPacketLen >= PackageBytesLen {
			// 读取到包长度
			if c.readPacketLen+messageLen < c.packetLen {
				// 包还没结束
				_ = copy(c.readBuff[c.readPacketLen:], data[0:])
				c.readPacketLen += messageLen
				c.DebugWF("wait packet end", zap.Any("packetLen", c.packetLen),
					zap.Any("readPacketLen", c.readPacketLen))
				return messageLen, nil
			} else {
				// 包结束
				readLen := c.packetLen - c.readPacketLen
				_ = copy(c.readBuff[c.readPacketLen:], data[0:readLen])
				c.processPacket(c.readBuff[0:c.packetLen])
				c.resetReadState()
				c.DebugWF("process packet", zap.Any("packetLen", c.packetLen))
				return readLen, nil
			}
		} else {
			readLen := PackageBytesLen - c.readPacketLen
			_ = copy(c.readBuff[c.readPacketLen:], data[0:readLen])
			c.readPacketLen += readLen
			if c.readPacketLen < PackageBytesLen {
				return readLen, nil
			}
			packetLen := int(binary.LittleEndian.Uint16(c.readBuff[0:PackageBytesLen]))
			if packetLenIsValid(packetLen) {
				c.ErrorWF("invalid packet len", zap.Any("packetLen", packetLen))
				return messageLen, errors.New("invalid packet len err")
			}
			c.packetLen = packetLen
			c.waitPacketEnd = true
			c.DebugWF("wait packet end", zap.Any("packetLen", c.packetLen),
				zap.Any("readPacketLen", c.readPacketLen))
			return readLen, nil
		}
	} else {
		// 读取包长度
		if messageLen < PackageBytesLen {
			// 包还没结束
			_ = copy(c.readBuff[c.readPacketLen:], data[0:])
			c.readPacketLen += messageLen
			c.DebugWF("wait packet end", zap.Any("packetLen", c.packetLen),
				zap.Any("readPacketLen", c.readPacketLen))
			return messageLen, nil
		}
		packetLen := int(binary.LittleEndian.Uint16(data[0:PackageBytesLen]))
		if packetLenIsValid(packetLen) {
			c.ErrorWF("invalid packet len", zap.Any("packetLen", packetLen))
			return messageLen, errors.New("invalid packet len err")
		}
		if packetLen <= messageLen {
			// 包结束
			c.processPacket(data[0:packetLen])
			c.resetReadState()
			c.DebugWF("process packet", zap.Any("packetLen", packetLen))
			return packetLen, nil
		}
		// 包还没结束
		c.packetLen = packetLen
		_ = copy(c.readBuff[c.readPacketLen:], data[0:])
		c.readPacketLen += messageLen
		c.waitPacketEnd = true
		c.DebugWF("wait packet end", zap.Any("packetLen", c.packetLen),
			zap.Any("readPacketLen", c.readPacketLen))
		return messageLen, nil
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.ErrorWF("read message error", zap.Error(err))
			}
			break
		}

		c.InfoWF("read message",
			zap.Any("message", len(message)),
			zap.Any("messageType", messageType),
		)

		if messageType != websocket.BinaryMessage {
			continue
		}

		if c.isClose.Load() {
			c.InfoWF("client is close", zap.Any("sessionId", c.sessionId))
			continue
		}

		messageLen := len(message)
		processLen := 0
		for processLen < messageLen {
			n, err := c.ProcsssBytes(message[processLen:])
			if err != nil {
				c.ErrorWF("process bytes error", zap.Error(err))
				return
			}
			processLen += n
		}

		// c.send <- message
	}
}

func (c *Client) processPacket(data []byte) {
	c.DebugWF("process packet", zap.Any("data", len(data)))
	// c.send <- data

	gDefaultTCPPkgCtl.Proc(fknet.NewTCPContext(context.TODO(), c), data)
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// if c.conn == nil {
			// 	c.InfoWF("client is close", zap.Any("sessionId", c.sessionId))
			// 	return
			// }
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				c.ErrorWF("write message NextWriter error", zap.Error(err))
				return
			}

			w.Write(message)
			c.InfoWF("write message", zap.Any("message", len(message)))
			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(<-c.send)
			// 	// c.InfoWF("write message loop")
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

// serveWs handles websocket requests from the peer.
func serveWs(ctx *app.RequestContext, logger fklog.FKLogI) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		clientLogger := logger.Clone(fmt.Sprintf("client-addr:%s", conn.RemoteAddr().String()))
		clientLogger.InfoWF("client connected", zap.Any("addr", conn.RemoteAddr().String()))
		client := &Client{conn: conn, send: make(chan []byte, 1024), FkTags: fknet.NewFkTags(), sessionId: hub.MakeSession()}
		client.SetTag("mySelf", client)
		client.SetTag("mySelfSession", client.sessionId)
		client.FKLogI = clientLogger
		hub.register <- client

		go client.writePump()
		client.readPump()
	})
	if err != nil {
		logger.ErrorWF("upgrade error", zap.Error(err))
	}
}

func serveJsonWs(ctx *app.RequestContext, logger fklog.FKLogI) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		clientLogger := logger.Clone(fmt.Sprintf("client-addr:%s", conn.RemoteAddr().String()))
		clientLogger.InfoWF("serveJsonWs client connected", zap.Any("addr", conn.RemoteAddr().String()))
		client := &Client{conn: conn, send: make(chan []byte, 1024), FkTags: fknet.NewFkTags(), sessionId: hub.MakeSession()}
		client.SetTag("mySelf", client)
		client.SetTag("mySelfSession", client.sessionId)
		client.FKLogI = clientLogger
		hub.register <- client

		go client.writeJsonPump()
		client.readJsonPump()
	})
	if err != nil {
		logger.ErrorWF("serveJsonWs upgrade error", zap.Error(err))
	}
}

func (c *Client) MockLogger(logger fklog.FKLogI) {
	c.FKLogI = logger.Clone("websocket_client")
}

// ReadData 从连接中读取数据（需结合现有读逻辑）
func (c *Client) ReadData() ([]byte, error) {
	// 注意：当前 readPump 已经在循环读取消息，实际需根据业务调整数据传递方式
	// 示例逻辑：从读取缓冲区获取数据（需根据实际协议调整）
	// 这里简化为返回空，具体实现需结合现有 ProcsssBytes/processPacket 逻辑
	return nil, nil
}

// SendData 发送字节数据到对端
func (c *Client) SendData(data []byte) error {
	select {
	case c.send <- data:
		return nil
	default:
		return errors.New("send channel is full")
	}
}

// Send 发送封装的消息体
func (c *Client) Send(msg fkpkg.PkgWriter) error {
	data, err := msg.Pack()
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}

// SendStripData 发送剥离后的数据（假设与 SendData 逻辑相同）
func (c *Client) SendStripData(data []byte) error {
	return c.SendData(data)
}

// SendRemain 发送剩余未发送的数据
func (c *Client) SendRemain(data []byte) error {
	return c.SendData(data)
}

// GetRemain 获取剩余数据通道
func (c *Client) GetRemain() <-chan []byte {
	return nil
}

// IsClosed 检查连接是否关闭
func (c *Client) IsClosed() bool {
	return c.isClose.Load()
}

// RemoteAddr 获取对端地址
func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// LocalAddr 获取本地地址
func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// Close 关闭连接
func (c *Client) Close() {
	c.DebugWF("client close", zap.Any("session", c.sessionId))
	if c.isClose.Load() {
		return
	}
	c.isClose.Store(true)
	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	}
}

func (c *Client) GetUserID() uint64 {
	foundUserID := uint64(0)
	tmpUserID := c.GetTag("userID")
	switch v := tmpUserID.(type) {
	case int64:
		foundUserID = uint64(v)
	case uint64:
		foundUserID = v
	}
	return foundUserID
}
