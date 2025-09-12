/*
 * @Author: majian
 * @Date: 2025-01-03 19:41:51
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 21:39:45
 */
package attr_calc

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/vardef"
	"maze_game_server/servers/maze_main_server/process/attr_calc/mazeattrlogic"
)

func RunPreviewDac(ctx context.Context, userId uint64, dacParam *mazeattrlogic.DACParam,
	careAttrs []int32, replaceSrc int32,
	repalceAttrs map[int32]int64, preViewId int64) (result map[int32]int64, err error) {
	dacParam.IsPreview = true
	dac := mazeattrlogic.NewDAC(ctx, userId)
	err = dac.InitData(ctx, dacParam)
	if err != nil {
		return
	}

	pp := &mazeattrlogic.PreviewParam{}
	pp.BuffSrc = replaceSrc
	pp.CareAttrs = careAttrs
	pp.PreviewId = preViewId
	pp.RepalceAttrs = repalceAttrs
	dac.SetPreviewInfo(ctx, pp)
	dac.Prepare()

	err = dac.Calc(ctx)
	if err != nil {
		return
	}
	result = dac.GetCareAttrs(careAttrs)
	return
}

func DoEquipPreview(ctx context.Context, userId uint64, careAttrs []int32,
	replaceSrc int32,
	repalceAttrs map[int32]int64, previewId int64) (result map[int32]int64, err error) {
	p := mazeattrlogic.NewDACParam()
	p.ChgType = constdef.MazeBuffChgTypeEquipDress
	p.ChgDesc = vardef.MazeBuffChgTypeDesc[p.ChgType]
	return RunPreviewDac(ctx, userId, p, careAttrs, replaceSrc, repalceAttrs, previewId)
}

type PreviewAttrInfo struct {
	PreviewId    int64           // 预览Id
	CareAttrs    []int32         // 预览的属性
	BuffSrc      int32           // buff来源
	RepalceAttrs map[int32]int64 // 要替换的属性
	Result       map[int32]int64 // 预览的结果
}

type PreviewPairAttrInfo struct {
	PreviewId     int64           // 预览Id
	CareAttrs     []int32         // 预览的属性
	BuffSrc       int32           // buff来源
	RepalceAttrs1 map[int32]int64 // 要替换的属性1
	Result1       map[int32]int64 // 预览的结果1
	RepalceAttrs2 map[int32]int64 // 要替换的属性2
	Result2       map[int32]int64 // 预览的结果2
}

// 批量预览
func DoEquipBatchPreview(ctx context.Context, userId uint64, previewList []*PreviewAttrInfo) (preResult []*PreviewAttrInfo, err error) {
	p := mazeattrlogic.NewDACParam()
	p.ChgType = constdef.MazeBuffChgTypeEquipDress
	p.ChgDesc = vardef.MazeBuffChgTypeDesc[p.ChgType]
	p.IsPreview = true
	dac := mazeattrlogic.NewDAC(ctx, userId)
	err = dac.InitData(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, previewInfo := range previewList {
		cp := dac.CloneData(ctx)

		pp := &mazeattrlogic.PreviewParam{}
		pp.BuffSrc = previewInfo.BuffSrc
		pp.CareAttrs = previewInfo.CareAttrs
		pp.PreviewId = previewInfo.PreviewId
		pp.RepalceAttrs = previewInfo.RepalceAttrs
		cp.SetPreviewInfo(ctx, pp)

		cp.Prepare()
		err = dac.Calc(ctx)
		if err != nil {
			return nil, err
		}
		previewInfo.Result = cp.GetCareAttrs(previewInfo.CareAttrs)
	}
	return previewList, nil
}

// 成对批量预览
func DoEquipBatchPairPreview(ctx context.Context, userId uint64, previewList []*PreviewPairAttrInfo) (result []*PreviewPairAttrInfo, err error) {
	p := mazeattrlogic.NewDACParam()
	p.ChgType = constdef.MazeBuffChgTypeEquipDress
	p.ChgDesc = vardef.MazeBuffChgTypeDesc[p.ChgType]
	p.IsPreview = true
	dac := mazeattrlogic.NewDAC(ctx, userId)
	err = dac.InitData(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, previewInfo := range previewList {
		cp := dac.CloneData(ctx)

		pp := &mazeattrlogic.PreviewParam{}
		pp.BuffSrc = previewInfo.BuffSrc
		pp.CareAttrs = previewInfo.CareAttrs
		pp.PreviewId = previewInfo.PreviewId
		pp.RepalceAttrs = previewInfo.RepalceAttrs1
		cp.SetPreviewInfo(ctx, pp)

		cp.Prepare()
		err = cp.Calc(ctx)
		if err != nil {
			return nil, err
		}
		previewInfo.Result1 = cp.GetCareAttrs(previewInfo.CareAttrs)

		cp = dac.CloneData(ctx)

		pp = &mazeattrlogic.PreviewParam{}
		pp.BuffSrc = previewInfo.BuffSrc
		pp.CareAttrs = previewInfo.CareAttrs
		pp.PreviewId = previewInfo.PreviewId
		pp.RepalceAttrs = previewInfo.RepalceAttrs2
		cp.SetPreviewInfo(ctx, pp)

		cp.Prepare()
		err = cp.Calc(ctx)
		if err != nil {
			return nil, err
		}
		previewInfo.Result2 = cp.GetCareAttrs(previewInfo.CareAttrs)
	}
	return previewList, nil
}
