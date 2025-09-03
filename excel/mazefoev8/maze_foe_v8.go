package mazefoev8

import (
	"context"
	"maze_game_server/config/GMazeFoeV8Cfg"
)

func GetMazeFoeConfig(ctx context.Context, monsterID int32) *GMazeFoeV8Cfg.MazeFoeV8ConfigRow {
	return GMazeFoeV8Cfg.GetWithCtx(ctx, monsterID)
}
