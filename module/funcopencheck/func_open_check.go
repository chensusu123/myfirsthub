package funcopencheck

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeActionCountV8Cfg"
)

// 检查结果
type FoResult struct {
	NeedLv   int32  // 需要人物等级
	LoseDesc string // 失败描述
	IsOpen   bool   // 是否可以开启
}

// 根据人物等级判断功能是否开启
func IsFuncOpen(ctx context.Context, funcId, mazeLevel int32) (result *FoResult, err error) {
	row := GMazeActionCountV8Cfg.GetWithCtx(ctx, funcId)
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
