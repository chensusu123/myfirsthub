// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package websocket_service_actor

import (
	"context"
	"encoding/binary"
	"errors"
	"log/slog"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"github.com/hertz-contrib/websocket"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet/fkpkg"
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
	writeActor    *actor.PID
	actorSystem   *actor.ActorSystem
	userId        uint64
	isJson        bool
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

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
		token := ctx.Query("jwt")

		if token == "" {
			return false
		}

		// 实现JWT验证逻辑
		return validateToken(token)
	},
}

type MyCustomClaims struct {
	jwt.StandardClaims
}

var secretKey = []byte("zU6W/(Y%,KX?-@q4m~tLy1_uhcekTNQg")

func VerifyWithCustomClaims(tokenString string) (*MyCustomClaims, error) {
	logger := fklog.AppLogger().Clone("VerifyWithCustomClaims")
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		logger.WarnWF("VerifyWithCustomClaims Error parsing token", zap.String("token", tokenString), zap.Error(err))
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		logger.WarnWF("VerifyWithCustomClaims Token is valid", zap.String("token", tokenString), zap.Error(err))
		return claims, nil
	}

	logger.InfoWF("VerifyWithCustomClaims Token is invalid", zap.String("token", tokenString), zap.Error(err))
	return nil, err
}

func validateToken(tokenIn string) bool {
	_, err := VerifyWithCustomClaims(tokenIn)
	return err == nil
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
	c.actorSystem.Root.Send(c.writeActor, data)
	return nil
}

// Send 发送封装的消息体
func (c *Client) Send(msg fkpkg.PkgWriter) error {
	data, err := msg.Pack()
	if err != nil {
		return err
	}
	c.actorSystem.Root.Send(c.writeActor, data)
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
	return c.userId
}

func (c *Client) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case []byte:
		ctx.Logger().Info("Received message", slog.Int("messageLen", len(msg)), slog.Int64("userID", int64(c.userId)),
			slog.Uint64("sessionId", c.sessionId))
		if c.isJson {
			c.actorProcessJsonPacket(ctx, msg)
		} else {
			messageLen := len(msg)
			processLen := 0
			for processLen < messageLen {
				n, err := c.ActorProcsssBytes(ctx, msg[processLen:])
				if err != nil {
					c.ErrorWF("process bytes error", zap.Error(err))
					return
				}
				processLen += n
			}
		}

	case *ClientLogin:
		ctx.Send(GetConnsMgrPID(), &ClientKick{UserId: msg.UserId, SessionId: msg.SessionId, wPID: msg.wPID})
		ctx.Logger().Info("client ClientLogin  message", slog.Uint64("UserId", msg.UserId), slog.Uint64("SessionId", msg.SessionId))
		msg.wPID = c.writeActor
		msg.rPID = ctx.Self()
		c.userId = msg.UserId
		ctx.Send(GetConnsMgrPID(), msg)
	}
}

func (c *Client) ActorProcsssBytes(ctx actor.Context, data []byte) (int, error) {
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
				c.actorProcessPacket(ctx, c.readBuff[0:c.packetLen])
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
			c.actorProcessPacket(ctx, data[0:packetLen])
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

func (c *Client) actorProcessPacket(actorCtx actor.Context, data []byte) {
	c.DebugWF("process packet", zap.Any("data", len(data)))
	gDefaultTCPPkgCtl.ProcActor(fknet.NewTCPContext(context.TODO(), c), actorCtx, data)
}

func (c *Client) actorProcessJsonPacket(actorCtx actor.Context, data []byte) {
	c.DebugWF("process json packet", zap.Any("data", len(data)))
	gDefaultTCPPkgCtl.ProcActorJson(fknet.NewTCPContext(context.TODO(), c), actorCtx, data)
}
