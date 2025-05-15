/*
 * @Author: majian
 * @Date: 2024-06-21 21:20:30
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 17:10:33
 */
package calcassembleattr

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeBuffData"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
)

// 打包普通属性
func PackCommonAttr(in *MazeEquipCache.EquipAttrInfo) *MazeBuffData.MazeBuffAttr {
	if in == nil {
		return nil
	}
	out := new(MazeBuffData.MazeBuffAttr)
	out.AttrId = proto.Int32(in.GetAttrId())
	out.AttrVal = proto.Int64(in.GetAttrValue())
	return out
}

// 打包附魔属性
// func PackTalentAttr(in *MazeEquipCache.EquipAttrInfo, enchantId int32) *MazeBuffData.MazeBuffAttr {
// 	if in == nil {
// 		return nil
// 	}
// 	out := new(MazeBuffData.MazeBuffAttr)
// 	out.AttrId = proto.Int32(in.GetAttrId())
// 	out.AttrVal = proto.Int64(in.GetAttrValue())
// 	out.AttrType = proto.Int32(int32(DollEquipAttr.ENUM_DOLL_EQUIP_ATTR_TYPE_ENCHANT))
// 	out.AttrRange = proto.Int32(enchantId) // 附魔类型
// 	return out
// }

// 追加属性
func AppendDollAttr(in *MazeBuffData.MazeBuffDb, item *MazeBuffData.MazeBuffAttr) {
	if item == nil {
		return
	}
	for _, attr := range in.MazeRealBuffs {
		if attr.GetAttrId() == item.GetAttrId() {
			attr.AttrVal = proto.Int64(item.GetAttrVal() + attr.GetAttrVal())
			return
		}
	}
	in.MazeRealBuffs = append(in.MazeRealBuffs, item)
}

// 打包属性
func PackMapAttr(in map[int32]int64) *MazeBuffData.MazeBuffDb {
	if in == nil {
		return nil
	}
	out := new(MazeBuffData.MazeBuffDb)
	for id, val := range in {
		out.MazeRealBuffs = append(out.MazeRealBuffs, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(val),
		})
	}
	return out
}

func PackMapAttrAll(real, show map[int32]int64) *MazeBuffData.MazeBuffDb {
	if real == nil && show == nil {
		return nil
	}
	out := new(MazeBuffData.MazeBuffDb)
	for id, val := range real {
		out.MazeRealBuffs = append(out.MazeRealBuffs, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(val),
		})
	}

	for id, val := range show {
		out.MazeShowBuffs = append(out.MazeShowBuffs, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(val),
		})
	}
	return out
}

func AppendDollAttrKv(in *MazeBuffData.MazeBuffDb, k int32, v int64) {
	for _, attr := range in.MazeRealBuffs {
		if attr.GetAttrId() == k {
			attr.AttrVal = proto.Int64(v + attr.GetAttrVal())
			return
		}
	}
	in.MazeRealBuffs = append(in.MazeRealBuffs, &MazeBuffData.MazeBuffAttr{AttrId: proto.Int32(k),
		AttrVal: proto.Int64(v)})
}

func AppendDollPanelAttrKv(in *MazeBuffData.MazeBuffDb, k int32, v int64) {
	for _, attr := range in.MazeShowBuffs {
		if attr.GetAttrId() == k {
			attr.AttrVal = proto.Int64(v + attr.GetAttrVal())
			return
		}
	}
	in.MazeShowBuffs = append(in.MazeShowBuffs, &MazeBuffData.MazeBuffAttr{AttrId: proto.Int32(k),
		AttrVal: proto.Int64(v)})
}

func AppendDollPanelAttr(in *MazeBuffData.MazeBuffDb, item *MazeBuffData.MazeBuffAttr) {
	if item == nil {
		return
	}
	for _, attr := range in.MazeShowBuffs {
		if attr.GetAttrId() == item.GetAttrId() {
			attr.AttrVal = proto.Int64(item.GetAttrVal() + attr.GetAttrVal())
			return
		}
	}
	in.MazeShowBuffs = append(in.MazeShowBuffs, item)
}
