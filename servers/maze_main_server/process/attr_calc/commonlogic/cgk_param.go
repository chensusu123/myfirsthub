/*
 * @Author: majian
 * @Date: 2024-07-18 17:45:17
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 20:47:34
 */
package commonlogic

import (
	"sync"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
)

var DollAttrSrcMapCfg sync.Map
var CleanIdLen int32 // 数据超过多长就可以清理0值ID

func init() {
	param.SafeMapParam(&DollAttrSrcMapCfg, "doll:attr:src:desc:map", func() {}, param.ConvInt32, param.ConvString, "人偶buff来源描述")
	param.Int32P(&CleanIdLen, "clean:attr:id:len", 100, "当属性长度超过多长可以清理0值ID")
}
