/*
 * @Author: majian
 * @Date: 2025-01-10 22:34:55
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-11 19:54:02
 */
package attr_calc

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeAttrCalcSvr"
	"go.uber.org/zap"
)

func OnMazeAttrPairPreviewRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnMazeAttrPairPreviewRQ")()
	userCtx := fkserver.NewUserContext(ctx.Context, uint64(shardingID), ctx.FKLogI)
	req := rqMsg.(*MazeAttrCalcSvr.MazeAttrPairPreviewRQ)
	res := rsMsg.(*MazeAttrCalcSvr.MazeAttrPairPreviewRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.InfoWF("OnMazeAttrPairPreviewRQ with", zap.Any("rq", req))

	defer func() {
		userCtx.InfoWF("OnMazeAttrPairPreviewRQ end ", zap.Any("req", req), zap.Any("res", res))
	}()
	// if !BreedVersionFC.IsDollVersion(userCtx, uint64(shardingID)) {
	// 	userCtx.WarnWF("OnMazeAttrPairPreviewRQ not doll version", zap.Int64("userID", shardingID))
	// 	return
	// }
	if req.GetAttrSrc() == 0 || len(req.GetCareAttrs()) == 0 {
		userCtx.WarnWF("OnMazeAttrPairPreviewRQ param less", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("less param")
		return
	}
	var preList []*PreviewPairAttrInfo
	for _, p := range req.GetPreviewParams() {
		preData := &PreviewPairAttrInfo{}
		preData.BuffSrc = req.GetAttrSrc()
		preData.CareAttrs = req.GetCareAttrs()
		preData.PreviewId = p.GetPreviewId()
		preData.RepalceAttrs1 = make(map[int32]int64)
		preData.RepalceAttrs2 = make(map[int32]int64)
		for _, attr := range p.GetReplaceAttrs1() {
			preData.RepalceAttrs1[attr.GetAttrId()] = attr.GetAttrVal()
		}
		for _, attr := range p.GetReplaceAttrs2() {
			preData.RepalceAttrs2[attr.GetAttrId()] = attr.GetAttrVal()
		}
		preList = append(preList, preData)
	}
	if len(preList) > 0 {
		preList, err = DoEquipBatchPairPreview(userCtx, uint64(shardingID), preList)
		if err != nil {
			userCtx.ErrorWF("OnMazeAttrPairPreviewRQ DoEquipBatchPairPreview fail",
				zap.Any("preList", preList))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return
		}
		for _, preRs := range preList {
			rsPb := &MazeAttrCalcSvr.PairPreviewResult{}
			rsPb.PreviewId = proto.Int64(preRs.PreviewId)
			for k, v := range preRs.Result1 {
				rsPb.Result1 = append(rsPb.Result1,
					&MazeAttrCalcSvr.PreviewAttr{
						AttrId:  proto.Int32(k),
						AttrVal: proto.Int64(v)})
			}

			for k, v := range preRs.Result2 {
				rsPb.Result2 = append(rsPb.Result2,
					&MazeAttrCalcSvr.PreviewAttr{
						AttrId:  proto.Int32(k),
						AttrVal: proto.Int64(v)})
			}
			res.PreviewResult = append(res.PreviewResult, rsPb)
		}
	}
	return nil
}
