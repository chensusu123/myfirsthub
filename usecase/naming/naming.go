package naming

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/discovery"
	"go.uber.org/zap"
)

// ClientSuite It is used to assemble multiple associated client's Options
type ClientSuite struct {
	fklog.FKLogI
	CfgName      string
	DstNameSpace string             // dest namespace for service discovery
	Resolver     discovery.Resolver // service discovery component
	// report service call result for circuitbreak
}

func NewClientSuite(cfgName string) *ClientSuite {
	return &ClientSuite{
		CfgName: cfgName,
	}
}

func (n *ClientSuite) GetServer(typeID int32, groupID int32) []string {
	if n.Resolver == nil {
		return []string{}
	}
	x, err := n.Resolver.Resolve(context.Background(), "minigame-test:config_data.rpc")
	if err != nil {
		n.ErrorWF("GetServer error", zap.Error(err))
		return []string{}
	}
	for _, v := range x.Instances {
		n.InfoWF("GetServer ", zap.Any("server", v.Address()))
		return []string{v.Address().String()}
	}
	// n.InfoWF("GetServer", zap.Any("server", x))
	return []string{}
}

func (n *ClientSuite) GetDB(typeID int32, groupID int32) string {
	n.InfoWF("GetDB", zap.Any("typeID", typeID), zap.Any("groupID", groupID))
	if n.Resolver == nil {
		return ""
	}
	x, err := n.Resolver.Resolve(context.Background(), "minigame-test:maze_main_server.redis")
	if err != nil {
		n.ErrorWF("GetDB resolve error", zap.Error(err))
		return ""
	}

	for _, v := range x.Instances {
		n.InfoWF("GetDB ", zap.Any("server", v.Address()))
		return v.Address().String()
	}
	// n.InfoWF("GetDB ", zap.Any("server", x))
	return ""
}

func (n *ClientSuite) Init() error {
	meta := make(map[string]string)
	meta["GroupID"] = fmt.Sprintf("%d", fkconfig.EnvVal.GroupID)
	o := ClientOptions{
		DstMetadata:  meta,
		SrcNamespace: "minigame-test",
		SrcMetadata:  meta,
	}
	r, err := NewPolarisResolver(o, n.CfgName)
	if err != nil {
		return err
	}
	n.Resolver = r
	logger := fklog.AppLogger().Clone("")
	n.FKLogI = logger
	return nil
}
