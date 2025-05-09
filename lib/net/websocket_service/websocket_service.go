package websocket_service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg_ctl"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/plate/freetk/fkcore"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkmonitor"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet/pkg_ctl"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// PlugTcpService 添加tcp服务
func PlugTcpRawService(regPack func()) {
	if regPack == nil {
		panic("PlugTcpRawService with nil regPack")
	}
	fkserver.AppServer.AddServeService(gGlobalTCPRawServer)
	fkconfig.SetTcpMonitor()
	regPack()
}

type tTCPRawService struct {
	fklog.FKLogI
	*WebsocketServer
	pkgCtl        pkg_ctl.TCPPkgCtler
	protocol      fknet.FkProtocol
	initFlag      atomic.Bool
	locker        sync.RWMutex
	baseSessionID int64
}

// newTCPServer 新建tcp服务器
func newTCPServer(pkgCtl pkg_ctl.TCPPkgCtler) *tTCPRawService {
	ts := &tTCPRawService{}
	ts.WebsocketServer = NewWebsocketServer()
	ts.pkgCtl = pkgCtl
	ts.initFlag.Store(false)
	return ts
}

// FKServiceI 服务接口
func (ts *tTCPRawService) Name() string {
	return "websocket-service"
}

func (ts *tTCPRawService) OnInit(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
	// 读取tcp配置
	tcpCfg := fkconfig.GetServerConfig()
	logger.InfoWF("websocket-service init begin.", zap.String("addr", tcpCfg.GetThriftRPCAddr()))
	// 检查tcp接口
	if tcpCfg.GetThriftRPCAddr() == "" {
		err = errors.New("websocket-service port is 0")
		fkfmt.Println("websocket-service port is 0.", tcpCfg)
		logger.ErrorWF("websocket-service port is 0.", zap.Any("addr", tcpCfg.GetThriftRPCAddr()))
		return
	}
	// 绑定时，不指定ip地址
	_, port, err := net.SplitHostPort(tcpCfg.GetThriftRPCAddr())
	if err != nil {
		fkfmt.Println("websocket-service split ip/port failed.", tcpCfg.GetThriftRPCAddr(), err)
		logger.ErrorWF("websocket-service split ip/port failed.", zap.String("addr", tcpCfg.GetThriftRPCAddr()), zap.Error(err))
		return
	}
	addr := ":" + port
	// 启动tcp服务器
	err = ts.WebsocketServer.Init(addr, func() fknet.FkProtocol {
		return ts
	})
	if err != nil {
		fkfmt.Println("websocket-service listen failed.", tcpCfg.GetThriftRPCAddr(), err)
		logger.ErrorWF("websocket-service listen failed.", zap.String("addr", tcpCfg.GetThriftRPCAddr()), zap.Error(err))
		return
	}
	// 初始化线程池
	ts.initFlag.Store(true)
	ts.pkgCtl.InitWorkGrop(cacheSize, threadCount, logger)
	fkfmt.Println("websocket-service init success.", tcpCfg.GetThriftRPCAddr(), tcpCfg)
	logger.InfoWF("websocket-service init success.", zap.String("addr", tcpCfg.GetThriftRPCAddr()))
	return
}

func (ts *tTCPRawService) OnStart(logger fklog.FKLogI, config fkcore.FkConfigerI) (err error) {
	logger.InfoWF("websocket-service start begin")
	// 启动服务
	err = ts.WebsocketServer.Start()
	if err != nil {
		fkfmt.Println("websocket-service start tcp server failed.")
		logger.ErrorWF("websocket-service start tcp server failed.", zap.Error(err))
		return
	}
	logger.InfoWF("websocket-service start success.")
	return
}

func (ts *tTCPRawService) OnStop(logger fklog.FKLogI) (err error) {
	logger.WarnWF("websocket-service stop begin.")
	go func() {
		ts.WebsocketServer.StopServer()
		logger.WarnWF("websocket-service stop end.")
	}()

	time.Sleep(500 * time.Millisecond)
	logger.InfoWF("websocket-service stop success.")
	return
}

func (ts *tTCPRawService) OnFinish(logger fklog.FKLogI) (err error) {
	logger.InfoWF("websocket-service finish success.")
	return
}

func (ts *tTCPRawService) OnTransportMade(conn fknet.FkTransport) {
	ts.WarnWF("websocket-service recv conn", zap.Stringer("addr", conn.RemoteAddr()))

	if ts.protocol != nil {
		ts.protocol.OnTransportMade(conn)
	}
}

