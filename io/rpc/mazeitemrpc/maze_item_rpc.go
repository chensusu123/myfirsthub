package mazeitemrpc

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc/stru"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc/thrift_rpc"
	"gitlab.ifreetalk.com/plate/protodef/MazeItemSvr"
	"go.uber.org/zap"
)

var itemRpc = thrift_rpc.AsyncRpc{}

func init() {
	// 迷宫物品管理rpc代理 20010
	fkconfig.RegisterNameNode("mazeitemrpcex", 20010, &itemRpc)
	// 17565
	// fkconfig.RegisterNode(un_cgk_svr_type.UN_CGK_SVR_TYPE_ITEM_SERVER, &itemRpc)
}

func DeductItemsRQ(logger fklog.FKLogI, req *MazeItemSvr.ConsumeItemRQ, res *MazeItemSvr.ConsumeItemRS) (err error) {
	logger.DebugWF("DeductItemsRQ send rq", zap.Any("rq", req))
	response, err := itemRpc.DealTwowayMessage(131439, req, req.GetUserId())
	err = CheckReplyType(response, err, 131440)
	if err != nil {
		logger.ErrorWF("DeductItemsRQ check res failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return err
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("DeductItemsRQ unmarshal failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return
	}
	logger.InfoWF("DeductItemsRQ recv rs", zap.Any("rs", res), zap.Any("req", req))
	return nil
}

func AddItemsRQ(logger fklog.FKLogI, req *MazeItemSvr.AddItemRQ, res *MazeItemSvr.AddItemRS) (err error) {
	logger.DebugWF("AddItemsRQ send rq", zap.String("rq", req.String()))
	response, err := itemRpc.DealTwowayMessage(131437, req, req.GetUserId())
	err = CheckReplyType(response, err, 131438)
	if err != nil {
		logger.ErrorWF("AddItemsRQ check res failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("[AddItemsRQ] unmarshal failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return
	}
	logger.InfoWF("AddItemsRQ recv rs", zap.String("rs", res.String()), zap.Any("req", req))
	return
}

func QueryItemsRQ(logger fklog.FKLogI, req *MazeItemSvr.QueryItemRQ, res *MazeItemSvr.QueryItemRS) (err error) {
	logger.DebugWF("QueryItemsRQ send rq", zap.String("rq", req.String()))
	response, err := itemRpc.DealTwowayMessage(131441, req, req.GetUserId())
	err = CheckReplyType(response, err, 131442)
	if err != nil {
		logger.ErrorWF("QueryItemsRQ check res failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return
	}

	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("[QueryItemsRQ] unmarshal failed", zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		return
	}
	logger.InfoWF("QueryItemsRQ recv rs", zap.String("rs", res.String()), zap.Any("req", req))
	return
}

func CheckReplyType(reply *stru.RPCResponse, err error, typeID int32) error {
	if err != nil {
		return err
	}
	if reply == nil {
		return stru.RPC_RESULT_NOT_SUCCESS
	}
	if typeID != reply.Type {
		return stru.RPC_RESULT_TYPE_NOT_MATCH
	}
	return nil
}
