// Copyright (c) nano Authors. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"maze_game_server/lib/nano/frame"
	"maze_game_server/lib/nano/internal/codec"
	"maze_game_server/lib/nano/internal/env"
	"maze_game_server/lib/nano/internal/log"
	"maze_game_server/lib/nano/internal/message"
	"maze_game_server/lib/nano/internal/packet"
	"maze_game_server/lib/nano/pipeline"
	"maze_game_server/lib/nano/scheduler"
	"maze_game_server/lib/nano/serialize"
	"maze_game_server/lib/nano/session"

	packCodec "maze_game_server/lib/codec"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	agentWriteBacklog = 16
)

var (
	// ErrBrokenPipe represents the low-level connection has broken.
	ErrBrokenPipe = errors.New("broken low-level pipe")
	// ErrBufferExceed indicates that the current session buffer is full and
	// can not receive more data.
	ErrBufferExceed = errors.New("session send buffer exceed")
)

type (
	// Agent corresponding a user, used for store raw conn information
	agent struct {
		// regular agent member
		session    *session.Session    // session
		conn       net.Conn            // low-level conn fd
		lastMid    uint64              // last message id
		state      int32               // current agent state
		chDie      chan struct{}       // wait for close
		chSend     chan pendingMessage // push message queue
		lastAt     int64               // last heartbeat unix time stamp
		decoder    *codec.Decoder      // binary decoder
		pcodec     frame.PacketProcessor
		serializer serialize.Serializer
		pipeline   pipeline.Pipeline

		rpcHandler rpcHandler
		srv        reflect.Value // cached session reflect.Value
	}

	pendingMessage struct {
		uid     int64
		ctx     context.Context
		typ     message.Type // message type
		route   string       // message route(push)
		mid     uint64       // response message id(response)
		payload interface{}  // payload
	}
)

// Create new agent instance
func newAgent(conn net.Conn, pipeline pipeline.Pipeline, pcodec frame.PacketCodec, rpcHandler rpcHandler) *agent {
	a := &agent{
		conn:       conn,
		state:      statusStart,
		chDie:      make(chan struct{}),
		lastAt:     time.Now().Unix(),
		chSend:     make(chan pendingMessage, agentWriteBacklog),
		decoder:    codec.NewDecoder(),
		pcodec:     pcodec.NewProcessor(),
		serializer: pcodec.Serializer(),
		pipeline:   pipeline,
		rpcHandler: rpcHandler,
	}

	// binding session
	s := session.New(a)
	a.session = s
	a.srv = reflect.ValueOf(s)

	return a
}

func (a *agent) send(m pendingMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrBrokenPipe
		}
	}()
	a.chSend <- m
	return
}

// LastMid implements the session.NetworkEntity interface
func (a *agent) LastMid() uint64 {
	return a.lastMid
}

