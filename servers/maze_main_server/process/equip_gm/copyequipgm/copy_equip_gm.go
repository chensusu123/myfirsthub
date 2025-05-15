/*
 * @Author: majian
 * @Date: 2024-12-16 14:27:45
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-30 15:18:26
 */
package copyequipgm

import (
	"net/http"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm/copyusers"
)

func RegGm(logger fklog.FKLogI) {
	gm.SafeHttpRegister(logger, "/CopyEquipData", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		uidFile := request.Form.Get("file")
		if uidFile == "" {
			_, _ = writer.Write([]byte("请指定文件名称"))
			return
		}
		sp := request.Form.Get("sp")
		copyUser := copyusers.NewCopyUsers()
		err := copyUser.LoadUser(logger, uidFile, sp)
		if err != nil {
			fkfmt.Println("load user error", err, "path", uidFile)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		fkfmt.Println("load user succ", "path", uidFile)
		copyUser.RangeUser(logger, RunCopyTask)
		_, _ = writer.Write([]byte("copy end"))
		fkfmt.Println("copy end", "path", uidFile)
	})
}
