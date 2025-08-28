package frame_service

import (
	"context"
	"sort"
	"time"

	"maze_game_server/common/errors"
	"maze_game_server/model/frame_model"
	"maze_game_server/pb/common/MazeRoom"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const GameTickMin = 33 // 单位ms

// 全局服务器实例
var RoomServer *frame_model.Server

func init() {
	RoomServer = frame_model.NewServer()
}

func checkGameTick(logger fklog.FKLogI, gameTick int32) error {
	// 最低帧数限制
	if gameTick < int32(GameTickMin) {
		logger.ErrorWF("checkGameTick gameTick invalid", zap.Int32("gameTick", gameTick))
		return errors.New("game tick invalid")
	}
	return nil
}

func getRoomById(logger fklog.FKLogI, roomId string) (room *frame_model.Room, err error) {
	if len(RoomServer.Rooms) == 0 {
		logger.ErrorWF("getRoomById RoomServer.Rooms is nil")
		return nil, errors.New("room not found")
	}

	room, ok := RoomServer.Rooms[roomId]
	if !ok {
		logger.ErrorWF("getRoomById roomId invalid", zap.String("roomId", roomId))
		return nil, errors.New("room not found")
	}
	return room, nil
}

func GetRoomByPlayerId(logger fklog.FKLogI, userId uint64) (room *frame_model.Room, err error) {
	roomId, ok := RoomServer.PlayerRoom[userId]
	if !ok {
		logger.ErrorWF("GetRoomByPlayerId userId invalid", zap.Uint64("userId", userId))
		return nil, errors.New("roomId not found")
	}
	room, ok = RoomServer.Rooms[roomId]
	if !ok {
		logger.ErrorWF("GetRoomByPlayerId roomId invalid", zap.String("roomId", roomId))
		return nil, errors.New("room not found")
	}
	return room, nil
}

func getRoomPlayer(logger fklog.FKLogI, room *frame_model.Room, userId uint64) (player *frame_model.Player, err error) {
	player, ok := room.MemberMap[userId]
	if !ok {
		logger.ErrorWF("getRoomPlayer userId invalid", zap.Uint64("userId", userId))
		return nil, errors.New("player not found in room")
	}
	return player, nil
}

func WritePump(logger fklog.FKLogI, userId uint64, message string) error {
	room, err := GetRoomByPlayerId(logger, userId)
	if err != nil {
		logger.ErrorWF("WritePump GetRoomByPlayerId failed", zap.Error(err))
		return errors.New("room not found")
	}

	player, err := getRoomPlayer(logger, room, userId)
	if err != nil {
		logger.ErrorWF("WritePump getRoomPlayer failed", zap.Error(err))
		return errors.New("player not found in room")
	}

	addPump(player, room, message)
	return nil
}

// 玩家消息添加到消息队列
func addPump(p *frame_model.Player, room *frame_model.Room, message string) {
	msg := frame_model.NewInput(message)
	var inputs []frame_model.Input
	inputs = append(inputs, *msg)
	// 添加到房间输入队列
	addInputs(room, p.ID, inputs)
}

// 获取或创建房间
func CreateRoom(logger fklog.FKLogI, userId uint64, userIdList []uint64, gameTick, startPercent, gameLastTime, udpReliabilityStrategy int32, roomExtInfo string, needSeed bool) (room *frame_model.Room, err error) {
	err = checkGameTick(logger, gameTick)
	if err != nil {
		logger.ErrorWF("CreateRoom checkGameTick failed", zap.Error(err))
		return nil, err
	}
	room, err = GetRoomByPlayerId(logger, userId)
	if err != nil {
		logger.ErrorWF("CreateRoom getRoomByPlayerId failed", zap.Error(err))
		return nil, err
	}

	RoomServer.Mutex.Lock()
	defer RoomServer.Mutex.Unlock()
	room = frame_model.NewRoomInfo(logger, userIdList, gameTick, startPercent, gameLastTime, udpReliabilityStrategy, roomExtInfo, needSeed)
	RoomServer.Rooms[room.RoomIdStr] = room
	addRoomMember(logger, room, userId, userIdList)
	return room, nil
}

func addRoomMember(logger fklog.FKLogI, room *frame_model.Room, userId uint64, userIdList []uint64) {
	master := frame_model.NewPlayer(logger, userId, frame_model.ROLE_MASTER)
	addPlayer(logger, room, master)
	for i := 0; i < len(userIdList); i++ {
		player := frame_model.NewPlayer(logger, userIdList[i], frame_model.ROLE_COMMON)
		addPlayer(logger, room, player)
	}
}

func JoinRoom(logger fklog.FKLogI, userId uint64, roomId string) (err error) {
	room, err := getRoomById(logger, roomId)
	if err != nil {
		logger.ErrorWF("JoinRoom getRoomById failed", zap.Uint64("userId", userId), zap.String("roomId", roomId))
		return errors.New("room not found")
	}
	player, _ := getRoomPlayer(logger, room, userId)
	if player != nil {
		logger.ErrorWF("JoinRoom getRoomPlayer already join room", zap.Uint64("userId", userId))
		return errors.New("already join room")
	}

	player = frame_model.NewPlayer(logger, userId, frame_model.ROLE_COMMON)
	addPlayer(logger, room, player)
	return nil
}

// 添加玩家到房间
func addPlayer(logger fklog.FKLogI, r *frame_model.Room, player *frame_model.Player) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.MemberMap[player.ID] = player
	logger.InfoWF("player %s addPlayer in room %s", zap.Uint64("playerId", player.ID), zap.String("roomId", r.RoomIdStr))
	// 如果是第一个玩家，启动房间
	if len(r.MemberMap) == 1 && !r.Running {
		start(logger, r)
	}
}

