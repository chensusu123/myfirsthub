package setseataskrpc

import (
	un_rpc_pack_type "gitlab.ifreetalk.com/plate/definition/uncgkconst"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"

	"go.uber.org/zap"
)

var setSeaTaskRpc = fkrpc.ThriftRPCClient{}

func init() {
	fkconfig.RegisterNameNode("SetSeaTaskRpc", 17762, &setSeaTaskRpc)
}

// SetSeaTask 设置定时器
func SetSeaTask(logger fklog.FKLogI, uid, session uint64, op SeaTaskSvr.TaskNotifyType, info *SeaTaskSvr.TaskInfo) (err error) {
	defer func() {
		logger.InfoWF("SetSeaTask", zap.Uint64("uid", uid), zap.Uint64("session", session),
			zap.String("op", op.String()), zap.String("info", info.String()), zap.Error(err))
	}()
	req := &SeaTaskSvr.SetSeaTaskRQ{}
	req.OpType = proto.Uint32(uint32(op))
	req.SessionId = &session
	req.TaskInfo = info

	response, err1 := setSeaTaskRpc.DealTwowayMessage(un_rpc_pack_type.UN_RPC_PACK_SET_SEA_TASK_RQ, req, uid)
	err = setSeaTaskRpc.CheckReplyType(response, err1, un_rpc_pack_type.UN_RPC_PACK_SET_SEA_TASK_RS)
	if err != nil {
		logger.ErrorWF("check SetSeaTask response failed, error", zap.Uint64("uid", uid), zap.Uint64("session", session),
			zap.String("op", op.String()), zap.String("info", info.String()), zap.Error(err))
		return err
	}
	res := &SeaTaskSvr.SetSeaTaskRS{}
	err = proto.Unmarshal(response.GetContent(), res)
	if err != nil {
		logger.ErrorWF("unmarshal SetSeaTask failed, error", zap.Uint64("uid", uid), zap.Uint64("session", session),
			zap.String("op", op.String()), zap.String("info", info.String()), zap.Error(err))
		return err
	}

	// 错误检查
	if res.ErrInfo != nil && res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		logger.ErrorWF("request failed.", zap.Uint64("uid", uid), zap.Uint64("session", session),
			zap.String("op", op.String()), zap.String("info", info.String()), zap.String("res", res.String()))
		return errors.New(res.GetErrInfo().String())
	}

	if res.GetSessionId() != req.GetSessionId() {
		logger.ErrorWF("session not match.", zap.Uint64("uid", uid), zap.Uint64("session", session),
			zap.String("op", op.String()), zap.String("info", info.String()), zap.String("res", res.String()))
		return errors.New("session not match")
	}

	return nil
}
