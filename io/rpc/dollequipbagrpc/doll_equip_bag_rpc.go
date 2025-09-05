package dollequipbagrpc

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/equip"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkrpc/thrift_rpc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var gRpcClient = thrift_rpc.AsyncRpc{}

func init() {
	//_ = fkconfig.RegisterNameNode("dollequipbagrpc", 19582, &gRpcClient)
}

// 添加装备
func MazeBagAddRQ(ctx context.Context, req *MazeEquipSvr.SvrAddMazeEquipRQ, res *MazeEquipSvr.SvrAddMazeEquipRS) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "MazeBagAddRQ start", zap.Any("req", req))
	err := equip.OnSvrAddMazeEquipRQ(ctx, int64(req.GetUserId()), req, res, "")
	if err != nil {
		logger.CtxError(ctx, "MazeBagAddRQ OnSvrAddMazeEquipRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.CtxError(ctx, "MazeBagAddRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
}

// 添加装备
func MazeBagAddRQWithOpData(ctx context.Context, req *MazeEquipSvr.SvrAddMazeEquipRQ, res *MazeEquipSvr.SvrAddMazeEquipRS, opData string) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "MazeBagAddRQWithOpData start", zap.Any("req", req))
	err := equip.OnSvrAddMazeEquipRQ(ctx, int64(req.GetUserId()), req, res, opData)
	if err != nil {
		logger.CtxError(ctx, "MazeBagAddRQWithOpData OnSvrAddMazeEquipRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.CtxError(ctx, "MazeBagAddRQWithOpData res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
}

// 更换装备rpc
func MazeEquipAssembleRQ(ctx context.Context, req *MazeEquipSvr.SvrMazeEquipAssembleRQ, res *MazeEquipSvr.SvrMazeEquipAssembleRS) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "MazeEquipAssembleRQ start", zap.Any("req", req))
	err := equip.OnSvrMazeEquipAssembleRQ(ctx, int64(req.GetUserId()), req, res)
	if err != nil {
		logger.CtxError(ctx, "MazeEquipAssembleRQ OnSvrMazeEquipAssembleRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.CtxError(ctx, "MazeEquipAssembleRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
	response, err := gRpcClient.DealTwowayMessage(131423, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131424)
	if err != nil {
		logger.CtxError(ctx, "DollEquipAssembleRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.CtxError(ctx, "DollEquipAssembleRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "DollEquipAssembleRQ recv rs",
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}

// 出售装备rpc
func MazeEquipSaleRQ(ctx context.Context, req *MazeEquipSvr.SvrMazeEquipSaleRQ, res *MazeEquipSvr.SvrMazeEquipSaleRS) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "MazeEquipSaleRQ start", zap.Any("req", req))
	err := equip.OnSvrDollEquipSaleRQ(context.TODO(), int64(req.GetUserId()), req, res)
	if err != nil {
		logger.CtxError(ctx, "MazeEquipSaleRQ OnSvrDollEquipSaleRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.CtxError(ctx, "MazeEquipSaleRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
	response, err := gRpcClient.DealTwowayMessage(131425, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131426)
	if err != nil {
		logger.CtxError(ctx, "DollEquipSaleRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.CtxError(ctx, "DollEquipSaleRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "DollEquipSaleRQ recv rs",
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}

// 实例化装备
func MazeBagInstanceRQ(ctx context.Context, req *MazeEquipSvr.SvrMazeEquipInstanceRQ, res *MazeEquipSvr.SvrMazeEquipInstanceRS) error {
	logger := fklog.ContextAppLogger(ctx)
	now := time.Now()
	response, err := gRpcClient.DealTwowayMessage(131431, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131432)
	if err != nil {
		logger.CtxError(ctx, "MazeBagInstanceRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.CtxError(ctx, "MazeBagInstanceRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "MazeBagInstanceRQ recv rs",
		zap.Duration("cost", time.Since(now)),
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}
