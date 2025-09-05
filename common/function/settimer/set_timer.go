package settimer

import (
	"context"
	"time"

	"maze_game_server/io/rpc/setseataskrpc"
	"maze_game_server/pb/server/SeaTaskSvr"

	"google.golang.org/protobuf/proto"
)

// 设置定时器
func SetTaskExpire(ctx context.Context, userID uint64, typ SeaTaskSvr.SEA_TASK_TYPE, timeout int64, data []byte) (err error) {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(typ))
	info.Context = data
	info.UserId = &userID
	info.Time = proto.Uint32(uint32(timeout))

	err = setseataskrpc.SetSeaTask(ctx, userID, uint64(time.Now().UnixNano()), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, info)
	return
}

func AddTimer(ctx context.Context, uid, session uint64, info *SeaTaskSvr.TaskInfo) error {
	return setseataskrpc.SetSeaTask(ctx, uid, session, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, info)
}

func RemoveTaskTimer(ctx context.Context, userID uint64, typ SeaTaskSvr.SEA_TASK_TYPE, timeout int64, data []byte) error {
	info := &SeaTaskSvr.TaskInfo{}
	info.Type = proto.Uint32(uint32(typ))
	info.Context = data
	info.UserId = &userID

	return setseataskrpc.SetSeaTask(ctx, userID, uint64(time.Now().UnixNano()), SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, info)
}
func RemoveTimer(ctx context.Context, uid, session uint64, info *SeaTaskSvr.TaskInfo) error {
	return setseataskrpc.SetSeaTask(ctx, uid, session, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, info)
}
