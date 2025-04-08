/*
 * @Author: majian
 * @Date: 2025-03-22 11:06:41
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:30:51
 */
package mazeenergyrpc

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc/thrift_rpc"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergySvr"
	"go.uber.org/zap"
)

var gRpcClient = thrift_rpc.AsyncRpc{}

func init() {
	_ = fkconfig.RegisterNameNode("mazeenergyrpc", 19354, &gRpcClient)
}

// 扣除迷宫体力接口
// 体力不足，返回错误码(80000 // 体力不足),并带回剩余的体力值
// 扣体力成功，返回剩余的体力值
func SubMazeEnergyRQ(logger fklog.FKLogI, req *MazeEnergySvr.SubMazeEnergyRQ, res *MazeEnergySvr.SubMazeEnergyRS) error {
	defer fkprometheus.DebugPMT("SubMazeEnergyRQ")()
	response, err := gRpcClient.DealTwowayMessage(131445, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131446)
	if err != nil {
		logger.ErrorWF("SubMazeEnergyRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("SubMazeEnergyRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("SubMazeEnergyRQ recv rs",
		zap.Any("req", req),
		zap.Any("rs", res))
	return nil
}

// 加迷宫体力接口
// 加之前体力已满，返回错误码(80001 // 体力已满),并带回剩余的体力值
// 加体力成功，返回剩余的体力值
func AddMazeEnergyRQ(logger fklog.FKLogI, req *MazeEnergySvr.AddMazeEnergyRQ, res *MazeEnergySvr.AddMazeEnergyRS) error {
	defer fkprometheus.DebugPMT("AddMazeEnergyRQ")()
	response, err := gRpcClient.DealTwowayMessage(131455, req, req.GetUserId())
	err = gRpcClient.CheckReplyType(response, err, 131456)
	if err != nil {
		logger.ErrorWF("AddMazeEnergyRQ check res failed", zap.Any("req", req), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("AddMazeEnergyRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
		return err
	}
	logger.InfoWF("AddMazeEnergyRQ recv rs",
		zap.Any("req", req),
		zap.Any("rs", res))
	return nil
}
