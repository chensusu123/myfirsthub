/*
 * @Author: majian
 * @Date: 2024-12-21 16:33:52
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-11 12:13:46
 */
package commonlogic

import (
	"bytes"
	"fmt"
)

type AttrWeight struct {
	SourceId int32 `json:"source_id"` // 来源
	Weight   int64 `json:"weight"`    // 分量
}

type AttrWeightMap struct {
	AwMap map[int32]*AttrWeight
}

type AttrCalcRecorder struct {
	AttrRecorderMap map[int32]*AttrWeightMap // 属性ID:所有来源分量
}

func NewAttrCalcRecorder() *AttrCalcRecorder {
	obj := &AttrCalcRecorder{}
	obj.AttrRecorderMap = make(map[int32]*AttrWeightMap)
	return obj
}

func NewAttrWeightMap() *AttrWeightMap {
	obj := &AttrWeightMap{}
	obj.AwMap = make(map[int32]*AttrWeight)
	return obj
}

func (m *AttrCalcRecorder) AddWeight(attrId, sourceId int32, val int64) {
	if wm, ok := m.AttrRecorderMap[attrId]; ok {
		wm.AwMap[sourceId] = &AttrWeight{Weight: val, SourceId: sourceId}
	} else {
		wmNew := NewAttrWeightMap()
		wmNew.AwMap[sourceId] = &AttrWeight{Weight: val, SourceId: sourceId}
		m.AttrRecorderMap[attrId] = wmNew
	}
}

func (m *AttrWeightMap) DumpWeightInfo() string {
	var bs bytes.Buffer
	for _, weight := range m.AwMap {
		bs.WriteString(fmt.Sprintf("%s(%d):%d|", GetSrcName(int(weight.SourceId)), weight.SourceId, weight.Weight))
	}
	return bs.String()
}

func (m *AttrCalcRecorder) DumpWeightInfo(attrId int32) string {
	if wm, ok := m.AttrRecorderMap[attrId]; ok {
		return wm.DumpWeightInfo()
	}
	return ""
}
