/*
 * @Author: majian
 * @Date: 2024-07-18 17:45:17
 * @Last Modified by: majian
 * @Last Modified time: 2024-07-24 15:12:29
 */
package attr_calc

import (
	"sync"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig/param"
)

var DollAttrSrcMapCfg sync.Map

func init() {
	param.SafeMapParam(&DollAttrSrcMapCfg, "doll:attr:src:desc:map", func() {}, param.ConvInt32, param.ConvString, "人偶buff来源描述")
}
