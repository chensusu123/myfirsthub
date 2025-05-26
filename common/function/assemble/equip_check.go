package assemble

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"maze_game_server/common/constdef"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/pb/server/MazeEquipCache"
)

// 是否有效的装备位
func IsValidEquipPos(pos uint32) bool {
	row := GMazeEquipPosRankV8Cfg.GetMazeEquipPosRankV8Config(int32(pos))
	return row != nil
}

// 编码装配field
func EnCodeAssembleEquipField(suitIndex, pos int32) string {
	return fmt.Sprintf("%s%d", constdef.AssemblePrefixEquip, suitIndex*100+pos)
}

// 解码装配field
func DecodeAssembleEquipField(field string) (suitIndex, posId int32) {
	if pos := strings.Index(field, constdef.AssemblePrefixEquip); pos == 0 {
		xy := fkutil.ToInt32(field[len(constdef.AssemblePrefixEquip):])
		suitIndex = xy / 100
		posId = xy % 100
	}
	return
}

// 是否有装配装备
func IsAssembleEquip(equip *MazeEquipCache.MazeEquipPosInfo) bool {
	if equip.GetEquipLoadInfo() != nil && equip.GetEquipLoadInfo().GetEquipGuid() > 0 {
		return true
	}
	return false
}

// 编码装配位field
func EnCodeAssemblePosField(pos int32) string {
	return fmt.Sprintf("%s%d", constdef.AssemblePrefixEquipPos, pos)
}

// 解码装配位field
func DecodeAssemblePosField(field string) (posId int32) {
	if pos := strings.Index(field, constdef.AssemblePrefixEquipPos); pos == 0 {
		posId = fkutil.ToInt32(field[len(constdef.AssemblePrefixEquipPos):])
	}
	return
}

func CloneAssemblePos(in *MazeEquipCache.MazeEquipPosInfo) *MazeEquipCache.MazeEquipPosInfo {
	out := &MazeEquipCache.MazeEquipPosInfo{}
	out.EquipInfo = in.EquipInfo // 装备信息不修改
	out.EquipLoadInfo = proto.Clone(in.GetEquipLoadInfo()).(*MazeEquipCache.MazeEquipPosDb)
	out.EquipPos = in.EquipPos // 装备位不会修改
	// 武力值不拷贝
	return out
}

func BatchCloneAssemblePos(in []*MazeEquipCache.MazeEquipPosInfo) []*MazeEquipCache.MazeEquipPosInfo {
	out := make([]*MazeEquipCache.MazeEquipPosInfo, 0, len(in))
	for _, pos := range in {
		out = append(out, CloneAssemblePos(pos))
	}
	return out
}

func IsFiveElemActivate(in int32) bool {
	return in&constdef.FiveElemActivate > 0
}

func IsHurtSuitActivate(in int32) bool {
	return in&constdef.HurtSuitActivate > 0
}
