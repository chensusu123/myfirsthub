/*
 * @Author: majian
 * @Date: 2024-12-13 11:26:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 17:51:17
 */
package mazebuffchgrrecordapi

import (
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"maze_game_server/common/function/maputil"
	"maze_game_server/common/structsdef"
	"maze_game_server/common/vardef"
	"maze_game_server/io/kafka/mazebuffchgrecord"
	"maze_game_server/pb/server/MazeBuffData"
)

func SendMazeBuffChgRecord(logger fklog.FKLogI, userId uint64, src, chgReason int32, attrDbOld, attrDbNew *MazeBuffData.MazeBuffDb) error {
	r := structsdef.MazeGameBuffAttrChgRecord{}
	r.UserId = userId
	r.ChgDesc = vardef.MazeBuffChgTypeDesc[chgReason]
	r.ChgReason = chgReason
	r.SrcType = src
	if attrDbOld != nil {
		r.OldVal = maputil.MapToString(DollAttrDbToMap(attrDbOld))
	}
	if attrDbNew != nil {
		r.NewVal = maputil.MapToString(DollAttrDbToMap(attrDbNew))
	}

	return mazebuffchgrecord.SendMazeBuffAttrRecord(logger, &r)
}

func DollAttrDbToMap(attrDB *MazeBuffData.MazeBuffDb) map[int32]int64 {
	result := make(map[int32]int64)
	cares := ParseAttrToMap(CareAttrIDs)
	for _, attr := range attrDB.GetMazeRealBuffs() {
		if attr.GetAttrId() == 0 {
			continue
		}
		if _, ok := cares[attr.GetAttrId()]; ok {
			result[attr.GetAttrId()] = attr.GetAttrVal()
		}
	}

	for _, attr := range attrDB.GetMazeShowBuffs() {
		if attr.GetAttrId() == 0 {
			continue
		}
		if _, ok := cares[attr.GetAttrId()]; ok {
			result[attr.GetAttrId()] = attr.GetAttrVal()
		}
	}
	return result
}

var CareAttrIDs string // 需要关心的属性Id
func init() {
	param.StringP(&CareAttrIDs, "care:attr:chg:list", "19,21,23", "需要关心的属性变化列表")
}

func ParseAttrToMap(careIds string) map[int32]struct{} {
	ids := strings.Split(CareAttrIDs, ",")
	ret := make(map[int32]struct{})
	for _, id := range ids {
		ret[fkutil.ToInt32(id)] = struct{}{}
	}
	return ret
}