// Push, implementation for session.NetworkEntity interface
func (a *agent) Push(ctx context.Context, route string, v interface{}) error {
	userID := a.session.UID()
	agentSession := a.session.ID()
	ctx, span := agentSendSpan(ctx, "push", 0, userID, agentSession)
	defer span.End()
	if a.status() == statusClosed {
		span.RecordError(ErrBrokenPipe)
		span.SetStatus(codes.Error, "agent is closed")
		return ErrBrokenPipe
	}

	if len(a.chSend) >= agentWriteBacklog {
		span.RecordError(ErrBufferExceed)
		span.SetStatus(codes.Error, "ErrBufferExceed")
		return ErrBufferExceed
	}

	if env.Debug {
		switch d := v.(type) {
		case []byte:
			log.Println(fmt.Sprintf("Type=Push, ID=%d, UID=%d, Route=%s, Data=%dbytes",
				a.session.ID(), a.session.UID(), route, len(d)))
		default:
			log.Println(fmt.Sprintf("Type=Push, ID=%d, UID=%d, Route=%s, Data=%+v",
				a.session.ID(), a.session.UID(), route, v))
		}
	}

	err := a.send(pendingMessage{ctx: ctx, typ: message.Push, route: route, payload: v, uid: userID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// RPC, implementation for session.NetworkEntity interface
func (a *agent) RPC(ctx context.Context, route string, v interface{}) error {
	if a.status() == statusClosed {
		return ErrBrokenPipe
	}

	// TODO: buffer
	data, err := message.Serialize(v, nil)
	if err != nil {
		return err
	}
	msg := &message.Message{
		Type:  message.Notify,
		Route: route,
		Data:  data,
	}
	a.rpcHandler(ctx, a.session, msg, true)
	return nil
}

// Response, implementation for session.NetworkEntity interface
// Response message to session
func (a *agent) Response(ctx context.Context, v interface{}) error {
	return a.ResponseMid(ctx, a.lastMid, v)
}

// ResponseMid, implementation for session.NetworkEntity interface
// Response message to session
func (a *agent) ResponseMid(ctx context.Context, mid uint64, v interface{}) error {
	userID := a.session.UID()
	agentSession := a.session.ID()
	ctx, span := agentSendSpan(ctx, "response", mid, userID, agentSession)
	defer span.End()
	if a.status() == statusClosed {
		span.RecordError(ErrBrokenPipe)
		span.SetStatus(codes.Error, "agent is closed")
		return ErrBrokenPipe
	}

	if mid <= 0 {
		span.RecordError(ErrSessionOnNotify)
		span.SetStatus(codes.Error, "ErrSessionOnNotify")
		return ErrSessionOnNotify
	}

	if len(a.chSend) >= agentWriteBacklog {
		span.RecordError(ErrBufferExceed)
		span.SetStatus(codes.Error, "ErrBufferExceed")
		return ErrBufferExceed
	}

	if env.Debug {
		switch d := v.(type) {
		case []byte:
			log.Println(fmt.Sprintf("Type=Response, ID=%d, UID=%d, MID=%d, Data=%dbytes",
				a.session.ID(), a.session.UID(), mid, len(d)))
		default:
			log.Println(fmt.Sprintf("Type=Response, ID=%d, UID=%d, MID=%d, Data=%+v",
				a.session.ID(), a.session.UID(), mid, v))
		}
	}

	err := a.send(pendingMessage{ctx: ctx, typ: message.Response, mid: mid, payload: v, uid: userID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// Close, implementation for session.NetworkEntity interface
// Close closes the agent, clean inner state and close low-level connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (a *agent) Close() error {
	if a.status() == statusClosed {
		return ErrCloseClosedSession
	}
	a.setStatus(statusClosed)

	if env.Debug {
		log.Println(fmt.Sprintf("Session closed, ID=%d, UID=%d, IP=%s",
			a.session.ID(), a.session.UID(), a.conn.RemoteAddr()))
	}

	// prevent closing closed channel
	select {
	case <-a.chDie:
		// expect
	default:
		close(a.chDie)
		scheduler.PushTask(func() { session.Lifetime.Close(a.session) })
	}

	return a.conn.Close()
}

// RemoteAddr, implementation for session.NetworkEntity interface
// returns the remote network address.
func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

// String, implementation for Stringer interface
func (a *agent) String() string {
	return fmt.Sprintf("Remote=%s, LastTime=%d", a.conn.RemoteAddr().String(), atomic.LoadInt64(&a.lastAt))
}

func (a *agent) status() int32 {
	return atomic.LoadInt32(&a.state)
}

func (a *agent) setStatus(state int32) {
	atomic.StoreInt32(&a.state, state)
}

func (a *agent) write() {
	ticker := time.NewTicker(env.Heartbeat)
	chWrite := make(chan WriteItem, agentWriteBacklog)

	// Logger
	logger := fklog.AppLogger().Clone("nano")
	var lastErr error
	_ = logger
	// clean func
	defer func() {
		ticker.Stop()
		close(a.chSend)
		close(chWrite)
		closeWriteSpan(chWrite)
		closeSendMsgSpan(a.chSend)

		if env.SessionMonitor != nil {
			env.SessionMonitor.OnClose(a.session, lastErr)
		}

		a.Close()
		if env.Debug {
			log.Println(fmt.Sprintf("Session write goroutine exit, SessionID=%d, UID=%d", a.session.ID(), a.session.UID()))
		}
		logger.DebugWF("session write goroutine exit",
			zap.Int64("session_id", a.session.ID()),
			zap.Int64("uid", a.session.UID()))
	}()

	for {
		select {
		case <-ticker.C:
			if a.pcodec == nil {
				deadline := time.Now().Add(-2 * env.Heartbeat).Unix()
				if atomic.LoadInt64(&a.lastAt) < deadline {
					log.Println(fmt.Sprintf("Session heartbeat timeout, LastTime=%d, Deadline=%d", atomic.LoadInt64(&a.lastAt), deadline))
					return
				}
				chWrite <- WriteItem{ctx: context.TODO(), data: hbd}
			}

		case dataWrite := <-chWrite:
			err := processDataWrite(a, &dataWrite)
			if err != nil {
				return
			}
			// data := dataWrite.data
			// ctx := dataWrite.ctx
			// span := trace.SpanFromContext(ctx)

			// span.AddEvent("conn.write")
			// // close agent while low-level conn broken
			// if wCount, err := a.conn.Write(data); err != nil {
			// 	lastErr = err
			// 	log.Println(err.Error())
			// 	logger.ErrorWF("nano write packet failed",
			// 		zap.Int("data_len", len(data)),
			// 		zap.String("remote_addr", a.conn.RemoteAddr().String()),
			// 		zap.Int("write_count", wCount), zap.Error(err))
			// 	span.RecordError(err)
			// 	span.SetStatus(codes.Error, "conn.Write failed.")
			// 	span.End()
			// 	return
			// } else {
			// 	logger.InfoWF("nano write packet",
			// 		zap.Int("data_len", len(data)),
			// 		zap.String("remote_addr", a.conn.RemoteAddr().String()),
			// 		zap.Int("write_count", wCount))
			// }
			// span.End()
		case data := <-a.chSend:
			err := processPendingMessage(a, data, chWrite)
			if err != nil {
				lastErr = err
				break
			}
			_ = err
			// ctx, span := sendMsgSpan(&data)
			// span.AddEvent("send.serialize")
			// payload, err := message.Serialize(data.payload, a.serializer)
			// if err != nil {
			// 	lastErr = err
			// 	span.RecordError(err)
			// 	span.SetStatus(codes.Error, "message.Serialize failed.")
			// 	switch data.typ {
			// 	case message.Push:
			// 		log.Println(fmt.Sprintf("Push: %s error: %s", data.route, err.Error()))
			// 	case message.Response:
			// 		log.Println(fmt.Sprintf("Response message(id: %d) error: %s", data.mid, err.Error()))
			// 	default:
			// 		// expect
			// 	}
			// 	span.End()
			// 	break
			// }

			// // construct message and encode
			// m := &message.Message{
			// 	Type:  data.typ,
			// 	Data:  payload,
			// 	Route: data.route,
			// 	ID:    data.mid,
			// }
			// if pipe := a.pipeline; pipe != nil {
			// 	err := pipe.Outbound().Process(a.session, m)
			// 	if err != nil {
			// 		lastErr = err
			// 		log.Println("broken pipeline", err.Error())
			// 		span.End()
			// 		break
			// 	}
			// }

			// var p []byte

			// if a.pcodec != nil {
			// 	callLogger := fklog.ContextAppLogger(ctx)
			// 	callLogger.CtxInfo(ctx, "nano process packet stop",
			// 		zap.Uint64("ID", m.ID),
			// 		zap.String("route", m.Route),
			// 		zap.String("remote_addr", a.conn.RemoteAddr().String()),
			// 		zap.Int("rs_data_len", len(m.Data)))
			// 	span.AddEvent("send.encode")
			// 	p, err = a.pcodec.Encode(m)
			// 	if err != nil {
			// 		span.RecordError(err)
			// 		span.SetStatus(codes.Error, "pcodec.Encode failed.")
			// 		lastErr = err
			// 		log.Println(err.Error())
			// 		span.End()
			// 		break
			// 	}
			// } else {
			// 	span.AddEvent("send.other.encode")
			// 	em, err := m.Encode()
			// 	if err != nil {
			// 		lastErr = err
			// 		log.Println(err.Error())
			// 		span.End()
			// 		break
			// 	}

			// 	// packet encode
			// 	span.AddEvent("send.pack.encode")
			// 	p, err = codec.Encode(packet.Data, em)
			// 	if err != nil {
			// 		lastErr = err
			// 		log.Println(err)
			// 		span.End()
			// 		break
			// 	}
			// }
			// // span.End()
			// span.AddEvent("send.to.chWrite")
			// chWrite <- WriteItem{ctx: ctx, data: p}

		case <-a.chDie: // agent closed signal
			return

		case <-env.Die: // application quit
			return
		}
	}
}

type WriteItem struct {
	ctx  context.Context
	data []byte
}

func sendMsgSpan(data *pendingMessage) (context.Context, trace.Span) {
	ctx := data.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	tracer := otel.Tracer("nano.send.message")
	sessionID, rqTime, rsID := packCodec.SplitSessionAndPackType(data.mid)
	ctx, span := tracer.Start(ctx, "agent.send.msg")
	span.SetAttributes(attribute.String("route", data.route),
		attribute.String("message.type", data.typ.String()),
		attribute.Int64("agent.session", int64(sessionID)),
		attribute.Int64("rsID", int64(rsID)),
		attribute.Int64("rqTime", int64(rqTime)),
		attribute.Int64("enduser.id", data.uid),
	)
	span.AddEvent("send.init")
	return ctx, span
}

func processPendingMessage(a *agent, data pendingMessage, chWrite chan WriteItem) (err error) {
	ctx, span := sendMsgSpan(&data)
	span.AddEvent("send.serialize")
	payload, err := message.Serialize(data.payload, a.serializer)
	allOK := false
	defer func() {
		if !allOK {
			callLogger := fklog.ContextAppLogger(ctx)
			callLogger.CtxError(ctx, "processPendingMessage failed",
				zap.Uint64("ID", data.mid),
				zap.String("route", data.route),
				zap.String("remote_addr", a.conn.RemoteAddr().String()),
				zap.Int("rs_data_len", len(payload)), zap.Error(err))
			span.End()
		}
	}()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "message.Serialize failed.")
		switch data.typ {
		case message.Push:
			log.Println(fmt.Sprintf("Push: %s error: %s", data.route, err.Error()))
		case message.Response:
			log.Println(fmt.Sprintf("Response message(id: %d) error: %s", data.mid, err.Error()))
		default:
			// expect
		}
		return err
	}

	// construct message and encode
	m := &message.Message{
		Type:  data.typ,
		Data:  payload,
		Route: data.route,
		ID:    data.mid,
	}

	if pipe := a.pipeline; pipe != nil {
		span.AddEvent("send.pipeline")
		err = pipe.Outbound().Process(a.session, m)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Outbound.process failed.")
			return err
		}
	}

	var p []byte

	if a.pcodec != nil {
		callLogger := fklog.ContextAppLogger(ctx)
		callLogger.CtxInfo(ctx, "nano process packet stop",
			zap.Uint64("ID", m.ID),
			zap.String("route", m.Route),
			zap.String("remote_addr", a.conn.RemoteAddr().String()),
			zap.Int("rs_data_len", len(m.Data)),
		)
		span.AddEvent("send.encode")
		p, err = a.pcodec.Encode(m)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "pcodec.Encode failed.")
			return err
		}
	} else {
		span.AddEvent("send.other.encode")
		em, err := m.Encode()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "other.msg.Encode failed.")
			return err
		}

		// packet encode
		span.AddEvent("send.pack.encode")
		p, err = codec.Encode(packet.Data, em)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "other.packet.Encode failed.")
			return err
		}
	}
	// span.End()
	span.AddEvent("send.to.chWrite")
	allOK = true
	chWrite <- WriteItem{ctx: ctx, data: p}
	return nil
}

func processDataWrite(a *agent, dataWrite *WriteItem) (err error) {
	data := dataWrite.data
	ctx := dataWrite.ctx
	span := trace.SpanFromContext(ctx)

	defer span.End()
	span.AddEvent("conn.write")
	// close agent while low-level conn broken
	if wCount, err := a.conn.Write(data); err != nil {
		fklog.ContextAppLogger(ctx).CtxError(ctx, "nano write packet failed",
			zap.Int("data_len", len(data)),
			zap.String("remote_addr", a.conn.RemoteAddr().String()),
			zap.Int("write_count", wCount), zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "conn.Write failed.")
		return err
	}
	return nil
}

func closeWriteSpan(chWrite chan WriteItem) {
	for v := range chWrite {
		ctx := v.ctx
		span := trace.SpanFromContext(ctx)
		span.End()
	}
}

func closeSendMsgSpan(chPendingMessage chan pendingMessage) {
	for v := range chPendingMessage {
		ctx := v.ctx
		span := trace.SpanFromContext(ctx)
		span.End()
	}
}
