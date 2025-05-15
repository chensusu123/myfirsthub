package httpalert

import (
	"fmt"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	//"gitlab.ifreetalk.com/maze-plate/io/common/alert"
)

// 使用http方式上报告警，默认有告警都会报出来
// online参数:false 发到golang内网告警群 true 发到告警中心/golang告警群，使用时找郭旺配置下策略,具体发送到哪个告警群,把内网和外网区分的标识发给他
// typ参数: 告警类型 服务内部保证唯一，区分不同的告警
// recevier参数: 接收者列表，多人用逗号分隔
// 告警示例 #190909195 告警 11:35:58 warship-gm-server-c0-g7-0 2024:04:29-11:35:58, ServerID:233731, ServerType:17712, Group:7 content [ 这是一条测试告警 ] @majian
func SendHttpAlert(logger fklog.FKLogI, typ int32, content, recevier string, online bool) error {
	sb := strings.Builder{}
	sb.WriteString(time.Now().Format("2006:01:02-15:04:05"))
	sb.WriteString(fmt.Sprintf(", ServerID:%d, ServerType:%d, Group:%d", fkconfig.EnvVal.ServerID, fkconfig.EnvVal.ServerType, fkconfig.EnvVal.GroupID))
	sb.WriteString(fmt.Sprintf(" content [ %s ]", content))

	typInfo := fmt.Sprintf("inner-%d", typ)
	if online {
		typInfo = fmt.Sprintf("online-%d", typ)
	}
	_ = typInfo
	return nil
	// return alert.NewAlert(logger, fkconfig.EnvVal.AppName, typInfo, fkconfig.EnvVal.HostName, recevier, sb.String())
}
