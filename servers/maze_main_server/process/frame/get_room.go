package frame

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/frame_model"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeRoom"
	frame_service "maze_game_server/services/frame"
)

// 获取房间信息
func (g *Frame) OnGetRoomRQ_10537_10538(s *session.Session, req *MazeRoom.MazeGetGameRoomInfoRQ) (err error) {
	defer fkprometheus.InfoPMT("OnGetRoomRQ")()
	ctx := s.Context()
	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeGetGameRoomInfoRS{}

	logger.InfoWF("OnGetRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnGetRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	room, err := frame_service.GetRoomByPlayerId(logger, userId)
	if err != nil {
		logger.ErrorWF("OnGetRoomRQ GetRoomByPlayerId fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	res.Data = roomPb(room)
	return nil
}

func roomPb(room *frame_model.Room) *MazeRoom.RoomInfo {
	info := &MazeRoom.RoomInfo{
		RoomIdStr:              proto.String(room.RoomIdStr),
		RoomState:              proto.Int32(room.RoomState),
		MaxMemberNum:           proto.Int32(room.MaxMemberNum),
		CreateTime:             proto.Int64(room.CreateTime),
		UpdateTimeStamp:        proto.Int64(room.UpdateTimeStamp),
		GameTick:               proto.Int32(room.GameTick),
		StartPercent:           proto.Int32(room.StartPercent),
		GameLastTime:           proto.Int32(room.GameLastTime),
		GameVersion:            proto.Int64(room.GameVersion),
		GameAccessInfo:         proto.String(room.GameAccessInfo),
		UdpReliabilityStrategy: proto.Int32(room.UdpReliabilityStrategy),
		RoomExtInfo:            proto.String(room.RoomExtInfo),
		Seed:                   proto.String(room.Seed),
	}
	for _, player := range room.MemberMap {
		info.MemberList = append(info.MemberList, &MazeRoom.RoomMember{
			ClientId: proto.Int32(int32(player.ID)),
			Role:     proto.Int32(player.Role),
		})
	}

	return info
}
