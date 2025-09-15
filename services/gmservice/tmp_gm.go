package gmservice

import (
	"fmt"
	"maze_game_server/model/friendmodel"
	"net/http"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

func (s *service) SetRecommendSize(writer http.ResponseWriter, request *http.Request) {
	// 外网线上环境不允许使用GM
	request.ParseForm()

	size := fkutil.ToUint64(request.Form.Get("size"))
	friendmodel.RecommendSize = int32(size)
	writer.Write([]byte(fmt.Sprintf("操作成功 当前推荐数量为%d", friendmodel.RecommendSize)))
}
