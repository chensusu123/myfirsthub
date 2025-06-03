package polarismessvc

import (
	"context"
	"fmt"

	"github.com/polarismesh/polaris-go"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
)

type PolarismesSvc struct {
	provider polaris.ProviderAPI
	fklog.FKLogI
	ctx          context.Context
	cancel       context.CancelFunc
	name         string
	isShutdown   bool
	namespace    string
	service      string
	host         string
	port         uint32
	token        string
	wsInstanceID string
	gmInstanceGm string
}

func NewPolarismesSvc() *PolarismesSvc {
	return &PolarismesSvc{name: "polarismes_svc"}
}

func (ts *PolarismesSvc) Name() string {
	return ts.name
}

func (ts *PolarismesSvc) OnInit(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
	ts.FKLogI = logger.Clone("polarismes_svc")
	ts.ctx, ts.cancel = context.WithCancel(context.Background())
	return
}

func (ts *PolarismesSvc) OnStart(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
	logger.InfoWF("polarismessvc-service start begin")

	provider, err := polaris.NewProviderAPIByFile("./conf.d/polaris.yaml")
	if err != nil {
		logger.ErrorWF("polarismessvc-service start failed.", zap.Error(err))
		return
	}
	// fkconfig.EnvVal.LogDir = logItem.LogDir
	ts.provider = provider
	ts.namespace = "minigame-test"
	// ts.service = "maze_main_server_dev.ws"
	ts.host = fkconfig.EnvVal.LocalIP

	ts.registerService()
	ts.registerServiceGm()
	logger.InfoWF("polarismessvc-service start success.")
	return
}

func (ts *PolarismesSvc) OnStop(logger fklog.FKLogI) (err error) {
	logger.WarnWF("polarismessvc-service stop begin.")
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.deregisterService()
	ts.deregisterServiceGM()
	logger.InfoWF("polarismessvc-service stop success.")
	return
}

func (ts *PolarismesSvc) OnFinish(logger fklog.FKLogI) (err error) {
	logger.InfoWF("polarismessvc-service finish success.")
	return
}

func convertMetadata() map[string]string {
	meta := make(map[string]string)
	meta["GroupID"] = fmt.Sprintf("%d", fkconfig.EnvVal.GroupID)
	meta["SharedID"] = fmt.Sprintf("%d", fkconfig.EnvVal.SharedID)
	meta["ServerType"] = fmt.Sprintf("%d", fkconfig.EnvVal.ServerType)

	meta["K8sNode"] = fkconfig.EnvVal.K8sNode
	meta["LocalIP"] = fkconfig.EnvVal.LocalIP
	meta["AppName"] = fkconfig.EnvVal.AppName

	// meta["internal-enable-set"] = "Y"
	// meta["internal-set-name"] = "xx.yy.sz"

	return meta
}

func (svr *PolarismesSvc) registerService() {
	version := fkserver.Version
	protocol := "ws"
	registerRequest := &polaris.InstanceRegisterRequest{}
	registerRequest.Service = fkconfig.EnvVal.AppName + ".ws"
	registerRequest.Namespace = svr.namespace
	registerRequest.Host = svr.host
	registerRequest.Port = 5998
	registerRequest.Metadata = convertMetadata()
	// registerRequest.ServiceToken = token
	registerRequest.SetTTL(5)
	registerRequest.Version = &version
	registerRequest.Protocol = &protocol
	resp, err := svr.provider.RegisterInstance(registerRequest)
	if err != nil {
		svr.ErrorWF("fail to register instance.", zap.Error(err))
		return
	}
	svr.wsInstanceID = resp.InstanceID
	svr.InfoWF("register instance success.", zap.String("instanceId", resp.InstanceID))
}

func (svr *PolarismesSvc) registerServiceGm() {
	version := fkserver.Version
	protocol := "http"

	registerRequest := &polaris.InstanceRegisterRequest{}
	registerRequest.Service = fkconfig.EnvVal.AppName + ".gm"
	registerRequest.Namespace = svr.namespace
	registerRequest.Host = svr.host
	registerRequest.Port = 5999
	registerRequest.Version = &version
	registerRequest.Protocol = &protocol
	registerRequest.Metadata = convertMetadata()
	// registerRequest.ServiceToken = token
	registerRequest.SetTTL(5)
	resp, err := svr.provider.RegisterInstance(registerRequest)
	if err != nil {
		svr.ErrorWF("fail to register instance.", zap.Error(err))
		return
	}
	svr.gmInstanceGm = resp.InstanceID
	svr.InfoWF("register instance success.", zap.String("instanceId", resp.InstanceID))
}

func (svr *PolarismesSvc) deregisterService() {
	deregisterRequest := &polaris.InstanceDeRegisterRequest{}
	deregisterRequest.Service = fkconfig.EnvVal.AppName + ".ws"
	deregisterRequest.Namespace = svr.namespace
	deregisterRequest.Host = svr.host
	deregisterRequest.Port = 5998
	deregisterRequest.InstanceID = svr.wsInstanceID
	// deregisterRequest.ServiceToken = token
	if err := svr.provider.Deregister(deregisterRequest); err != nil {
		svr.ErrorWF("fail to deregister instance.", zap.Error(err))
		return
	}
	svr.InfoWF("deregister instance success.")
}

func (svr *PolarismesSvc) deregisterServiceGM() {
	deregisterRequest := &polaris.InstanceDeRegisterRequest{}
	deregisterRequest.Service = fkconfig.EnvVal.AppName + ".gm"
	deregisterRequest.Namespace = svr.namespace
	deregisterRequest.Host = svr.host
	deregisterRequest.Port = 5999
	deregisterRequest.InstanceID = svr.gmInstanceGm
	// deregisterRequest.ServiceToken = token
	if err := svr.provider.Deregister(deregisterRequest); err != nil {
		svr.ErrorWF("fail to deregister instance.", zap.Error(err))
		return
	}
	svr.InfoWF("deregister instance success.")
}
