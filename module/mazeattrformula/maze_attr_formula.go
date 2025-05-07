/*
 * @Author: majian
 * @Date: 2024-11-30 16:34:31
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 15:15:42
 */
package mazeattrformula

import (
	"errors"
	"fmt"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttributeFormulaV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func CalcMazeFormulaAttr(logger fklog.FKLogI, formulaId int32,
	attrs map[int32]int64) (m int64, desc string, err error) {
	var p1, p2, p3, p4, p5, p11, p12 float64
	var p8, p9, p10 float64
	var p8s, p9s, p10s string

	cfg := GMazeAttributeFormulaV8Cfg.Get(formulaId)
	if cfg == nil {
		logger.ErrorWF("CalcDollFormulaAttr cannot find config", zap.Int32("formulaId", formulaId))
		return 0, desc, errors.New("no found formula config")
	}

	//先取出参数
	if attr, ok := attrs[cfg.Parameter_1]; ok {
		p1 = float64(attr)
		desc += fmt.Sprintf("p1_%d_%d|", cfg.Parameter_1, attr)
	}
	if attr, ok := attrs[cfg.Parameter_2]; ok {
		p2 = float64(attr)
		desc += fmt.Sprintf("p2_%d_%d|", cfg.Parameter_2, attr)
	}

	if attr, ok := attrs[cfg.Parameter_3]; ok {
		p3 = float64(attr)
		desc += fmt.Sprintf("p3_%d_%d|", cfg.Parameter_3, attr)
	}

	if attr, ok := attrs[cfg.Parameter_4]; ok {
		p4 = float64(attr)
		desc += fmt.Sprintf("p4_%d_%d|", cfg.Parameter_4, attr)
	}

	if attr, ok := attrs[cfg.Parameter_5]; ok {
		p5 = float64(attr)
		desc += fmt.Sprintf("p5_%d_%d|", cfg.Parameter_5, attr)
	}

	if attr, ok := attrs[cfg.Parameter_11]; ok {
		p11 = float64(attr)
		desc += fmt.Sprintf("p11_%d_%d|", cfg.Parameter_11, attr)
	}

	if attr, ok := attrs[cfg.Parameter_12]; ok {
		p12 = float64(attr)
		desc += fmt.Sprintf("p12_%d_%d|", cfg.Parameter_12, attr)
	}
	desc += "p8_"
	p8, p8s = CalcMultiplyParam(cfg.Parameter_8, attrs)
	desc += p8s
	desc += "p9_"
	p9, p9s = CalcMultiplyParam(cfg.Parameter_9, attrs)
	desc += p9s
	desc += "p10_"
	p10, p10s = CalcMultiplyParam(cfg.Parameter_10, attrs)
	desc += p10s

	//攻击、防御、生命值 = {[参数1 * （1+参数2/10000) * （1+参数8/10000） * （1+参数3/10000) + 参数4]* (1+参数5/10000)*(1+参数9/10000)*(1+参数10/10000) +参数11}*（1+参数12/10000）
	m1 := p1 * (1 + p2/10000.0) * p8 * (1 + p3/10000.0)
	m2 := (m1 + p4) * (1 + p5/10000.0) * p9 * p10
	m = int64((m2 + p11) * (1 + p12/10000.0))
	logger.InfoWF("CalcMazeFormulaAttr result",
		zap.Int32("formulaId", formulaId),
		zap.Int64("val", m),
		zap.Float64("m1", m1),
		zap.Float64("m2", m2),
		zap.String("desc", desc),
	)
	return m, desc, nil
}

func CalcMultiplyParam(params []int32, attrs map[int32]int64) (float64, string) {
	var f = 1.0
	var attrDesc string
	for _, attrId := range params {
		if attrId == 0 {
			continue
		}
		if val, ok := attrs[attrId]; ok {
			f = f * (1 + float64(val)/10000.0)
			attrDesc += fmt.Sprintf("%d_%d|", attrId, val)
		}
	}
	return f, attrDesc
}

// 获取攻防血相关的公式属性参数Id列表
func GetGFXFormulaParamAttrs(formulaId int32) []int32 {
	cfg := GMazeAttributeFormulaV8Cfg.Get(formulaId)
	if cfg == nil {
		return nil
	}
	var params []int32
	if cfg.Parameter_1 > 0 {
		params = append(params, cfg.Parameter_1)
	}
	if cfg.Parameter_2 > 0 {
		params = append(params, cfg.Parameter_2)
	}
	if cfg.Parameter_3 > 0 {
		params = append(params, cfg.Parameter_3)
	}
	if cfg.Parameter_4 > 0 {
		params = append(params, cfg.Parameter_4)
	}

	if cfg.Parameter_5 > 0 {
		params = append(params, cfg.Parameter_5)
	}
	if cfg.Parameter_11 > 0 {
		params = append(params, cfg.Parameter_11)
	}

	if cfg.Parameter_12 > 0 {
		params = append(params, cfg.Parameter_12)
	}
	for _, p := range cfg.Parameter_8 {
		if p > 0 {
			params = append(params, p)
		}
	}

	for _, p := range cfg.Parameter_9 {
		if p > 0 {
			params = append(params, p)
		}
	}

	for _, p := range cfg.Parameter_10 {
		if p > 0 {
			params = append(params, p)
		}
	}
	return params
}
