/*
 * @Author: majian
 * @Date: 2024-03-13 18:00:04
 * @Last Modified by: majian
 * @Last Modified time: 2024-03-13 18:00:45
 */

package httpalert

import (
	"fmt"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func TestAlert(t *testing.T) {
	var logger = fklog.AppLogger().Clone("httpalert")
	err := SendHttpAlert(logger, 1, "http_alert_test", "majian", false)
	fmt.Println("err", err)
}
