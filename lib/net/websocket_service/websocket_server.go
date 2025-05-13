package websocket_service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"go.uber.org/zap"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type WebsocketServer struct {
	addr   string
	server bool
	pf     fknet.FkProtocolFactory
	wg     sync.WaitGroup
	fklog.FKLogI
	h     server.Hertz
	isRun bool
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{server: false}
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

	h := server.Default(server.WithHostPorts(addr))

	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		serveWs(ctx, t.FKLogI)
	})

	t.server = true
	t.pf = pf
	t.h = *h
	return nil
}

func (t *WebsocketServer) Start() error {
	t.InfoWF("WebsocketServer Start begin server", zap.String("addr", t.addr))
	go hub.run(t.FKLogI)
	go t.beginServer()
	return nil
}
