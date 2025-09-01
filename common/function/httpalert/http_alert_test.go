/*
 * @Author: majian
 * @Date: 2024-03-13 18:00:04
 * @Last Modified by: majian
 * @Last Modified time: 2024-03-13 18:00:45
 */

package httpalert

import (
	"context"
	"fmt"
	"testing"
)

func TestAlert(t *testing.T) {
	ctx := context.Background()
	err := SendHttpAlert(ctx, 1, "http_alert_test", "majian", false)
	fmt.Println("err", err)
}
