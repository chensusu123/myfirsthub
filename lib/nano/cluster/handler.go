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
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	packCodec "maze_game_server/lib/codec"
	"maze_game_server/lib/nano/cluster/clusterpb"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/frame"
	"maze_game_server/lib/nano/internal/codec"
	"maze_game_server/lib/nano/internal/env"
	"maze_game_server/lib/nano/internal/log"
	"maze_game_server/lib/nano/internal/message"
	"maze_game_server/lib/nano/internal/packet"
	"maze_game_server/lib/nano/nanometrics"
	"maze_game_server/lib/nano/pipeline"
	"maze_game_server/lib/nano/scheduler"
	"maze_game_server/lib/nano/serialize"
	"maze_game_server/lib/nano/session"

	"github.com/gorilla/websocket"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/logidutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	// cached serialized data
	hrd []byte // handshake response data
	hbd []byte // heartbeat packet data
)

type rpcHandler func(ctx context.Context, session *session.Session, msg *message.Message, noCopy bool)

// CustomerRemoteServiceRoute customer remote service route
type CustomerRemoteServiceRoute func(service string, session *session.Session, members []*clusterpb.MemberInfo) *clusterpb.MemberInfo

func cache() {
	hrdata := map[string]interface{}{
		"code": 200,
		"sys": map[string]interface{}{
			"heartbeat":  env.Heartbeat.Seconds(),
			"servertime": time.Now().UTC().Unix(),
		},
	}
	if dict, ok := message.GetDictionary(); ok {
		hrdata = map[string]interface{}{
			"code": 200,
			"sys": map[string]interface{}{
				"heartbeat":  env.Heartbeat.Seconds(),
				"servertime": time.Now().UTC().Unix(),
				"dict":       dict,
			},
		}
	}
	// data, err := json.Marshal(map[string]interface{}{
	// 	"code": 200,
	// 	"sys": map[string]float64{
	// 		"heartbeat": env.Heartbeat.Seconds(),
	// 	},
	// })
	data, err := json.Marshal(hrdata)
	if err != nil {
		panic(err)
	}

	hrd, err = codec.Encode(packet.Handshake, data)
	if err != nil {
		panic(err)
	}

	hbd, err = codec.Encode(packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}

type LocalHandler struct {
	localServices map[string]*component.Service // all registered service
	localHandlers map[string]*component.Handler // all handler method

	mu             sync.RWMutex
	remoteServices map[string][]*clusterpb.MemberInfo

	pipeline    pipeline.Pipeline
	currentNode *Node
	// Custom packet encoder/decoder
	pcodec    frame.PacketCodec
	taskCount atomic.Int64
	userCount atomic.Int64
}

func NewHandler(currentNode *Node, pipeline pipeline.Pipeline, pcodec frame.PacketCodec) *LocalHandler {
	h := &LocalHandler{
		localServices:  make(map[string]*component.Service),
		localHandlers:  make(map[string]*component.Handler),
		remoteServices: map[string][]*clusterpb.MemberInfo{},
		pipeline:       pipeline,
		currentNode:    currentNode,
		pcodec:         pcodec,
	}

	return h
}

func (h *LocalHandler) register(comp component.Component, opts []component.Option) error {
	s := component.NewService(comp, opts)

	if _, ok := h.localServices[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all localHandlers
	h.localServices[s.Name] = s
	for name, handler := range s.Handlers {
		n := fmt.Sprintf("%s.%s", s.Name, name)
		log.Println("Register local handler", n)
		h.localHandlers[n] = handler
	}
	return nil
}

func (h *LocalHandler) initRemoteService(members []*clusterpb.MemberInfo) {
	for _, m := range members {
		h.addRemoteService(m)
	}
}

func (h *LocalHandler) addRemoteService(member *clusterpb.MemberInfo) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, s := range member.Services {
		log.Println("Register remote service", s)
		h.remoteServices[s] = append(h.remoteServices[s], member)
	}
}

func (h *LocalHandler) delMember(addr string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for name, members := range h.remoteServices {
		for i, maddr := range members {
			if addr == maddr.ServiceAddr {
				if i >= len(members)-1 {
					members = members[:i]
				} else {
					members = append(members[:i], members[i+1:]...)
				}
			}
		}
		if len(members) == 0 {
			delete(h.remoteServices, name)
		} else {
			h.remoteServices[name] = members
		}
	}
}

func (h *LocalHandler) LocalService() []string {
	var result []string
	for service := range h.localServices {
		result = append(result, service)
	}
	sort.Strings(result)
	return result
}

func (h *LocalHandler) RemoteService() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []string
	for service := range h.remoteServices {
		result = append(result, service)
	}
	sort.Strings(result)
	return result
}

