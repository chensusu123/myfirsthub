/*
@Author: xiaobo
@Date: 2025/3/21 16:42
@Description:
*/

package item

import (
	"maze_game_server/common/itemclass/mazebag"
	"maze_game_server/common/itemclass/mazecommonvalue"
)

func init() {
	mazebag.Register(GloRegIns)

	mazecommonvalue.Register(GloRegIns)
}
