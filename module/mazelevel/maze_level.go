/*
 * @Author: majian
 * @Date: 2025-04-17 11:02:31
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-17 16:05:00
 */
package mazelevel

import (
	"context"
	"maze_game_server/io/redis/mazeuserlevelredis"
)

func GetMazelLevel(ctx context.Context, userId uint64) (level int64, err error) {
	level, err = mazeuserlevelredis.GetUserLevel(ctx, userId)
	if err == nil && level == 0 {
		level = 1
	}
	return
}