func (ts *tTCPRawService) OnTransportLost(conn fknet.FkTransport) {
	ts.WarnWF("websocket-service lose conn", zap.Stringer("addr", conn.RemoteAddr()))

	if ts.protocol != nil {
		ts.protocol.OnTransportLost(conn)
	}
}

func (ts *tTCPRawService) OnTransportData(conn fknet.FkTransport, data []byte) {
	ts.pkgCtl.Proc(fknet.NewTCPContext(context.TODO(), conn), data)
}

// SetPkgCtrl 替换tcp包处理接口
// 如果调用了这个接口. 那么就需要使用替换接口的注册包接口. (处理非es包等.所以这里只定义了处理接口.不定义注册接口)
func (ts *tTCPRawService) SetPkgCtrl(ctl pkg_ctl.TCPPkgCtler) error {
	// 服务已经初始化.不允许替换接口
	if ts.initFlag.Load() {
		panic("websocket-service is inited before. You Shoud Call SetPkgCtrl in func init().")
	}
	ts.pkgCtl = ctl
	if ts.FKLogI != nil {
		ts.FKLogI.InfoWF("websocket-service update pkg ctrl.")
	}
	fkfmt.Println("websocket-service update pkg ctrl.")
	return nil
}

func (ts *tTCPRawService) SetNetProtocol(prt fknet.FkProtocol) error {
	ts.protocol = prt
	return nil
}

func init() {
}

var (
	gGlobalTCPRawServer    *tTCPRawService
	cacheSize, threadCount uint32
)

func init() {
	// 登陆包忽略
	gDefaultTCPPkgCtl.IgnoreProc(43001)
	gDefaultTCPPkgCtl.IgnoreProc(43003)
	// 初始化tcp服务器
	gGlobalTCPRawServer = newTCPServer(&gDefaultTCPPkgCtl)
	param.Uint32P(&cacheSize, "cache:size", 1000, "tcp包缓冲大小")
	param.Uint32P(&threadCount, "thread:count", 1024, "处理线程数量")

	// 添加tcp服务
}

// SetPkgCtrl 替换tcp包处理接口
// 如果调用了这个接口. 那么就需要使用替换接口的注册包接口. (处理非es包等.所以这里只定义了处理接口.不定义注册接口)
func SetPkgCtrl(ctl pkg_ctl.TCPPkgCtler) error {
	return gGlobalTCPRawServer.SetPkgCtrl(ctl)
}

func SetNetProtocol(prt fknet.FkProtocol) error {
	return gGlobalTCPRawServer.SetNetProtocol(prt)
}

// 默认tcp包处理接口. 处理es包.
var gDefaultTCPPkgCtl raw_pkg_ctl.SvrEsRawPackageCtrl

// RegProc reg msg proc - 没有使用SetPkgCtrl替换情况下可以使用.
func RegProc(packType uint16, pf raw_pkg_ctl.SvrTCPEsRawProc) (err error) {
	err = gDefaultTCPPkgCtl.RegProc(packType, pf)
	if err != nil {
		fkfmt.Println("websocket-service reg tcp packet ", packType, " failed.", err)
		fkutil.PanicIfError(err)
	}
	// 注册包ID
	fkconfig.RegServicePacket(uint32(packType))
	return nil
}

