package websocket_service

import (
	"maze_game_server/lib/net/raw_pkg_ctl"
	websocket_service_impl "maze_game_server/lib/net/websocket_service_actor"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet/pkg_ctl"
	"google.golang.org/protobuf/proto"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// PlugTcpService 添加tcp服务
func PlugTcpRawService(regPack func()) {
	if regPack == nil {
		panic("PlugTcpRawService with nil regPack")
	}
	websocket_service_impl.PlugTcpRawService(regPack)
}

// SetPkgCtrl 替换tcp包处理接口
// 如果调用了这个接口. 那么就需要使用替换接口的注册包接口. (处理非es包等.所以这里只定义了处理接口.不定义注册接口)
func SetPkgCtrl(ctl pkg_ctl.TCPPkgCtler) error {
	return nil
}

func SetNetProtocol(prt fknet.FkProtocol) error {
	return nil
}

// RegProc reg msg proc - 没有使用SetPkgCtrl替换情况下可以使用.
func RegProc(packType uint16, pf raw_pkg_ctl.SvrTCPEsRawProc) (err error) {
	return nil
}

// IgnoreProc ignore msg proc
func IgnoreProc(packType uint16) (err error) {
	return nil
}

type Unmarshaler interface {
	Unmarshal(data []byte) error
}

type Marshaller interface {
	Marshal() ([]byte, error)
}

// RegProcSimple 注册tcp接口
func RegProcSimple(rqID uint16, rqMsg proto.Message, rsID uint16, rsMsg proto.Message,
	deal func(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) error,
) error {
	return websocket_service_impl.RegProcSimple(rqID, rqMsg, rsID, rsMsg, deal)
}

// func SendBytes(logger fklog.FKLogI, shardingID uint64, packType uint16, pack []byte) error {
// 	return nil
// }

func SendPacket(logger fklog.FKLogI, shardingID uint64, packType uint16, pack interface{}) error {
	return websocket_service_impl.SendPacket(logger, shardingID, packType, pack)
}

func MockOnInit(logger fklog.FKLogI, addr string) (err error) {
	// 读取tcp配置
	websocket_service_impl.MockOnInit(logger, addr)
	return
}

func MockOnStart(logger fklog.FKLogI) (err error) {
	logger.InfoWF("websocket-service start begin")
	websocket_service_impl.MockOnStart(logger)
	logger.InfoWF("websocket-service start success.")
	return
}
