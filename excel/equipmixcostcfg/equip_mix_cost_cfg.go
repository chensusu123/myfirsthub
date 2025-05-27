// @Author: ZhaoXiming 2025/3/24 15:21
// @Desc:

package equipmixcostcfg

import (
	"sort"
	"sync"

	"maze_game_server/config/GMazeEquipMixV8Cfg"
	"maze_game_server/pb/common/MazeCommon"

	"google.golang.org/protobuf/proto"
)

var (
	equipData = make(map[int32]*CfgData)
	lock      sync.RWMutex
)

type CfgData struct {
	Cost []*MazeCommon.MazeItem
	Cfg  *GMazeEquipMixV8Cfg.MazeEquipMixV8ConfigRow
}

func init() {
	GMazeEquipMixV8Cfg.RegisterLoadedCallBack("equip_mix_cfg", func(config *GMazeEquipMixV8Cfg.MazeEquipMixV8Config) {

		m := make(map[int32]*CfgData, len(config.ConfigRows))
		for _, row := range config.ConfigRows {
			cost := make([]*MazeCommon.MazeItem, 0, len(row.Cost))
			for k, v := range row.Cost {
				cost = append(cost, &MazeCommon.MazeItem{
					ItemId: proto.Int32(k),
					Count:  proto.Int64(v),
				})
			}
			if len(cost) > 0 {
				sort.SliceStable(cost, func(i, j int) bool {
					return cost[i].GetItemId() < cost[j].GetItemId()
				})
			}
			m[row.Order] = &CfgData{
				Cost: cost,
				Cfg:  row,
			}
		}

		lock.Lock()
		equipData = m
		lock.Unlock()

	})
}

func GetEquipCost(lv int32) (cfg *CfgData) {
	lock.RLock()
	defer lock.RUnlock()
	cfg = equipData[lv]
	return
}
