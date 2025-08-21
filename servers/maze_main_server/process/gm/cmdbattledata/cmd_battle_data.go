/*
 * @Author: majian
 * @Date: 2025-04-03 10:48:23
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-03 19:14:19
 */
package cmdbattledata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maze_game_server/services/tempbuffservice"
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/common/function/gm"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/servers/maze_main_server/process/game"
)

func RegBattleDataGm(logger fklog.FKLogI) {
	gm.SafeHttpRegister(logger, "/DumpBattleData", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		barrierId := fkutil.ToInt32(request.Form.Get("barrierId"))
		logger.SetLogId(time.Now().UnixNano())
		logger.SetUid(userId)
		logger.InfoWF("DumpBattleData begin")
		battleData, e := game.GetMazeBattleData(logger, userId, barrierId)
		if e != nil {
			writer.Write([]byte(e.Error()))
			return
		}
		var bs bytes.Buffer
		bs.WriteString("战斗数据:\n")
		as, _ := json.Marshal(battleData.GetAreaInfos())
		bs.WriteString(fmt.Sprintf("区域怪物数据:%s\n", string(as)))
		rc, _ := json.Marshal(battleData.GetRoleConfigInfo())
		bs.WriteString(fmt.Sprintf("角色配置数据:%s\n", string(rc)))
		em, _ := json.Marshal(battleData.GetEliteMonsterInfos())
		bs.WriteString(fmt.Sprintf("精英怪数据:%s\n", string(em)))

		tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userId, barrierId)
		if err != nil {
			logger.ErrorWF("DumpBattleData GetBarrierTempBuff err", zap.Error(err))
		} else {
			tb, _ := json.Marshal(tempBuffInfo.TotalBuff)
			bs.WriteString(fmt.Sprintf("临时buff数据:%s\n", string(tb)))
		}
		attrMap := game.GetAttrType()
		bs.WriteString("属性ID<->枚举映射关系:\n")
		for k, v := range attrMap {
			var attrName string
			attrCfg := GMazeAttributeV8Cfg.Get(k)
			if attrCfg != nil {
				attrName = attrCfg.Name
			}
			bs.WriteString(fmt.Sprintf("%d(%s)->%d\n", k, attrName, v))
		}
		writer.Write(bs.Bytes())
	})
}