// IgnoreProc ignore msg proc
func IgnoreProc(packType uint16) (err error) {
	err = gDefaultTCPPkgCtl.IgnoreProc(packType)
	if err != nil {
		fkfmt.Println("websocket-service ingore tcp packet ", packType, " failed.", err)
		fkutil.PanicIfError(err)
	}
	// 注册包ID
	fkconfig.RegServicePacket(uint32(packType))
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
	// 反射类型.用于新建对象
	rqMsgType := reflect.TypeOf(rqMsg).Elem()
	rsMsgType := reflect.TypeOf(rsMsg).Elem()
	// 类型名称
	rqTypeString := rqMsgType.String()
	rsTypeString := rsMsgType.String()
	// 函数名
	funcName := fkutil.FunName(deal)
	// 监控
	dealMonitor := fkmonitor.DefaultTimeCheckMonitor(funcName+"Success", funcName+"Failed")
	// dealMonitor := fkmonitor.DefaultFrequencyMonitor(funcName)
	// 告警提示信息
	alertTip := fmt.Sprintf("rq:%d %s", rqID, rqTypeString)

	// 封包函数
	handler := func(ctx fknet.TCPContext, es *raw_pkg.StruSvrEsRawBaseHead, data []byte) (err error) {
		// 监控
		defer dealMonitor.Start(func() bool {
			return err == nil
		})()
		// defer dealMonitor.Incr()

		// 新建请求包
		rq := reflect.New(rqMsgType).Elem().Addr().Interface().(proto.Message)
		rs := reflect.New(rsMsgType).Elem().Addr().Interface().(proto.Message)

		// 解包
		err = proto.Unmarshal(data, rq)
		if err != nil {
			fkalert.Alert(fkalert.AlertTCPUnMarshalFailed, alertTip)
			// 解包失败.丢弃.
			ctx.ErrorWF("unmarshal tcp request failed.",
				zap.Uint16("rqID", rqID), zap.Uint16("rsID", rsID), zap.String("func", funcName),
				zap.String("rqType", rqTypeString), zap.String("rsType", rsTypeString),
				zap.Int("len", len(data)),
				zap.Uint64("sessionID", uint64(es.SessionID)),
				zap.Error(err))
			return
		}

		pkg := &raw_pkg.StruSvrEsRawBaseHead{}
		pkg.PackType = rsID
		pkg.SessionID = es.SessionID
		pkg.EsRqTime = es.EsRqTime

		needRS := false
		// 回包
		defer func() {
			if !needRS {
				return
			}
			dataPB, errMarshal := proto.Marshal(rs)
			if errMarshal != nil {
				ctx.ErrorWF("send resMarshalpond failed.",
					zap.Uint16("rqID", rqID), zap.Uint16("rsID", rsID), zap.String("func", funcName),
					zap.String("rqType", rqTypeString), zap.String("rsType", rsTypeString),
					zap.Int("len", len(data)),
					zap.Uint64("sessionID", uint64(es.SessionID)),
					zap.Error(errMarshal))
				return
			}
			pkg.Data = dataPB
			pkg.EsRsTime = uint64(time.Now().UnixNano())
			err = ctx.Send(pkg)
			if err != nil {
				ctx.ErrorWF("send respond failed.",
					zap.Uint16("rqID", rqID), zap.Uint16("rsID", rsID), zap.String("func", funcName),
					zap.String("rqType", rqTypeString), zap.String("rsType", rsTypeString),
					zap.Int("len", len(data)),
					zap.Uint64("sessionID", uint64(es.SessionID)),
					zap.Error(err))
				return
			}
		}()

		foundUserID := uint64(0)
		if es.PackType != 5183 {
			// 非登录包
			tmpUserID := ctx.GetTag("userID")
			switch v := tmpUserID.(type) {
			case int64:
				foundUserID = uint64(v)
			case uint64:
				foundUserID = v
			}
			if foundUserID == 0 {
				ctx.WarnWF("user not login", zap.Uint16("rqID", rqID))
				return
			}
			ctx.SetUid(foundUserID)
		}

		needRS = true
		// 处理请求
		err = deal(ctx, foundUserID, rq, rs)
		if err != nil {
			ctx.ErrorWF("deal rpc request failed.",
				zap.Uint16("rqID", rqID), zap.Uint16("rsID", rsID), zap.String("func", funcName),
				zap.String("rqType", rqTypeString), zap.String("rsType", rsTypeString),
				zap.Int("len", len(data)), zap.Uint64("sharedingID", foundUserID),
				zap.Uint64("sessionID", uint64(es.SessionID)),
				zap.Error(err))
		}
		if es.PackType == 5183 {
			// 登录包
			tmpUserID := ctx.GetTag("userID")
			switch v := tmpUserID.(type) {
			case int64:
				foundUserID = uint64(v)
			case uint64:
				foundUserID = v
			}
			if foundUserID != 0 {
				// 登录成功
				hub.OnLogin(foundUserID, ctx)
			}
		}
		return
	}

	err := gDefaultTCPPkgCtl.RegProc(rqID, raw_pkg_ctl.SvrTCPEsRawProc(handler))
	if err != nil {
		fkfmt.Println("websocket-service reg tcp packet ", rqID, " failed.", err)
		fkutil.PanicIfError(err)
	}
	// 注册包ID
	fkconfig.RegServicePacket(uint32(rqID))
	return nil
}

