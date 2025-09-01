/*
 * @Author: majian
 * @Date: 2025-03-22 11:06:41
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:30:51
 */
package mazeenergyrpc

// import (
// 	"google.golang.org/protobuf/proto"
// 	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
// 	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
// 	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
// 	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkrpc/thrift_rpc"
// 	"maze_game_server/pb/common/MazeEnergySvr"
// 	"go.uber.org/zap"
// )

// var gRpcClient = thrift_rpc.AsyncRpc{}

// func init() {
// 	_ = fkconfig.RegisterNameNode("mazeenergyrpc", 19354, &gRpcClient)
// }

// // 扣除迷宫体力接口
// // 体力不足，返回错误码(80000 // 体力不足),并带回剩余的体力值
// // 扣体力成功，返回剩余的体力值
// func SubMazeEnergyRQ(ctx context.Context, req *MazeEnergySvr.SubMazeEnergyRQ, res *MazeEnergySvr.SubMazeEnergyRS) error {
// 	defer fkprometheus.DebugPMT("SubMazeEnergyRQ")()
// 	response, err := gRpcClient.DealTwowayMessage(131445, req, req.GetUserId())
// 	err = gRpcClient.CheckReplyType(response, err, 131446)
// 	if err != nil {
// 		logger.CtxError(ctx,"SubMazeEnergyRQ check res failed", zap.Any("req", req), zap.Error(err))
// 		return err
// 	}

// 	err = proto.Unmarshal(response.GetContent(), res)
// 	if err != nil {
// 		logger.CtxError(ctx,"SubMazeEnergyRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
// 		return err
// 	}
// 	logger.CtxInfo(ctx,"SubMazeEnergyRQ recv rs",
// 		zap.Any("req", req),
// 		zap.Any("rs", res))
// 	return nil
// }

// // 加迷宫体力接口
// // 加之前体力已满，返回错误码(80001 // 体力已满),并带回剩余的体力值
// // 加体力成功，返回剩余的体力值
// func AddMazeEnergyRQ(ctx context.Context, req *MazeEnergySvr.AddMazeEnergyRQ, res *MazeEnergySvr.AddMazeEnergyRS) error {
// 	defer fkprometheus.DebugPMT("AddMazeEnergyRQ")()
// 	response, err := gRpcClient.DealTwowayMessage(131455, req, req.GetUserId())
// 	err = gRpcClient.CheckReplyType(response, err, 131456)
// 	if err != nil {
// 		logger.CtxError(ctx,"AddMazeEnergyRQ check res failed", zap.Any("req", req), zap.Error(err))
// 		return err
// 	}

// 	err = proto.Unmarshal(response.GetContent(), res)
// 	if err != nil {
// 		logger.CtxError(ctx,"AddMazeEnergyRQ unmarshal failed", zap.Any("req", req), zap.Error(err))
// 		return err
// 	}
// 	logger.CtxInfo(ctx,"AddMazeEnergyRQ recv rs",
// 		zap.Any("req", req),
// 		zap.Any("rs", res))
// 	return nil
// }
