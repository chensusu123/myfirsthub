/*
@Author: xiaobo
@Date: 2025/3/21 16:42
@Description:
*/

package item

import (
	"maze_game_server/services/bagservice"
	"maze_game_server/services/itemservice"
	"maze_game_server/services/moneyservice"
)

func init() {
	bagservice.Register(itemservice.GloRegIns)
	moneyservice.Register(itemservice.GloRegIns)
}
