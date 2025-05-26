/*
 * @Author: majian
 * @Date: 2024-08-22 17:35:59
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 15:42:41
 * @Desc 装备套装信息打包
 */
package equipsuittopb

import (
	"errors"
	"fmt"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/config/GMazeEquipSuiteAttrV8Cfg"
	"maze_game_server/config/GMazeEquipSuiteInfoV8Cfg"
	"maze_game_server/config/GMazeEquipSuiteNameV8Cfg"
	"maze_game_server/excel/mazeequipinfocfgex"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
)

func PackEquipSuitCliPb(logger fklog.FKLogI, pos int32, equips []*MazeEquipCache.MazeEquipPosInfo, suitId int32, dollLv int32) (cliSuitInfo *MazeGameEquip.EquipSuitInfo, err error) {
	suitPosMap := make(map[int32]*MazeEquipCache.MazeEquipPosInfo) // 跟选定装备成套的装备 <pos,guid>
	if suitId <= 0 {
		return nil, nil
	}

	var suitNum int32 // 激活套装的数量
	for _, equip := range equips {
		if equip.GetEquipInfo().GetSuitId() == suitId && suitId > 0 {
			suitPosMap[equip.GetEquipPos().GetPos()] = equip
			suitNum++
		}
	}

	cliSuitInfo = &MazeGameEquip.EquipSuitInfo{}

	cliSuitInfo.SuitId = proto.Int32(suitId)
	desV8Row := GMazeEquipSuiteInfoV8Cfg.GetMazeEquipSuiteInfoV8Config(suitId)
	if desV8Row == nil {
		logger.ErrorWF("PackEquipSuitCliPb cannot find cfg",
			zap.Int32("suitId", suitId),
			zap.String("sheet", GMazeEquipSuiteInfoV8Cfg.GetConfigDesc()))
		return nil, errors.New("cannot find suit cfg")
	}
	cliSuitInfo.SuitName = proto.String(desV8Row.Suite_name)

	// 套装位置打包
	for _, suitPos := range desV8Row.Pos_list {
		suitPosPb := MazeGameEquip.SuitPos{}
		suitPosPb.PosId = proto.Int32(suitPos)
		equipInfo := suitPosMap[suitPos]
		var equipName string

		// rankRow := GDollEquipPosRankV8Cfg.GetDollEquipPosRankV8Config(suitPos)
		// if rankRow != nil {
		// 	posName = rankRow.Name
		// }
		if equipInfo != nil {
			suitPosPb.Guid = proto.Int64(equipInfo.GetEquipLoadInfo().GetEquipGuid())
			var subType int32
			equipCfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(equipInfo.GetEquipLoadInfo().GetEquipId())
			if equipCfg != nil {
				subType = equipCfg.Pos_sub_type
			}
			_, equipName = pbutil.GetDollEquipName(equipInfo.GetEquipInfo(), subType)

		} else {
			equipName = FindSuitEquipName(logger, suitPos, dollLv, suitId)
		}
		suitPosPb.SuitEquipName = proto.String(equipName)
		cliSuitInfo.SuitPosList = append(cliSuitInfo.SuitPosList, &suitPosPb)
	}
	sort.Slice(cliSuitInfo.SuitPosList, func(i, j int) bool {
		return cliSuitInfo.SuitPosList[i].GetPosId() <= cliSuitInfo.SuitPosList[j].GetPosId()
	})

	// 套装数量打包
	allRow := GMazeEquipSuiteAttrV8Cfg.GetAllMazeEquipSuiteAttrV8Config()
	for _, row := range allRow {
		if row.Suite_id != suitId {
			continue
		}
		suitNumPb := &MazeGameEquip.SuitNum{}
		suitNumPb.Num = proto.Int32(row.Affix_num)
		suitNumPb.SuitEffect = proto.String(row.Add_attr_desc)
		if suitNum >= row.Affix_num { // 已激活
			suitNumPb.Activate = proto.Int32(1)
		}
		cliSuitInfo.SuitNumList = append(cliSuitInfo.SuitNumList, suitNumPb)
	}
	sort.Slice(cliSuitInfo.SuitNumList, func(i, j int) bool {
		return cliSuitInfo.SuitNumList[i].GetNum() <= cliSuitInfo.SuitNumList[j].GetNum()
	})
	return
}

// 未获得的套装装备获取名字
// 先找不高于当前等级最接近的等级的装备
// 如果没有,再找高于当前等级最接近当前等级的装备
func FindSuitEquipName(logger fklog.FKLogI, pos int32, lv int32, suitId int32) string {
	allRows := mazeequipinfocfgex.GetEquipCfgRows(pos, constdef.EquipQualityOrange)
	var needEquipLessCfg, needEquipMoreCfg []*GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow

	for _, row := range allRows {
		if row.Suite_id[suitId] <= 0 {
			continue
		}
		// 过滤没有配置套装名字的装备
		// suitEquipCfg := GMazeEquipSuiteNameV8Cfg.Get(row.Equipment_id)
		// if suitEquipCfg == nil {
		// 	continue
		// }
		if row.Level <= lv {
			needEquipLessCfg = append(needEquipLessCfg, row)
		} else {
			needEquipMoreCfg = append(needEquipMoreCfg, row)
		}
	}

	var aimEquipId int32
	var aimName string
	defer func() {
		logger.InfoWF("FindSuitEquipName",
			zap.String("aimName", aimName),
			zap.Int32("aimEquipId", aimEquipId),
			zap.Int32("suitId", suitId),
			zap.Int32("lv", lv),
			zap.Int32("pos", pos))
	}()
	if len(needEquipLessCfg) > 0 {
		sort.Slice(needEquipLessCfg, func(i, j int) bool {
			if needEquipLessCfg[i].Level != needEquipLessCfg[j].Level {
				return needEquipLessCfg[i].Level > needEquipLessCfg[j].Level
			} else {
				return needEquipLessCfg[i].Equipment_id <= needEquipLessCfg[j].Equipment_id
			}
		})
		aimEquipId = needEquipLessCfg[0].Equipment_id
	} else {
		if len(needEquipMoreCfg) > 0 {
			sort.Slice(needEquipMoreCfg, func(i, j int) bool {
				if needEquipMoreCfg[i].Level != needEquipMoreCfg[j].Level {
					return needEquipMoreCfg[i].Level < needEquipMoreCfg[j].Level
				} else {
					return needEquipMoreCfg[i].Equipment_id <= needEquipMoreCfg[j].Equipment_id
				}
			})
			aimEquipId = needEquipMoreCfg[0].Equipment_id
		}
	}
	if aimEquipId > 0 {
		suitEquipCfg := GMazeEquipSuiteNameV8Cfg.Get(aimEquipId)
		if suitEquipCfg != nil {
			aimName = suitEquipCfg.Suite_equip_name[suitId]
		} else {
			aimName = fmt.Sprintf("未配置套装名(%d)", aimEquipId)
		}
	}
	return aimName
}
