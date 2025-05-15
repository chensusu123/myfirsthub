package funcopencheck

import (
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeActionCountV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
)

// 检查结果
type FoResult struct {
	NeedLv   int32  // 需要人物等级
	LoseDesc string // 失败描述
	IsOpen   bool   // 是否可以开启
}

// 根据人物等级判断功能是否开启
func IsFuncOpen(funcId, mazeLevel int32) (result *FoResult, err error) {
	row := GMazeActionCountV8Cfg.Get(funcId)
	if row == nil {
		return nil, errors.New("no found cfg")
	}
	result = &FoResult{}
	// 7 表示人物等级类型
	result.NeedLv = row.Count_level_v8[10]
	result.LoseDesc = row.Lose_desc_v8
	result.IsOpen = mazeLevel >= result.NeedLv
	return result, nil
}
