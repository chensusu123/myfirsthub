package raw_pkg_ctl

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet/fkpkg"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/workergroup"
	"maze_game_server/lib/net/raw_pkg"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// SvrTCPEsRawProc 服务端 es包头tcp处理函数
type SvrTCPEsRawProc func(ctx fknet.TCPContext, es *raw_pkg.StruSvrEsRawBaseHead, data []byte) error

// SvrEsRawPackageCtrl es pb包管理器
type SvrEsRawPackageCtrl struct {
	*workergroup.FkWorkGroup
	route     sync.Map
	jsonRoute sync.Map
}

type MiniJsonMsg struct {
	MsgType uint16 `json:"msg_type"`
}

type ErrorJsonMsg struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
}

func (ctl *SvrEsRawPackageCtrl) SetPkgReader(fkpkg.PkgReader) {
}

// RegProc reg msg proc
func (ctl *SvrEsRawPackageCtrl) RegProc(packType uint16, pf SvrTCPEsRawProc) error {
	if _, load := ctl.route.LoadOrStore(packType, pf); load {
		return fmt.Errorf("packType %d is already registered", packType)
	}
	return nil
}

func (ctl *SvrEsRawPackageCtrl) RegJsonProc(packType uint16, pf SvrTCPEsRawProc) error {
	if _, load := ctl.jsonRoute.LoadOrStore(packType, pf); load {
		return fmt.Errorf("packType %d is already registered", packType)
	}
	return nil
}

func (ctl *SvrEsRawPackageCtrl) IgnoreJsonProc(packType uint16) error {
	if _, load := ctl.jsonRoute.LoadOrStore(packType, SvrTCPEsRawProc(func(ctx fknet.TCPContext, es *raw_pkg.StruSvrEsRawBaseHead, data []byte) error {
		ctx.DebugWF("ignore msg.", zap.Uint16("packType", packType))
		return nil
	})); load {
		return fmt.Errorf("packType %d is already registered or ignored", packType)
	}
	return nil
}

// IgnoreProc ignore msg proc
func (ctl *SvrEsRawPackageCtrl) IgnoreProc(packType uint16) error {
	if _, load := ctl.route.LoadOrStore(packType, SvrTCPEsRawProc(func(ctx fknet.TCPContext, es *raw_pkg.StruSvrEsRawBaseHead, data []byte) error {
		ctx.DebugWF("ignore msg.", zap.Uint16("packType", packType))
		return nil
	})); load {
		return fmt.Errorf("packType %d is already registered or ignored", packType)
	}
	return nil
}

// InitWorkGrop 初始化线程池
func (ctl *SvrEsRawPackageCtrl) InitWorkGrop(cacheSize, pipeSize uint32, logger fklog.FKLogI) (err error) {
	////初始化的时候打印tcp包注册信息
	requestIDs := make([]int, 0)
	ctl.route.Range(func(key, value interface{}) bool {
		requestID := key.(uint16)
		requestIDs = append(requestIDs, int(requestID))
		// 添加rs包，rq+1
		requestIDs = append(requestIDs, int(requestID+1))
		return true
	})
	// 排序，主要是便于运维快速添加id
	sort.Ints(requestIDs[:])
	tcpRegisterProm := "register tcp packet id: "
	for _, value := range requestIDs {
		tcpRegisterProm += fmt.Sprintf(" %d", value)
	}
	logger.WarnWF(tcpRegisterProm)

	// ctl.FkWorkGroup = workergroup.NewFkWrokGroup(int(cacheSize), int(pipeSize), logger)
	return
}

// StopWorkGroup 停止工作线程
func (ctl *SvrEsRawPackageCtrl) StopWorkGroup() {
	if ctl.FkWorkGroup != nil {
		ctl.FkWorkGroup.Stop()
	}
}

// Proc 处理tcp包
func (ctl *SvrEsRawPackageCtrl) Proc(ctx fknet.TCPContext, data []byte) (err error) {
	// 防止异常中断
	defer fkalert.RecoverAlertException()
	//
	ctx.FKLogI.SetLogId(time.Now().UnixNano())

	msg := raw_pkg.StruSvrEsRawBaseHead{}
	err = msg.UnPack(data)
	msg.EsRqTime = uint64(time.Now().UnixNano())
	if err != nil {
		ctx.ErrorWF("recv unpack package.", zap.Error(err))
		return
	}
	v, ok := ctl.route.Load(msg.PackType)
	if !ok {
		ctx.WarnWF("recv unreg msg.", zap.Uint16("packType", msg.PackType), zap.Int("len", len(data)))
		// ctx.SendData(data)
		return
	}
	pf, ok := v.(SvrTCPEsRawProc)
	if !ok {
		ctx.WarnWF("invalid reg proc.(func convert failed.)", zap.Uint16("packType", msg.PackType), zap.Int("len", len(data)))
		return
	}
	// 设置用户ID
	ctx.FKLogI.SetUid(uint64(msg.SessionID))

	err = pf(ctx, &msg, msg.Data)
	if err != nil {
		ctx.WarnWF("deal proc failed.", zap.Uint16("packType", msg.PackType),
			zap.Uint32("SessionID", msg.SessionID),
			zap.Int("len", len(data)), zap.Error(err))
		return
	}
	return
}

// Proc 处理tcp包
func (ctl *SvrEsRawPackageCtrl) ProcJson(ctx fknet.TCPContext, data []byte) (err error) {
	// 防止异常中断
	defer fkalert.RecoverAlertException()
	//
	ctx.FKLogI.SetLogId(time.Now().UnixNano())

	msg := MiniJsonMsg{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		ctx.ErrorWF("recv unpack package.", zap.Error(err))
		packt := ErrorJsonMsg{
			MsgType: int(1),
			Data:    "json.Unmarshal failed",
		}
		data, err = json.Marshal(packt)
		if err == nil {
			ctx.SendData(data)
		}
		return
	}
	v, ok := ctl.jsonRoute.Load(msg.MsgType)
	if !ok {
		ctx.WarnWF("recv unreg msg.", zap.Uint16("packType", msg.MsgType), zap.Int("len", len(data)))
		packt := ErrorJsonMsg{
			MsgType: int(1),
			Data:    fmt.Sprintf("msg %d unreg", msg.MsgType),
		}
		data, err = json.Marshal(packt)
		if err == nil {
			ctx.SendData(data)
		}
		return
	}
	pf, ok := v.(SvrTCPEsRawProc)
	if !ok {
		ctx.WarnWF("invalid reg proc.(func convert failed.)", zap.Uint16("packType", msg.MsgType), zap.Int("len", len(data)))
		return
	}
	// 设置用户ID
	// ctx.FKLogI.SetUid(uint64(msg.SessionID))

	tcpPacket := raw_pkg.StruSvrEsRawBaseHead{
		PackType: msg.MsgType,
	}
	err = pf(ctx, &tcpPacket, data)
	if err != nil {
		ctx.WarnWF("deal proc failed.", zap.Uint16("packType", msg.MsgType),
			zap.Int("len", len(data)), zap.Error(err))
		return
	}
	return
}