func (h *LocalHandler) handle(conn net.Conn, r *http.Request, pcodec frame.PacketCodec) {
	uerCount := h.userCount.Add(1)
	nanometrics.UserCountGauge.Set(float64(uerCount))
	var agentSessionID int64
	var closeNoraml bool

	// Select a packet codec
	if pcodec == nil {
		pcodec = h.pcodec
	}
	handleLogger := fklog.AppLogger().Clone("nano_handle")

	// create a client agent and startup write gorontine
	agent := newAgent(conn, h.pipeline, pcodec, h.remoteProcess)
	agentSessionID = agent.session.ID()

	handleLogger.InfoWF("agent handle entry",
		zap.Int64("agentSessionID", agentSessionID),
		zap.String("remote_addr", agent.conn.RemoteAddr().String()),
	)

	defer func() {
		fkalert.RecoverAlertException()
		uerCount = h.userCount.Add(-1)
		nanometrics.UserCountGauge.Set(float64(uerCount))
		agent.session.SetClose()
		fklog.AppLogger().InfoWF("agent handle end",
			zap.Int64("agentSessionID", agentSessionID),
			zap.Bool("closeNoraml", closeNoraml),
			zap.String("remote_addr", agent.conn.RemoteAddr().String()),
			zap.Int64("enduser.id", agent.session.UID()))
	}()
	// Init session
	// 将Websocket连接请求Header中的数据转存至Session
	if r != nil {
		// Header: X-Forwarded-For
		addr := r.Header.Get("X-Forwarded-For")
		if len(addr) > 0 {
			agent.session.Set("ClientAddr", addr)
		}
	}

	// Logger
	logger := fklog.AppLogger().Clone("nano")

	h.currentNode.storeSession(agent.session)

	if env.SessionMonitor != nil {
		env.SessionMonitor.OnCreate(context.Background(), agent.session)
	}

	var lastErr error

	// startup write goroutine
	go agent.write()

	// guarantee agent related resource be destroyed
	defer func() {
		closelogger := logger.Clone("nano")
		closelogger.SetLogId(logidutil.GenerateLogID())
		closelogger.SetUid(uint64(agent.session.UID()))
		ctx := fklog.ContextWithLogger(context.Background(), closelogger)

		ctx, span := closeHandleSpan(ctx, agent, "read.close")
		defer span.End()
		request := &clusterpb.SessionClosedRequest{
			SessionId: agent.session.ID(),
		}
		members := h.currentNode.cluster.remoteAddrs()
		for _, remote := range members {
			pool, err := h.currentNode.rpcClient.getConnPool(remote)
			if err != nil {
				closelogger.CtxWarn(ctx, "Cannot retrieve connection pool for address", zap.Error(err), zap.String("remote", remote))
				continue
			}
			client := clusterpb.NewMemberClient(pool.Get())
			_, err = client.SessionClosed(ctx, request)
			if err != nil {
				closelogger.CtxWarn(ctx, "Cannot closed session in remote address", zap.Error(err), zap.String("remote", remote))
				continue
			}
		}

		if env.SessionMonitor != nil {
			env.SessionMonitor.OnClose(ctx, agent.session, lastErr)
		}

		agent.Close()
		closelogger.CtxDebug(ctx, "Session read goroutine exit",
			zap.Int64("agent.session", agent.session.ID()),
			zap.Int64("enduser.id", agent.session.UID()),
			zap.Any("lastErr", lastErr),
		)
		closeNoraml = true
	}()

	// read loop
	buf := make([]byte, 2048)
	remoteAddr := agent.conn.RemoteAddr().String()
	sessionID := agent.session.ID()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			lastErr = err
			log.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}

		if agent.pcodec != nil {
			loggerLoop := logger.Clone("nano")
			loggerLoop.SetLogId(logidutil.GenerateLogID())
			loggerLoop.SetUid(uint64(agent.session.UID()))
			ctx := fklog.ContextWithLogger(context.Background(), loggerLoop)

			ctx, span := receiveSpan(ctx, n)

			// Must working
			agent.setStatus(statusWorking)
			span.SetAttributes(attribute.String("remote_addr", remoteAddr),
				attribute.Int("recv_data_len", n),
				attribute.Int64("agent.session", sessionID),
				attribute.Int64("enduser.id", agent.session.UID()),
			)
			span.AddEvent("agent.pcodec.decode")
			msgs, packets, err := agent.pcodec.Decode(buf[:n])
			if err != nil {
				lastErr = err
				log.Println(err.Error())
				span.RecordError(err)
				span.SetStatus(codes.Error, "agent.pcodec.Decode failed.")
				// process message decoded
				for index, m := range msgs {
					loggerLoop.CtxInfo(ctx, "nano process packet start", zap.Uint64("ID", m.ID),
						zap.Int64("agentSessionID", agentSessionID),
						zap.String("remote_addr", agent.conn.RemoteAddr().String()),
						zap.String("route", m.Route),
						zap.Uint16("PackLen", packets[index].PackLen),
						zap.Uint16("PackType", packets[index].PackType),
						zap.Uint32("SessionID", packets[index].SessionID),
						zap.Uint64("RqTime", packets[index].EsRqTime),
						zap.Uint8("CompressType", packets[index].CompressType),
						zap.Int("rq_data_len", len(m.Data)),
					)
					h.processMessage(ctx, agent, m)
				}
				fklog.ContextAppLogger(ctx).CtxError(ctx, "nano message processed", zap.Error(err))
				span.End()
				return
			}

			span.AddEvent("process.message", trace.EventOption(trace.WithAttributes(attribute.Int("msg_count", len(msgs)))))
			// process all message
			for index, m := range msgs {
				loggerLoop.CtxInfo(ctx, "nano process packet start", zap.Uint64("ID", m.ID),
					zap.String("route", m.Route),
					zap.String("remote_addr", agent.conn.RemoteAddr().String()),
					zap.Uint16("PackLen", packets[index].PackLen),
					zap.Uint16("PackType", packets[index].PackType),
					zap.Uint32("SessionID", packets[index].SessionID),
					zap.Uint64("RqTime", packets[index].EsRqTime),
					zap.Uint8("CompressType", packets[index].CompressType),
					zap.Int("rq_data_len", len(m.Data)),
				)
				h.processMessage(ctx, agent, m)
			}

			span.End()
		} else {
			// TODO(warning): decoder use slice for performance, packet data should be copy before next Decode
			packets, err := agent.decoder.Decode(buf[:n])
			if err != nil {
				lastErr = err
				log.Println(err.Error())

				// process packets decoded
				for _, p := range packets {
					if err := h.processPacket(agent, p); err != nil {
						lastErr = err
						log.Println(err.Error())
						return
					}
				}
				return
			}

			// process all packets
			for _, p := range packets {
				if err := h.processPacket(agent, p); err != nil {
					lastErr = err
					log.Println(err.Error())
					return
				}
			}
		}
	}
}

