package naming

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/discovery"
	namingI "gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
	"go.uber.org/zap"
)

// ClientSuite It is used to assemble multiple associated client's Options
type ClientSuite struct {
	fklog.FKLogI
	CfgName      string
	DstNameSpace string             // dest namespace for service discovery
	Resolver     discovery.Resolver // service discovery component
	// report service call result for circuitbreak
	InitFunList []InitFunc
}

type InitFunc func(namingI.NamingI) error

func NewClientSuite(cfgName string, initFuncList []InitFunc) *ClientSuite {
	return &ClientSuite{
		CfgName:     cfgName,
		InitFunList: initFuncList,
	}
}

func (n *ClientSuite) GetServer(svrName string, groupID int32) ([]string, error) {
	if n.Resolver == nil {
		return []string{}, nil
	}
	desc := fmt.Sprintf("%s:%s", "minigame-test", svrName)
	x, err := n.Resolver.Resolve(context.Background(), desc)
	if err != nil {
		n.ErrorWF("GetServer error", zap.Error(err))
		return []string{}, nil
	}
	for _, v := range x.Instances {
		n.InfoWF("GetServer ", zap.Any("server", v.Address()))
		return []string{v.Address().String()}, nil
	}
	// n.InfoWF("GetServer", zap.Any("server", x))
	return []string{}, nil
}

func (n *ClientSuite) GetDB(dbName string, groupID int32) (string, error) {
	n.InfoWF("GetDB", zap.Any("dbName", dbName), zap.Any("groupID", groupID))
	if n.Resolver == nil {
		return "", nil
	}
	desc := fmt.Sprintf("%s:%s", "minigame-test", dbName)
	x, err := n.Resolver.Resolve(context.Background(), desc)
	if err != nil {
		n.ErrorWF("GetDB resolve error", zap.Error(err))
		return "", nil
	}

	for _, v := range x.Instances {
		n.InfoWF("GetDB ", zap.Any("dbName", dbName), zap.Any("server", v.Address()))
		return v.Address().String(), nil
	}
	// n.InfoWF("GetDB ", zap.Any("server", x))
	return "", nil
}

func (n *ClientSuite) GetDBUserAndPassword(dbName string, groupID int32) (string, string, error) {
	n.InfoWF("GetDBUserAndPassword", zap.Any("dbName", dbName), zap.Any("groupID", groupID))

	return "majiange", "162>wind", nil
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
	for _, f := range n.InitFunList {
		err = f(n)
		if err != nil {
			return err
		}
	}
	return nil
}
