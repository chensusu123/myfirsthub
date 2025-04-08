package dollequipbagrpc

import (
	"time"

	"gitlab.ifreetalk.com/plate/protodef/MazeEquipSvr"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc/thrift_rpc"
	"go.uber.org/zap"
)

var gRpcClient = thrift_rpc.AsyncRpc{}

func init() {
	_ = fkconfig.RegisterNameNode("dollequipbagrpc", 19582, &gRpcClient)
}

// 添加装备
func MazeBagAddRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrAddMazeEquipRQ, res *MazeEquipSvr.SvrAddMazeEquipRS) error {
	now := time.Now()
	response, err := gRpcClient.DealTwowayMessage(131421, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131422)
	if err != nil {
		logger.ErrorWF("MazeBagAddRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("MazeBagAddRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("MazeBagAddRQ recv rs",
		zap.Duration("cost", time.Since(now)),
		zap.Any("req", req), zap.Any("rs", res))
	return nil
}

// 更换装备rpc
func MazeEquipAssembleRQ(logger fklog.FKLogI, req *MazeEquipSvr.SvrMazeEquipAssembleRQ, res *MazeEquipSvr.SvrMazeEquipAssembleRS) error {
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
