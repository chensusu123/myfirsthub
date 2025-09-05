package main

import (
	"context"
	"errors"

	"maze_game_server/lib/codec"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeRobGuaJi"
	"maze_game_server/servers/maze_es/svrheadertrace"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/nanotrace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type (
	HttpDemo2 struct{}
)

func (h *HttpDemo2) Handle(hSvr *server.Hertz) error {
	hSvr.GET("/config_manager", func(ctx context.Context, c *app.RequestContext) {
		GetData(ctx)
		GetDataCtx(ctx)
		c.String(consts.StatusOK, "config_manager ")
	})
	hSvr.GET("/span", func(ctx context.Context, c *app.RequestContext) {
		SpancCall1(ctx)
		SpancCall3(ctx)
		c.String(consts.StatusOK, "span ")
	})

	return nil
}

func GetData(ctx context.Context) {
	GetDataCtx(ctx)
}

func GetDataCtx(ctx context.Context) {
	config_manager.MissRecord(ctx, "errr", 2)
}

func SpancCall1(ctx context.Context) {
	span := nanotrace.NewSimpleTrace("SpancCall1")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)

	_, routes := Components()
	xx := codec.NewEsPacketCodec(routes, codec.WithSerializer(codec.NewProtobufSerializer()))

	k1 := make([]byte, 0)
	data := svrheadertrace.YYY(ctx, 10490, &MazeRobGuaJi.MazeRobGuaJiRQ{
		RobUser: proto.Uint64(1000000000000000000),
	}, xx)

	k1 = append(k1, data...)

	a, b, c := xx.Decode(k1)
	_, _, _ = a, b, c
	fklog.ContextAppLogger(ctx).CtxInfo(ctx, "",
		zap.Any("ac", c), zap.Any("b", b))

	// SpancCall2(ctx)
}

func SpancCall2(ctx context.Context) {
	span := nanotrace.NewSimpleTrace("SpancCall2")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	span.RecordError(ctx, errors.New("SpancCall2 err"))
}

func SpancCall3(ctx context.Context) {
	span := nanotrace.NewSimpleTrace("SpancCall3")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
}

func Components() (comps *component.Components, routes *codec.Routes) {
	comps = &component.Components{}
	routes = &codec.Routes{}
	// 注册nano组件与自定义包解析路由
	reg := func(comp component.Component) {
		comps.Register(comp)
		routes.Register(comp)
	}

	{
		reg(NewRob()) // IM聊天组件
	}

	return
}

type Rob struct {
	component.Base
}

func NewRob() *Rob {
	return &Rob{}
}

func (*Rob) OnMazeRobGuaJiRQ_10490_10491(s *session.Session, req *MazeRobGuaJi.MazeRobGuaJiRQ) (err error) {
	return nil
}
