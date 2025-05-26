// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package websocket_service_raw

import (
	"sync"
	"sync/atomic"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"go.uber.org/zap"
	"maze_game_server/lib/net/raw_pkg"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients     map[*Client]struct{}
	clientsLock sync.RWMutex

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	login      chan *ClientLogin
	clientMaps map[uint64]*Client
	fklog.FKLogI

	sendDataQueue chan *SendDataMsg

	sessionId atomic.Uint64
	isRun     bool
}

type ClientLogin struct {
	UserId    uint64
	SessionId uint64
	Client    *Client
}

type SendDataMsg struct {
	UserId uint64
	Type   uint16
	DataV2 interface{}
}

var hub = newHub()

func newHub() *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]struct{}),
		clientMaps:    make(map[uint64]*Client),
		login:         make(chan *ClientLogin),
		sendDataQueue: make(chan *SendDataMsg, 1024),
	}
}

func (h *Hub) MakeSession() uint64 {
	session := h.sessionId.Add(1)
	return session
}

func (h *Hub) run(logger fklog.FKLogI) {
	if h.isRun {
		return
	}
	h.isRun = true
	h.FKLogI = logger.Clone("websocket_hub")
	h.sessionId.Store(1)
	for {
		select {
		case client := <-h.register:
			h.Register(client)
		case client := <-h.unregister:
			h.Unregister(client)
		case login := <-h.login:
			h.Login(login)
		case sendData := <-h.sendDataQueue:
			h.SendData(sendData)
		}
	}
}

func (h *Hub) Login(login *ClientLogin) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.InfoWF("Login client ", zap.Any("userID", login.UserId),
		zap.Any("new sessionId", login.SessionId),
	)

	if client, ok := h.clientMaps[login.UserId]; ok {
		h.InfoWF("Login client kick out", zap.Any("userID", login.UserId),
			zap.Any("old sessionId", client.sessionId),
			zap.Any("new sessionId", login.SessionId),
			zap.Any("RemoteAddr", client.conn.RemoteAddr()))
		client.Close()
		delete(h.clientMaps, login.UserId)
	}

	h.clientMaps[login.UserId] = login.Client
}

func (h *Hub) Register(client *Client) {
	h.InfoWF("Register client", zap.Any("RemoteAddr", client.conn.RemoteAddr()), zap.Any("sessionId", client.sessionId))
	h.AddClient(client)
}

func (h *Hub) AddClient(client *Client) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	h.clients[client] = struct{}{}
}

func (h *Hub) Unregister(client *Client) {
	userID := client.GetUserID()
	h.InfoWF("Unregister client", zap.Any("userID", userID), zap.Any("sessionId", client.sessionId))
	h.DelClient(client)
}

func (h *Hub) DelClient(client *Client) {
	h.clientsLock.Lock()
	defer h.clientsLock.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
	userID := client.GetUserID()
	if clientFound, ok := h.clientMaps[userID]; ok {
		if clientFound == client {
			delete(h.clientMaps, userID)
		}
	}
}

func (h *Hub) OnLogin(userID uint64, ctx fknet.TCPContext) {
	client := ctx.GetTag("mySelf").(*Client)
	clientSession := ctx.GetTag("mySelfSession").(uint64)
	h.InfoWF("OnLogin client", zap.Any("userID", userID),
		zap.Any("sessionId", clientSession),
	)
	h.login <- &ClientLogin{
		UserId:    userID,
		SessionId: clientSession,
		Client:    client,
	}
}

func (h *Hub) SendData(info *SendDataMsg) {
	h.InfoWF("SendData entry", zap.Any("userID", info.UserId), zap.Any("packetType", info.Type))
	h.clientsLock.RLock()
	var data []byte
	var err error
	if client, ok := h.clientMaps[info.UserId]; ok {
		h.clientsLock.RUnlock()
		if client.isJosn {
			sendPacket := &NoramlJsonMsg{
				MsgType: int(info.Type),
				Data:    info.DataV2,
			}
			data, err = json.Marshal(sendPacket)
			if err != nil {
				h.ErrorWF("SendData Marshal json failed",
					zap.Any("err", err),
					zap.Any("userID", info.UserId),
					zap.Any("packetType", info.Type))
				return
			}
		} else {
			dataPB, err := proto.Marshal(info.DataV2.(proto.Message))
			pkg := &raw_pkg.StruSvrEsRawBaseHead{}
			pkg.SessionID = uint32(0)
			pkg.PackType = info.Type
			pkg.Data = dataPB
			pkg.EsRsTime = uint64(time.Now().Unix())
			pkg.SetTeaflag()
			data, err = pkg.Pack()
			if err != nil {
				h.ErrorWF("SendData Marshal pb failed",
					zap.Any("err", err),
					zap.Any("userID", info.UserId),
					zap.Any("packetType", info.Type))
				return
			}
		}

		h.InfoWF("SendData", zap.Any("userID", info.UserId),
			zap.Any("sessionId", client.sessionId),
			zap.Any("packetType", info.Type),
		)
		err := client.SendData(data)
		if err != nil {
			h.InfoWF("SendData error", zap.Any("userID", info.UserId),
				zap.Any("sessionId", client.sessionId),
				zap.Any("packetType", info.Type),
				zap.Error(err),
			)
		} else {
			h.InfoWF("SendData success", zap.Any("userID", info.UserId),
				zap.Any("sessionId", client.sessionId),
				zap.Any("packetType", info.Type),
			)
		}
		return
	} else {
		h.InfoWF("SendData client not found", zap.Any("userID", info.UserId), zap.Any("packetType", info.Type))
	}
	h.clientsLock.RUnlock()
}

// func (h *Hub) SendDataByUserID(logger fklog.FKLogI, userID uint64, data []byte) {
// 	h.sendDataQueue <- &SendDataMsg{
// 		UserId: userID,
// 		Data:   data,
// 	}
// }

func (h *Hub) SendDataPacketUserID(logger fklog.FKLogI, userID uint64, packType uint16, pack interface{}) {
	h.sendDataQueue <- &SendDataMsg{
		UserId: userID,
		DataV2: pack,
		Type:   packType,
	}
}
