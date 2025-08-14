package dollmappuzzlenewcfgex

import (
	"go.uber.org/zap"
	"maze_game_server/config/GDollMapPuzzleNewV8Cfg"
	"maze_game_server/lib/log"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"
)

type AreaInfo struct {
	StageId   int32
	AreaId    int32
	AreaIndex int32
}

type DollMapPuzzleNewV8ConfigEx struct {
	BarrierConfigMap map[int32][]*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow
	BarrierAreaMap   map[int32][]*AreaInfo
}

func init() {
	GDollMapPuzzleNewV8Cfg.RegisterDollMapPuzzleNewV8InitCallBack("dollmappuzzlenewcfgex", loadConfigByBarrier)
}

var gConfigDataEx *DollMapPuzzleNewV8ConfigEx

// 代表刷怪区域类型
const (
	ShowType = 2
	SubType  = 1
)

func loadConfigByBarrier(f *GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8Config) {
	logger := log.Clone("dollmappuzzlenewcfgex", 0, 0)
	temp := &DollMapPuzzleNewV8ConfigEx{
		BarrierConfigMap: make(map[int32][]*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow),
		BarrierAreaMap:   make(map[int32][]*AreaInfo),
	}
	for _, i := range f.GetAll() {
		barrierConfigs := temp.BarrierConfigMap[i.Level]
		barrierConfigs = append(barrierConfigs, i)
		temp.BarrierConfigMap[i.Level] = barrierConfigs

		if i.Show_type == ShowType && i.Sub_type == SubType {
			areaList := temp.BarrierAreaMap[i.Level]
			splitArr := strings.Split(i.Monster_area, "_")
			if len(splitArr) != 2 {
				logger.ErrorWF("Monster_area config err.", zap.Any("config", i))
				continue
			}
			areaId, err := strconv.ParseInt(splitArr[0], 10, 32)
			if err != nil {
				logger.ErrorWF("Monster_area config err.", zap.Any("config", i))
				continue
			}
			areaIndex, err := strconv.ParseInt(splitArr[1], 10, 32)
			if err != nil {
				logger.ErrorWF("Monster_area config err.", zap.Any("config", i))
				continue
			}

			areaList = append(areaList, &AreaInfo{StageId: i.Stage, AreaId: int32(areaId), AreaIndex: int32(areaIndex)})
			temp.BarrierAreaMap[i.Level] = areaList
		}
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigDataEx)), unsafe.Pointer(temp))
}

// 根据关卡ID获取所有配置
func GetBarrierConfigs(barrier int32) []*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow {
	return gConfigDataEx.BarrierConfigMap[barrier]
}

// 根据关卡ID获取刷怪区域配置
func GetBarrierAreaInfos(barrier int32) []*AreaInfo {
	return gConfigDataEx.BarrierAreaMap[barrier]
}

// 获取通过的区域
func GetPassAreaInfos(barrierId, stageId int32) []*AreaInfo {
	areaInfos := GetBarrierAreaInfos(barrierId)
	passArea := make([]*AreaInfo, 0)
	for _, i := range areaInfos {
		if i.StageId != 0 || i.StageId > stageId {
			continue
		}
		passArea = append(passArea, i)
	}
	return passArea
}
