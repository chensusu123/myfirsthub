/*
 * @Author: majian
 * @Date: 2025-01-10 22:15:04
 * @Last Modified by: majian
 * @Last Modified time: 2025-01-10 22:34:58
 */
package attr_calc

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/protodef/MazeAttrCalcSvr"
	"go.uber.org/zap"
)

func OnMazeAttrPreviewRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnMazeAttrPreviewRQ")()
	userCtx := fkserver.NewUserContext(ctx.Context, uint64(shardingID), ctx.FKLogI)
	req := rqMsg.(*MazeAttrCalcSvr.MazeAttrPreviewRQ)
	res := rsMsg.(*MazeAttrCalcSvr.MazeAttrPreviewRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.InfoWF("OnMazeAttrPreviewRQ with", zap.Any("rq", req))

	defer func() {
		userCtx.InfoWF("OnMazeAttrPreviewRQ end ", zap.Any("req", req), zap.Any("res", res))
	}()

	if req.GetAttrSrc() == 0 || len(req.GetCareAttrs()) == 0 {
		userCtx.WarnWF("OnMazeAttrPreviewRQ param less", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("less param")
		return
	}
	var preList []*PreviewAttrInfo
	for _, p := range req.GetPreviewParams() {
		preData := &PreviewAttrInfo{}
		preData.BuffSrc = req.GetAttrSrc()
		preData.CareAttrs = req.GetCareAttrs()
		preData.PreviewId = p.GetPreviewId()
		preData.RepalceAttrs = make(map[int32]int64)
		for _, attr := range p.GetReplaceAttrs() {
			preData.RepalceAttrs[attr.GetAttrId()] = attr.GetAttrVal()
		}
		preList = append(preList, preData)
	}
	if len(preList) > 0 {
		preList, err = DoEquipBatchPreview(userCtx, uint64(shardingID), preList)
		if err != nil {
			userCtx.ErrorWF("OnMazeAttrPreviewRQ DoEquipBatchPreview fail",
				zap.Any("preList", preList))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return
		}
		for _, preRs := range preList {
			rsPb := &MazeAttrCalcSvr.PreviewResult{}
			rsPb.PreviewId = proto.Int64(preRs.PreviewId)
			for k, v := range preRs.Result {
				rsPb.Result = append(rsPb.Result,
					&MazeAttrCalcSvr.PreviewAttr{
						AttrId:  proto.Int32(k),
						AttrVal: proto.Int64(v)})
			}
			res.PreviewResult = append(res.PreviewResult, rsPb)
		}
	}
	return nil
}
