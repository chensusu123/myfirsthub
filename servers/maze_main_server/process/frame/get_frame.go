package frame

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/frame_model"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeRoom"
	frame_service "maze_game_server/services/frame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 分片拉取房间信息RQ
func (g *Frame) OnGetFrameRQ_10533_10534(s *session.Session, req *MazeRoom.MazeGetFrameRQ) (err error) {
	defer fkprometheus.InfoPMT("OnGetFrameRQ")()
	ctx := s.Context()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeGetFrameRS{}

	logger.CtxInfo(ctx, "OnGetFrameRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnGetFrameRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnGetFrameRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	frameList, hasMore, err := frame_service.GetFrame(ctx, userId, req.GetBeginFrameId(), req.GetEndFrameId())
	if err != nil {
		return err
	}

	info := newFrameInfo(frameList)
	frameData := &MazeRoom.FrameData{}
	frameData.HasMore = proto.Bool(hasMore)
	frameData.FrameList = append(frameData.GetFrameList(), info)
	return nil
}

func newFrameInfo(arr []*frame_model.FrameData) *MazeRoom.FrameInfo {
	info := &MazeRoom.FrameInfo{
		FrameId: proto.Int64(0),
	}
	for _, input := range arr {
		data := newPkgInfo(input)
		info.PkgList = append(info.PkgList, data)
	}
	return info
}

func newPkgInfo(data *frame_model.FrameData) *MazeRoom.PkgInfo {
	info := &MazeRoom.PkgInfo{}
	for userId, inputList := range data.Inputs {
		info.OpenId = proto.Uint64(userId)
		for _, input := range inputList {
			info.ActionList = append(info.ActionList, input.Data)
		}
	}
	info.B = proto.Bool(false)

	return info
}
