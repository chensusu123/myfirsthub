package websocket_service_actor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"go.uber.org/zap"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type WebsocketServer struct {
	addr   string
	server bool
	pf     fknet.FkProtocolFactory
	wg     sync.WaitGroup
	fklog.FKLogI
	h           server.Hertz
	isRun       bool
	actorSystem *actor.ActorSystem
	connsMgrPID *actor.PID
}

func (t *WebsocketServer) GetActorSystem() *actor.ActorSystem {
	return t.actorSystem
}

func NewWebsocketServer() *WebsocketServer {
	actorSystem := actor.NewActorSystem(actor.WithLoggerFactory(zapAdapterLogging))
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewConnActorMgr(actorSystem)
	})
	connsMgrPID := actorSystem.Root.Spawn(props)
	return &WebsocketServer{server: false, actorSystem: actorSystem, connsMgrPID: connsMgrPID}
}

func (t *WebsocketServer) StopServer() {
	t.server = false
	t.wg.Wait()
}

func (t *WebsocketServer) IsOpening() bool {
	return t.server
}

func (t *WebsocketServer) GetAddr() string {
	return t.addr
}

func (t *WebsocketServer) OpenServer(addr string, pf fknet.FkProtocolFactory) error {
	if len(t.addr) != 0 {
		return errors.New("server addr has beed setted")
	}

	t.addr = addr
	t.FKLogI = fklog.AppLogger().Clone(fmt.Sprintf("server-addr:%s", addr))

	t.server = true
	t.pf = pf
	// t.h = h
	// go t.beginServer()
	return nil
}

func (t *WebsocketServer) beginServer() {
	for t.server {
		t.InfoWF("WebsocketServer begin server", zap.String("addr", t.addr))
		t.h.Spin()
		break
	}
}

func (t *WebsocketServer) Init(addr string, pf fknet.FkProtocolFactory) error {
	if len(t.addr) != 0 {
		return errors.New("server addr has beed setted")
	}

	t.addr = addr
	t.FKLogI = fklog.AppLogger().Clone(fmt.Sprintf("server-addr:%s", addr))
	// system := actor.NewActorSystem(actor.WithLoggerFactory(zapAdapterLogging))
	// t.actorSystem = system

	actorSystem := actor.NewActorSystem(actor.WithLoggerFactory(zapAdapterLogging))
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewConnActorMgr(actorSystem)
	})
	connsMgrPID := actorSystem.Root.Spawn(props)
	t.actorSystem = actorSystem
	t.connsMgrPID = connsMgrPID

	h := server.Default(server.WithHostPorts(addr))

	h.GET("/pb", func(c context.Context, ctx *app.RequestContext) {
		serveActorWs(ctx, t.FKLogI, t.actorSystem, t.connsMgrPID, false)
	})

	h.GET("/json", func(c context.Context, ctx *app.RequestContext) {
		serveActorWs(ctx, t.FKLogI, t.actorSystem, t.connsMgrPID, true)
	})

	t.server = true
	t.pf = pf
	t.h = *h
	return nil
}

func (t *WebsocketServer) Start() error {
	t.InfoWF("WebsocketServer Start begin server", zap.String("addr", t.addr))
	go t.beginServer()
	return nil
}

func GetActorSystem() *actor.ActorSystem {
	return gGlobalTCPRawServer.actorSystem
}

func GetConnsMgrPID() *actor.PID {
	return gGlobalTCPRawServer.connsMgrPID
}
