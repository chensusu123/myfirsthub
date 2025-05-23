// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package websocket_service_actor

import (
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
)

type ClientLogin struct {
	UserId    uint64
	SessionId uint64
	Client    *Client

	wPID *actor.PID
	rPID *actor.PID
}

type ClientLogout struct {
	UserId    uint64
	SessionId uint64
	wPID      *actor.PID
	rPID      *actor.PID
}

type ClientKick struct {
	UserId    uint64
	SessionId uint64
	wPID      *actor.PID
}

type SendDataMsg struct {
	UserId     uint64
	Data       interface{}
	PacketType uint16
}
type Hub struct {
	sessionId atomic.Uint64
}

func (h *Hub) MakeSession() uint64 {
	session := h.sessionId.Add(1)
	return session
}

var hub = newHub()

func newHub() *Hub {
	return &Hub{}
}
