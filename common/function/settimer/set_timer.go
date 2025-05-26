package settimer

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/io/rpc/setseataskrpc"
	"maze_game_server/pb/server/SeaTaskSvr"
)

// 设置定时器
func SetTaskExpire(logger fklog.FKLogI, userID uint64, typ SeaTaskSvr.SEA_TASK_TYPE, timeout int64, data []byte) (err error) {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(typ))
	info.Context = data
	info.UserId = &userID
	info.Time = proto.Uint32(uint32(timeout))

	err = setseataskrpc.SetSeaTask(logger, userID, uint64(time.Now().UnixNano()), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, info)
	return
}

func AddTimer(logger fklog.FKLogI, uid, session uint64, info *SeaTaskSvr.TaskInfo) error {
	return setseataskrpc.SetSeaTask(logger, uid, session, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, info)
}

func RemoveTaskTimer(logger fklog.FKLogI, userID uint64, typ SeaTaskSvr.SEA_TASK_TYPE, timeout int64, data []byte) error {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(typ))
	info.Context = data
	info.UserId = &userID

	return setseataskrpc.SetSeaTask(logger, userID, uint64(time.Now().UnixNano()), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, info)
}
func RemoveTimer(logger fklog.FKLogI, uid, session uint64, info *SeaTaskSvr.TaskInfo) error {
	return setseataskrpc.SetSeaTask(logger, uid, session, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, info)
}
