// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file may have been modified by CloudWeGo authors. All CloudWeGo
// Modifications are Copyright 2022 CloudWeGo Authors.

package websocket_service

import (
	"context"
	"time"

	"github.com/hertz-contrib/websocket"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"go.uber.org/zap"
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readJsonPump() {
	defer fkalert.RecoverAlertException()
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

		if messageType != websocket.TextMessage {
			continue
		}

		if c.isClose.Load() {
			c.InfoWF("client is close", zap.Any("sessionId", c.sessionId))
			continue
		}

		c.processJsonPacket(message)
	}
}

func (c *Client) processJsonPacket(data []byte) {
	c.DebugWF("process json packet", zap.Any("data", len(data)))
	// c.send <- data

	gDefaultTCPPkgCtl.ProcJson(fknet.NewTCPContext(context.TODO(), c), data)
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeJsonPump() {
	defer fkalert.RecoverAlertException()
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
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
