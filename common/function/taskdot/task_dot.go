/*
 * @Author: majian
 * @Date: 2024-04-02 14:00:32
 * @Last Modified by: majian
 * @Last Modified time: 2024-08-13 14:35:29
 */
package taskdot

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/io/kafka_interface/common/TaskDataKafka"
)

// 精简不需要的字段
func SendTask(logger fklog.FKLogI, uid uint64, taskId int32, itemId int64, value int64) {
	TaskDataKafka.SendData(logger, uid, taskId, itemId, value, "", 0)
}