func ExistRoom(logger fklog.FKLogI, userId uint64) (err error) {
	room, err := GetRoomByPlayerId(logger, userId)
	if err != nil {
		logger.ErrorWF("ExistRoom GetRoomByPlayerId failed", zap.Uint64("userId", userId))
		return errors.New("already exist room")
	}

	_, err = getRoomPlayer(logger, room, userId)
	if err != nil {
		logger.ErrorWF("ExistRoom getRoomPlayer failed", zap.Uint64("userId", userId))
		return errors.New("already exist room")
	}

	removePlayer(logger, room, userId)
	return nil
}

// 从房间移除玩家
func removePlayer(logger fklog.FKLogI, r *frame_model.Room, playerID uint64) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, exists := r.MemberMap[playerID]; exists {
		delete(r.MemberMap, playerID)
		delete(r.InputQueue, playerID) // 是否保留之前的同步信息，先按照不保留
		logger.InfoWF("player %s removePlayer in room %s", zap.Uint64("playerId", playerID), zap.String("roomId", r.RoomIdStr))
		// 如果房间空了，停止房间
		if len(r.MemberMap) == 0 && r.Running {
			stop(logger, r)
		}
	}
}

// 添加输入到队列
func addInputs(r *frame_model.Room, playerID uint64, inputs []frame_model.Input) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].Client <= inputs[j].Client
	})
	if _, exists := r.InputQueue[playerID]; !exists {
		r.InputQueue[playerID] = []frame_model.Input{}
	}
	r.InputQueue[playerID] = append(r.InputQueue[playerID], inputs...)
}

// 启动房间
func start(logger fklog.FKLogI, r *frame_model.Room) {
	if r.Running {
		return
	}
	r.Running = true
	r.Ticker = time.NewTicker(time.Duration(1000/r.GameTick) * time.Millisecond)
	go func() {
		for range r.Ticker.C {
			updateFrame(logger, r)
		}
	}()
	logger.InfoWF("room %s start success, tick: %d FPS", zap.String("roomId", r.RoomIdStr), zap.Any("tick", r.GameTick))
}

// 停止房间
func stop(logger fklog.FKLogI, r *frame_model.Room) {
	if !r.Running {
		return
	}
	r.Running = false
	if r.Ticker != nil {
		r.Ticker.Stop()
	}
	logger.InfoWF("room %s stop success", zap.String("roomId", r.RoomIdStr))
}

// 更新帧
func updateFrame(logger fklog.FKLogI, r *frame_model.Room) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// 增加帧号
	r.CurrentFrame++

	// 创建新的帧数据
	frameData := frame_model.NewFrameData(r.CurrentFrame)

	// 复制输入队列到帧数据
	for playerID, inputs := range r.InputQueue {
		if len(inputs) > 0 {
			frameData.Inputs[playerID] = make([]frame_model.Input, len(inputs))
			copy(frameData.Inputs[playerID], inputs)
			// 清空该玩家的输入队列
			r.InputQueue[playerID] = []frame_model.Input{}
		}
	}

	// 缓存每一帧数据
	r.FrameDataList = append(r.FrameDataList, frameData)

	idPack := &MazeRoom.MazeFrameSyncID{
		FrameId:    proto.Int32(int32(frameData.FrameID)),
		AccessInfo: proto.String(r.GameAccessInfo),
	}
	for _, inputList := range frameData.Inputs {
		for _, input := range inputList {
			idPack.Bytes = append(idPack.Bytes, input.Data)
		}
	}

	// 广播帧数据给所有玩家
	for _, player := range r.MemberMap {
		// push其他人
		_ = online.ClusterPush(context.TODO(), player.ID, 10544, idPack)
	}
}

func GetFrame(logger fklog.FKLogI, userId uint64, beginId, endId int32) (frameDataList []*frame_model.FrameData, hasMore bool, err error) {
	room, err := GetRoomByPlayerId(logger, userId)
	if err != nil {
		logger.ErrorWF("GetFrame getRoomByPlayerId invalid", zap.Error(err))
		return nil, false, errors.New("no found room")
	}

	queue := room.InputQueue
	if len(queue) == 0 {
		logger.ErrorWF("GetFrame no input queue in room", zap.Error(err))
		return nil, false, errors.New("no input queue in room")
	}

	if endId > int32(room.CurrentFrame) {
		logger.ErrorWF("GetFrame endId out of range", zap.Error(err))
		return nil, false, errors.New("endId out of range")
	}

	sort.Slice(room.FrameDataList, func(i, j int) bool {
		return room.FrameDataList[i].FrameID < room.FrameDataList[j].FrameID
	})

	list := room.FrameDataList[beginId:endId]
	if len(list) == 0 {
		logger.ErrorWF("GetFrame is nil", zap.Error(err))
		return nil, false, errors.New("getFrame is nil")
	}

	hasMore = endId < int32(room.CurrentFrame)
	return list, hasMore, nil
}
