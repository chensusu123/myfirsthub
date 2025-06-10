package business

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
)

const (
	alertOpenExcelFailed fkalert.AlertID = fkalert.CustomAlert + iota
	alertPbUnmarsal
	alertRequestDbIDNotZero
	alertSheetNotExists
	alertGetColumnDataFailed
	alertDataTooBig
	// AlertTest 临时.测试告警
	AlertTest
)

// alertMsg 告警信息
var alertMsg = map[fkalert.AlertID]string{
	alertOpenExcelFailed:     "打开excel文件失败",
	alertPbUnmarsal:          "pb解析失败",
	alertRequestDbIDNotZero:  "服务提供excel文件数据.不能读取数据库.",
	alertSheetNotExists:      "查询的sheet表不存在",
	alertGetColumnDataFailed: "读取指定列数据失败,列不匹配.请检查版本",
	alertDataTooBig:          "配置数据太大了.一条都发不出去,请修改配置",
	AlertTest:                "测试告警",
}

func init() {
	// 添加告警初始化
	fkalert.AppendAlertMsg(alertMsg)
	// 测试提示信息注册上报
	fkconfig.RegTipInfo(1, "测试注册上报", "ENUM_CGK_ERROR_test_tip_info")
}

// var monitorLoadRPC = fkmonitor.DefaultTimeCheckMonitor("RPC.LoadSuccess", "RPC.LoadFailed")
// var monitorLoadFile = fkmonitor.DefaultTimeMonitor("FileMonitor.ParseExcel")
