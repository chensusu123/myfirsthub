/*
 * @Author: majian
 * @Date: 2024-08-10 11:53:39
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 15:03:52
 */
package equip

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/protodef/KafkaMsgNotify"
	"go.uber.org/zap"
)

func OnKafkaTcpMsgRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnKafkaTcpMsgRQ")()
	req := request.(*KafkaMsgNotify.KafkaMsgDistributeRQ)
	res := response.(*KafkaMsgNotify.KafkaMsgDistributeRS)

	res.ErrInfo = errors.NO_ERROR
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnKafkaTcpMsgRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnKafkaTcpMsgRQ with", zap.Any("req", req))

	if req.GetUserId() == 0 {
		userCtx.WarnWF("OnKafkaTcpMsgRQ invald userId", zap.Any("msg", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的用户Id")
		return
	}
	if req.GetMsgType() <= 0 {
		userCtx.WarnWF("OnKafkaTcpMsgRQ invald msg type", zap.Any("msg", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的消息类型")
		return
	}
	if len(req.GetMsgData()) == 0 {
		userCtx.WarnWF("OnKafkaTcpMsgRQ invald msg content", zap.Any("msg", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的消息类型")
		return
	}

	if f, ok := GTcpMsgCallBackMap[req.GetMsgType()]; ok {
		f(userCtx, shardingID, req.GetMsgData())
		userCtx.InfoWF("OnKafkaTcpMsgRQ has reg func", zap.Int32("msgType", req.GetMsgType()))
	} else {
		userCtx.InfoWF("OnKafkaTcpMsgRQ no has reg func", zap.Int32("msgType", req.GetMsgType()))
	}
	return nil
}

type TcpMsgCallBackFunc func(logger fklog.FKLogI, userId uint64, msg []byte) error

var GTcpMsgCallBackMap = make(map[int32]TcpMsgCallBackFunc)

func RegTcpMsgCallBackFunc(msgType int32, f TcpMsgCallBackFunc) {
	if _, ok := GTcpMsgCallBackMap[msgType]; ok {
		fkfmt.Println("RegTcpMsgCallBackFunc already reg", zap.Int32("typ", msgType))
		return
	}
	GTcpMsgCallBackMap[msgType] = f
}