func (h *LocalHandler) processPacket(agent *agent, p *packet.Packet) error {
	switch p.Type {
	case packet.Handshake:
		if err := env.HandshakeValidator(agent.session, p.Data); err != nil {
			return err
		}

		if _, err := agent.conn.Write(hrd); err != nil {
			return err
		}

		agent.setStatus(statusHandshake)
		if env.Debug {
			log.Println(fmt.Sprintf("Session handshake Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		}

	case packet.HandshakeAck:
		agent.setStatus(statusWorking)
		if env.Debug {
			log.Println(fmt.Sprintf("Receive handshake ACK Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		}

	case packet.Data:
		if agent.status() < statusWorking {
			return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, remote=%s",
				agent.conn.RemoteAddr().String())
		}

		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		h.processMessage(context.TODO(), agent, msg)

	case packet.Heartbeat:
		// expected
	}

	agent.lastAt = time.Now().Unix()
	return nil
}

func (h *LocalHandler) findMembers(service string) []*clusterpb.MemberInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.remoteServices[service]
}

func (h *LocalHandler) remoteProcess(ctx context.Context, session *session.Session, msg *message.Message, noCopy bool) {
	index := strings.LastIndex(msg.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("nano/handler: invalid route %s", msg.Route))
		return
	}

	service := msg.Route[:index]
	members := h.findMembers(service)
	if len(members) == 0 {
		log.Println(fmt.Sprintf("nano/handler: %s not found(forgot registered?)", msg.Route))
		return
	}

	// Select a remote service address
	// 1. if exist customer remote service route ,use it, otherwise use default strategy
	// 2. Use the service address directly if the router contains binding item
	// 3. Select a remote service address randomly and bind to router
	var remoteAddr string
	if h.currentNode.Options.RemoteServiceRoute != nil {
		if addr, found := session.Router().Find(service); found {
			remoteAddr = addr
		} else {
			member := h.currentNode.Options.RemoteServiceRoute(service, session, members)
			if member == nil {
				log.Println(fmt.Sprintf("customize remoteServiceRoute handler: %s is not found", msg.Route))
				return
			}
			remoteAddr = member.ServiceAddr
			session.Router().Bind(service, remoteAddr)
		}
	} else {
		if addr, found := session.Router().Find(service); found {
			remoteAddr = addr
		} else {
			remoteAddr = members[rand.Intn(len(members))].ServiceAddr
			session.Router().Bind(service, remoteAddr)
		}
	}
	pool, err := h.currentNode.rpcClient.getConnPool(remoteAddr)
	if err != nil {
		log.Println(err)
		return
	}
	data := msg.Data
	if !noCopy && len(msg.Data) > 0 {
		data = make([]byte, len(msg.Data))
		copy(data, msg.Data)
	}

	// Retrieve gate address and session id
	gateAddr := h.currentNode.ServiceAddr
	sessionId := session.ID()
	switch v := session.NetworkEntity().(type) {
	case *acceptor:
		gateAddr = v.gateAddr
		sessionId = v.sid
	}

	client := clusterpb.NewMemberClient(pool.Get())
	switch msg.Type {
	case message.Request:
		request := &clusterpb.RequestMessage{
			GateAddr:  gateAddr,
			SessionId: sessionId,
			Id:        msg.ID,
			Route:     msg.Route,
			Data:      data,
		}
		_, err = client.HandleRequest(ctx, request)
	case message.Notify:
		request := &clusterpb.NotifyMessage{
			GateAddr:  gateAddr,
			SessionId: sessionId,
			Route:     msg.Route,
			Data:      data,
		}
		_, err = client.HandleNotify(ctx, request)
	}
	if err != nil {
		log.Println(fmt.Sprintf("Process remote message (%d:%s) error: %+v", msg.ID, msg.Route, err))
	}
}

func (h *LocalHandler) processMessage(ctx context.Context, agent *agent, msg *message.Message) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		fkalert.RecoverAlertException()
		cancel()
		if ctx.Err() != nil {
			fklog.ContextAppLogger(ctx).ErrorWF("process message timeout", zap.Error(ctx.Err()))
		}
	}()

	var lastMid uint64
	switch msg.Type {
	case message.Request:
		lastMid = msg.ID
	case message.Notify:
		lastMid = 0
	default:
		log.Println("Invalid message type: " + msg.Type.String())
		return
	}

	sessionID, _, rsID := packCodec.SplitSessionAndPackType(msg.ID)
	tracer := otel.Tracer("nano.process.message")
	ctx, span := tracer.Start(ctx, msg.Route)

	span.SetAttributes(
		attribute.Int64("agent.session", agent.session.ID()),
		attribute.Int64("enduser.id", agent.session.UID()),
		attribute.Int64("packet.session", int64(sessionID)),
		attribute.Int64("packet.id", int64(rsID)),
	)
	// agent.session.SetContext(ctx)
	handler, found := h.localHandlers[msg.Route]
	if !found {
		span.AddEvent("nano.remote.process")
		h.remoteProcess(ctx, agent.session, msg, false)
		span.End()
	} else {
		span.AddEvent("nano.local.process")
		h.localProcess(ctx, handler, lastMid, agent.session, agent.serializer, msg)
	}
}

func (h *LocalHandler) handleWS(conn *websocket.Conn, r *http.Request, pcodec frame.PacketCodec) {
	c, err := newWSConn(conn)
	if err != nil {
		log.Println(err)
		return
	}
	go h.handle(c, r, pcodec)
}

func (h *LocalHandler) localProcess(ctx context.Context, handler *component.Handler, lastMid uint64, session *session.Session, serializer serialize.Serializer, msg *message.Message) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("nano.local.process.begin")
	if pipe := h.pipeline; pipe != nil {
		err := pipe.Inbound().Process(session, msg)
		if err != nil {
			log.Println("Pipeline process failed: " + err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, "Inbound().Process")
			span.End()
			return
		}
	}

	payload := msg.Data
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		if serializer == nil {
			serializer = env.Serializer
		}
		span.AddEvent("nano.unmarshal")
		err := serializer.Unmarshal(payload, data)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Unmarshal failed.")
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, payload))
			span.End()
			return
		}
	}

	if env.Debug {
		log.Println(fmt.Sprintf("UID=%d, Message={%s}, Data=%+v", session.UID(), msg.String(), data))
	}

	args := []reflect.Value{handler.Receiver, reflect.ValueOf(session), reflect.ValueOf(data)}
	task := func() {
		switch v := session.NetworkEntity().(type) {
		case *agent:
			v.lastMid = lastMid
		case *acceptor:
			v.lastMid = lastMid
		}
		// span := trace.SpanFromContext(ctx)
		span.AddEvent("nano.func.call.begin")
		session.SetContext(ctx)
		defer func() {
			if err := recover(); err != nil {
				fklog.ContextAppLogger(ctx).ErrorWF("local process panic", zap.Any("err", err))
			}
			span.AddEvent("nano.local.process.end")
			session.SetContext(context.TODO())
			span.End()
			h.taskCount.Add(-1)
			session.TaskCountDec()
		}()

		if session.IsClose() {
			span.AddEvent("session.isclose")
			return
		}

		result := handler.Method.Func.Call(args)
		span.AddEvent("nano.func.call.end")

		if len(result) > 0 {
			if err := result[0].Interface(); err != nil {
				log.Println(fmt.Sprintf("Service %s error: %+v", msg.Route, err))
				span.RecordError(err.(error))
				span.SetStatus(codes.Error, "handler.Method.Func.Call failed.")
			}
		}
	}

	index := strings.LastIndex(msg.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("nano/handler: invalid route %s", msg.Route))
		span.SetStatus(codes.Error, "nano/handler: invalid route")
		span.End()
		return
	}

	// A message can be dispatch to global thread or a user customized thread
	service := msg.Route[:index]
	if s, found := h.localServices[service]; found && s.SchedName != "" {
		sched := session.Value(s.SchedName)
		if sched == nil {
			log.Println(fmt.Sprintf("nanl/handler: cannot found `schedular.LocalScheduler` by %s", s.SchedName))
			span.SetStatus(codes.Error, "nanl/handler: cannot found schedular.LocalScheduler")
			span.End()
			return
		}

		local, ok := sched.(scheduler.LocalScheduler)
		_ = local
		if !ok {
			log.Println(fmt.Sprintf("nanl/handler: Type %T does not implement the `schedular.LocalScheduler` interface",
				sched))
			span.SetStatus(codes.Error, "nanl/handler: cannot found schedular.LocalScheduler")
			span.End()
			return
		}
		span.AddEvent("nano.schedule.task")
		taskCount := h.taskCount.Add(1)
		// sesstionTaskCount := session.TaskCountInc()
		sesstionTaskCount := session.PushTask(task)
		span.SetAttributes(attribute.Int64("nano.current.task.count", taskCount))
		span.SetAttributes(attribute.Int64("nano.session.task.count", sesstionTaskCount))
		span.SetAttributes(attribute.String("nano.task.scheduler.name", service))
		// local.Schedule(task)
	} else {
		span.AddEvent("nano.schedule.task")
		taskCount := h.taskCount.Add(1)
		sesstionTaskCount := session.PushTask(task)
		span.SetAttributes(attribute.Int64("nano.current.task.count", taskCount))
		span.SetAttributes(attribute.Int64("nano.session.task.count", sesstionTaskCount))
		span.SetAttributes(attribute.String("nano.task.scheduler.name", "global"))
		// scheduler.PushTask(task)
		// session.PushTask(task)
		// session.session.scheduler.PushTask(task)
	}
}