func SendPb(logger fklog.FKLogI, shardingID uint64, packType uint16, pack proto.Message) error {
	var sessionID int64
	pb, err := proto.Marshal(pack)
	if err != nil {
		logger.ErrorWF("SendPb Marshal failed",
			zap.Any("err", err),
			zap.Any("shardingID", shardingID), zap.Any("sessionID", sessionID),
			zap.Any("packType", packType), zap.Any("pack", pack))
		return err
	}

	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.SessionID = uint32(sessionID)
	pkg.PackType = packType
	pkg.EsRsTime = uint64(time.Now().Unix())
	pkg.Data = pb
	pkg.SetTeaflag()

	data, err := pkg.Pack()
	if err != nil {
		logger.ErrorWF("SendPb Pack failed",
			zap.Any("err", err),
			zap.Any("shardingID", shardingID), zap.Any("sessionID", sessionID),
			zap.Any("packType", packType), zap.Any("pack", pack))
		return err
	}

	err = gGlobalTCPRawServer.SendData(logger, int64(shardingID), sessionID, data)
	if err != nil {
		logger.ErrorWF("SendPb SendData failed",
			zap.Any("err", err),
			zap.Any("shardingID", shardingID), zap.Any("sessionID", sessionID),
			zap.Any("packType", packType), zap.Any("pack", pack))
		return err
	}
	return err
}

func SendBytes(logger fklog.FKLogI, shardingID uint64, packType uint16, pack []byte) error {
	var sessionID int64
	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.SessionID = uint32(sessionID)
	pkg.PackType = packType
	pkg.Data = pack
	pkg.EsRsTime = uint64(time.Now().Unix())
	pkg.SetTeaflag()

	data, err := pkg.Pack()
	if err != nil {
		logger.ErrorWF("SendData Pack failed",
			zap.Any("err", err),
			zap.Any("shardingID", shardingID), zap.Any("sessionID", sessionID),
			zap.Any("packType", packType), zap.Any("pack", pack))
		return err
	}

	err = gGlobalTCPRawServer.SendData(logger, int64(shardingID), sessionID, data)
	if err != nil {
		logger.ErrorWF("SendData SendData failed",
			zap.Any("err", err),
			zap.Any("shardingID", shardingID), zap.Any("sessionID", sessionID),
			zap.Any("packType", packType), zap.Any("pack", pack))
		return err
	}
	return err
}

func (ts *tTCPRawService) SendData(logger fklog.FKLogI, userID int64, sessionID int64, data []byte) error {
	if ts == nil {
		return errors.New("tTCPService == nil")
	}
	hub.SendDataByUserID(logger, uint64(userID), data)
	return nil
}

func MockOnInit(logger fklog.FKLogI, addr string) (err error) {
	// 读取tcp配置
	tcpCfg := fkconfig.GetServerConfig()
	logger.InfoWF("websocket-service init begin.", zap.String("addr", tcpCfg.GetThriftRPCAddr()))
	gGlobalTCPRawServer.FKLogI = logger.Clone("websocket-service")
	// 启动tcp服务器
	err = gGlobalTCPRawServer.WebsocketServer.Init(addr, func() fknet.FkProtocol {
		return gGlobalTCPRawServer
	})
	if err != nil {
		fkfmt.Println("websocket-service listen failed.", tcpCfg.GetThriftRPCAddr(), err)
		logger.ErrorWF("websocket-service listen failed.", zap.String("addr", tcpCfg.GetThriftRPCAddr()), zap.Error(err))
		return
	}
	// 初始化线程池
	gGlobalTCPRawServer.initFlag.Store(true)
	gGlobalTCPRawServer.pkgCtl.InitWorkGrop(cacheSize, threadCount, logger)
	fkfmt.Println("websocket-service init success.", tcpCfg.GetThriftRPCAddr(), tcpCfg)
	logger.InfoWF("websocket-service init success.", zap.String("addr", tcpCfg.GetThriftRPCAddr()))
	return
}

func MockOnStart(logger fklog.FKLogI) (err error) {
	logger.InfoWF("websocket-service start begin")
	// 启动服务
	err = gGlobalTCPRawServer.WebsocketServer.Start()
	if err != nil {
		fkfmt.Println("websocket-service start tcp server failed.")
		logger.ErrorWF("websocket-service start tcp server failed.", zap.Error(err))
		return
	}
	logger.InfoWF("websocket-service start success.")
	return
}
