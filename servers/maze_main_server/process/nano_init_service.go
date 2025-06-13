package process

import (
	"maze_game_server/lib/codec"
	"maze_game_server/lib/nano"
	"maze_game_server/lib/nano/serialize/json"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/registry"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/utils"
	"gitlab.ifreetalk.com/maze-plate/freetk/plateregistry"
	"go.uber.org/zap"
)

type NanoInitService struct {
	addr    string
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
	listenAddr := ":5998"
	ns.nlisten = func() {
		nano.Listen(listenAddr,
			// nano.WithDebugMode(),

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
	ns.addr = fkconfig.EnvVal.LocalIP + listenAddr
	return
}

// OnStart implements fkcore.FKServiceI.
func (ns *NanoInitService) OnStart(logger fklog.FKLogI, config fkconfig.FkConfigerI) error {
	go ns.nlisten()
	tags := fkserver.GetRegistryMetadata()

	Info := &registry.Info{
		Namespace:   fkconfig.EnvVal.Namespace,
		ServiceName: fkconfig.EnvVal.AppName + ".ws",
		Addr:        utils.NewNetAddr("tcp", ns.addr),
		Tags:        tags,
	}

	regSvrErr := plateregistry.Registry().Register(Info)
	logger.InfoWF("nano init service start, addr: ", zap.Any("regSvrErr", regSvrErr), zap.Any("regInfo", Info))
	return nil
}

// OnStop implements fkcore.FKServiceI.
func (ns *NanoInitService) OnStop(logger fklog.FKLogI) error {
	Info := &registry.Info{
		Namespace:   fkconfig.EnvVal.Namespace,
		ServiceName: fkconfig.EnvVal.AppName + ".ws",
	}

	plateregistry.Registry().Deregister(Info)
	return nil
}

// OnFinish implements fkcore.FKServiceI.
func (ns *NanoInitService) OnFinish(logger fklog.FKLogI) error {
	return nil
}
