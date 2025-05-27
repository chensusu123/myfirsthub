package process

import (
	"maze_game_server/lib/codec"
	"maze_game_server/usecase/online"

	"github.com/lonng/nano"
	"github.com/lonng/nano/serialize/json"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type NanoInitService struct {
	nlisten func()
}

// Name implements fkcore.FKServiceI.
func (ns *NanoInitService) Name() string {
	return "nano-init-service"
}

// OnInit implements fkcore.FKServiceI.
func (ns *NanoInitService) OnInit(logger fklog.FKLogI, config fkconfig.FkConfigerI) (err error) {
	// Nano组件与路由
	comps, routes := Components()
	ns.nlisten = func() {
		nano.Listen(":5998",
			nano.WithDebugMode(),

			// 启用WebSocket协议
			nano.WithIsWebsocket(true),
			nano.WithWSPath("/json",
				// 以下Serializer与PacketCodec作用于局部
				codec.NewJsonPacketCodec(routes, codec.WithSerializer(json.NewSerializer())),
			),
			nano.WithWSPath("/pb",
				// 以下Serializer与PacketCodec作用于局部
				codec.NewEsPacketCodec(routes, codec.WithSerializer(codec.NewProtobufSerializer())),
			),
			nano.WithSessionMonitor(online.SessionMonitor()),
			nano.WithComponents(comps),
		)
	}
	return
}

// OnStart implements fkcore.FKServiceI.
func (ns *NanoInitService) OnStart(logger fklog.FKLogI, config fkconfig.FkConfigerI) error {
	go ns.nlisten()
	return nil
}

// OnStop implements fkcore.FKServiceI.
func (ns *NanoInitService) OnStop(logger fklog.FKLogI) error {
	return nil
}

// OnFinish implements fkcore.FKServiceI.
func (ns *NanoInitService) OnFinish(logger fklog.FKLogI) error {
	return nil
}
