package timertask

import (
	"fmt"
	"time"

	"maze_game_server/io/rpc/setseataskrpc"
	"maze_game_server/pb/server/SeaTaskSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
)

func newSessionID() uint64 {
	return uint64(time.Now().UnixNano())
}

func SetTask(logger fklog.FKLogI, uid uint64, timeout int64, format string, args ...interface{}) (err error) {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(227))
	info.Context = []byte(fmt.Sprintf(format, args...))
	info.UserId = proto.Uint64(uid)
	info.Time = proto.Uint32(uint32(timeout))

	err = setseataskrpc.SetSeaTask(logger, uid, newSessionID(), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, info)
	return
}

func StopTask(logger fklog.FKLogI, uid uint64, timeout int64, format string, args ...interface{}) (err error) {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(227))
	info.Context = []byte(fmt.Sprintf(format, args...))
	info.UserId = proto.Uint64(uid)
	info.Time = proto.Uint32(uint32(timeout))

	err = setseataskrpc.SetSeaTask(logger, uid, newSessionID(), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, info)
	return
}
