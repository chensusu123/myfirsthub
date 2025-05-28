package dollequipbagrpc

import (
	"time"

	"maze_game_server/pb/server/MazeEquipSvr"

	"maze_game_server/common/errors"
	"maze_game_server/servers/maze_main_server/process/equip"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkrpc/thrift_rpc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var gRpcClient = thrift_rpc.AsyncRpc{}

func init() {
	// _ = fkconfig.RegisterNameNode("dollequipbagrpc", 19582, &gRpcClient)
}

// 添加装备
func MazeBagAddRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrAddMazeEquipRQ, res *MazeEquipSvr.SvrAddMazeEquipRS) error {
	logger.InfoWF("MazeBagAddRQ start", zap.Any("req", req))
	err := equip.OnSvrAddMazeEquipRQ(logger, int64(req.GetUserId()), req, res, "")
	if err != nil {
		logger.ErrorWF("MazeBagAddRQ OnSvrAddMazeEquipRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("MazeBagAddRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
}

// 添加装备
func MazeBagAddRQWithOpData(logger fklog.FKLogI, req *MazeEquipSvr.SvrAddMazeEquipRQ, res *MazeEquipSvr.SvrAddMazeEquipRS, opData string) error {
	logger.InfoWF("MazeBagAddRQWithOpData start", zap.Any("req", req))
	err := equip.OnSvrAddMazeEquipRQ(logger, int64(req.GetUserId()), req, res, opData)
	if err != nil {
		logger.ErrorWF("MazeBagAddRQWithOpData OnSvrAddMazeEquipRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("MazeBagAddRQWithOpData res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
}

// 更换装备rpc
func MazeEquipAssembleRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrMazeEquipAssembleRQ, res *MazeEquipSvr.SvrMazeEquipAssembleRS) error {
	logger.InfoWF("MazeEquipAssembleRQ start", zap.Any("req", req))
	err := equip.OnSvrMazeEquipAssembleRQ(logger, int64(req.GetUserId()), req, res)
	if err != nil {
		logger.ErrorWF("MazeEquipAssembleRQ OnSvrMazeEquipAssembleRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("MazeEquipAssembleRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
	response, err := gRpcClient.DealTwowayMessage(131423, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131424)
	if err != nil {
		logger.ErrorWF("DollEquipAssembleRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("DollEquipAssembleRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("DollEquipAssembleRQ recv rs",
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}

// 出售装备rpc
func MazeEquipSaleRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrMazeEquipSaleRQ, res *MazeEquipSvr.SvrMazeEquipSaleRS) error {
	logger.InfoWF("MazeEquipSaleRQ start", zap.Any("req", req))
	err := equip.OnSvrDollEquipSaleRQ(logger, int64(req.GetUserId()), req, res)
	if err != nil {
		logger.ErrorWF("MazeEquipSaleRQ OnSvrDollEquipSaleRQ failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	if res.GetErrInfo() != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("MazeEquipSaleRQ res failed", zap.Any("req", req), zap.Error(err))
	}
	return err
	response, err := gRpcClient.DealTwowayMessage(131425, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131426)
	if err != nil {
		logger.ErrorWF("DollEquipSaleRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("DollEquipSaleRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("DollEquipSaleRQ recv rs",
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}

// 实例化装备
func MazeBagInstanceRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrMazeEquipInstanceRQ, res *MazeEquipSvr.SvrMazeEquipInstanceRS) error {
	now := time.Now()
	response, err := gRpcClient.DealTwowayMessage(131431, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131432)
	if err != nil {
		logger.ErrorWF("MazeBagInstanceRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("MazeBagInstanceRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("MazeBagInstanceRQ recv rs",
		zap.Duration("cost", time.Since(now)),
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}
