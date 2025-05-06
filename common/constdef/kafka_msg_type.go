/*
 * @Author: majian
 * @Date: 2024-08-12 20:08:21
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-19 20:51:21
 * @Desc kafka 消息分发消息类型定义
 */
package constdef

const (
	KafkaMDTSex = 13 // 性别变化
)

var (
	KafkaMDTSexDesc = "maze_mdtsex_consumer"
)
var KafkaNameMap = map[string]int32{
	KafkaMDTSexDesc: KafkaMDTSex, // 性别变化
}

var KafkaNameMapRe = map[int32]string{
	KafkaMDTSex: KafkaMDTSexDesc,
}

func GetMsgType(name string) int32 {
	return KafkaNameMap[name]
}
