package item

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeBag"
	"maze_game_server/services/bagservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

type Item struct {
	component.Base
}

func NewItem() *Item {
	return &Item{}
}

// 获取迷宫背包列表
func (i *Item) OnMazeBagListRQ_10400_10401(s *session.Session, req *MazeBag.MazeBagListRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "OnMazeBagListRQ start", zap.Any("req", req))
	res := &MazeBag.MazeBagListRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	uid := uint64(s.UID())
	defer fkprometheus.DebugPMT("OnMazeBagListRQ")()
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeBagListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	bagItemMap, err := bagservice.GlobalBagService.GetAllBagItem(ctx, uid)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取数据失败")
		logger.CtxError(ctx, "OnMazeBagListRQ GetAllBagItem err", zap.Error(err))
		return
	}
	res.Items = make([]*MazeBag.MazeBagItem, 0, len(bagItemMap))
	for id, count := range bagItemMap {
		if id <= 0 || count <= 0 {
			continue
		}
		res.Items = append(res.Items, itemutil.BuildMazeBagItem(ctx, id, count))
	}
	return
}

// 重置迷宫背包列表
func (i *Item) OnResetMazeBagRQ_10402_10403(s *session.Session, req *MazeBag.ResetMazeBagRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "OnResetMazeBagRQ start", zap.Any("req", req))

	res := &MazeBag.ResetMazeBagRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	uid := uint64(s.UID())

	defer fkprometheus.DebugPMT("OnResetMazeBagRQ")()
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnResetMazeBagRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	err = bagservice.GlobalBagService.DelAllBagItem(ctx, uid)
	if err != nil {
		logger.CtxError(ctx, "OnResetMazeBagRQ DelAllBagItem err", zap.Error(err))
		return
	}

	return
}
