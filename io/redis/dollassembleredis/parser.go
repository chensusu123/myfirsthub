/*
 * @Author: majian
 * @Date: 2024-08-14 14:11:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 17:10:29
 */
package dollassembleredis

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
)

func unpackFieldToPb(field string, in []byte, pb *MazeEquipCache.MazeAssembleDb) error {
	if in == nil {
		return nil
	}
	var find = true
	switch field {
	case constdef.AssemblePrefixCurAssembleSuitIndex: // 当前生效套装Id
		ret, e := fkutil.Bytes2Int64(in)
		if e != nil {
			return e
		}
		pb.CurSuitIndex = proto.Int32(int32(ret))
	case constdef.AssemblePrefixSwitchSuitTime: // 当前套装切换时间
		ret, e := fkutil.Bytes2Int64(in)
		if e != nil {
			return e
		}
		pb.SwitchSuitTime = proto.Int64(ret)
	default:
		find = false
	}
	if find {
		return nil
	}

	pos := assemble.DecodeAssemblePosField(field)
	if pos > 0 { // 解析装备位信息
		slotInfo := &MazeEquipCache.MazeEquipSlotDb{}
		e := proto.Unmarshal(in, slotInfo)
		if e != nil {
			return e
		}
		posInfo := &MazeEquipCache.MazeEquipPosInfo{}
		posInfo.EquipPos = slotInfo
		pb.MazeEquips = append(pb.MazeEquips, posInfo)
	}
	return nil
}

func unpackFieldsToPb(rs map[string][]byte, pb *MazeEquipCache.MazeAssembleDb) (err error) {
	for fieldName, val := range rs {
		err = unpackFieldToPb(fieldName, val, pb)
		if err != nil {
			return
		}
	}
	return
}

func packFieldFromPb(field string, pb *MazeEquipCache.MazeAssembleDb) (out interface{}, err error) {
	var find = true
	switch field {
	case constdef.AssemblePrefixCurAssembleSuitIndex: // 当前生效套装
		out = pb.GetCurSuitIndex()
	case constdef.AssemblePrefixSwitchSuitTime: // 套装切换时间
		out = pb.GetSwitchSuitTime()
	default:
		find = false
	}
	if find {
		return
	}
	// 解析装备位
	pos := assemble.DecodeAssemblePosField(field)
	if pos > 0 {
		posList := pb.GetMazeEquips()
		for _, posInfo := range posList {
			if posInfo.GetEquipPos() != nil {
				if posInfo.GetEquipPos().GetPos() == pos {
					r, e := proto.Marshal(posInfo.GetEquipPos())
					if e != nil {
						err = e
						return
					}
					out = r
					return
				}
			}
		}
	}
	err = errors.New("no found parser")
	return
}

func packFieldsFromPb(pb *MazeEquipCache.MazeAssembleDb, fields []string) (rs map[string]interface{}, err error) {
	rs = make(map[string]interface{})
	for _, field := range fields {
		ret, err1 := packFieldFromPb(field, pb)
		if err1 != nil {
			err = err1
			return
		}
		rs[field] = ret
	}
	return
}
