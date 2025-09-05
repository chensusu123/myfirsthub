package mazeskillattrcfg

import (
	"sync/atomic"
	"unsafe"

	"maze_game_server/config/GMazeSkillInfoV8Cfg"
)

type MazeSkillInfoV8ConfigEx struct {
	AttrSkillMap map[int32]*GMazeSkillInfoV8Cfg.MazeSkillInfoV8ConfigRow
}

func init() {
	GMazeSkillInfoV8Cfg.RegisterMazeSkillInfoV8InitCallBack("mazeskillattrcfg", loadDollAttrOrderByGroup)
}

var gConfigData *MazeSkillInfoV8ConfigEx

func loadDollAttrOrderByGroup(f *GMazeSkillInfoV8Cfg.MazeSkillInfoV8Config) {
	gTmp := &MazeSkillInfoV8ConfigEx{}
	gTmp.AttrSkillMap = make(map[int32]*GMazeSkillInfoV8Cfg.MazeSkillInfoV8ConfigRow)
	for _, row := range f.ConfigRows {
		gTmp.AttrSkillMap[row.Skill_attr_id] = row
	}
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(gTmp))
}

// 获取属性对应获得技能
func GetAttrSkill(attrID int32) *GMazeSkillInfoV8Cfg.MazeSkillInfoV8ConfigRow {
	return gConfigData.AttrSkillMap[attrID]
}
