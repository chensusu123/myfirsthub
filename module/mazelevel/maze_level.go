/*
 * @Author: majian
 * @Date: 2025-04-17 11:02:31
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-17 16:05:00
 */
package mazelevel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
)

func GetMazelLevel(logger fklog.FKLogI, userId uint64) (level int64, err error) {
	level, err = mazeuserlevelredis.GetUserLevel(logger, userId)
	if err == nil && level == 0 {
		level = 1
	}
	return
}
