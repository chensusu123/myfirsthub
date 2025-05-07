/*
 * @Author: majian
 * @Date: 2024-04-30 13:44:58
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-25 22:38:22
 */
package equipaassemblegm

import (
	"fmt"

	"gitlab.ifreetalk.com/maze/maze_equip_server/common/function/packtopb"
	"gitlab.ifreetalk.com/maze/maze_equip_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_equip_server/io/rpc/dollequipbagrpc"
	"gitlab.ifreetalk.com/maze/maze_equip_server/module/effectequip"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosRankV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/plate/protodef/MazeGameEquip"
)

type EquipParam struct {
	DressLv    int32 // 穿戴等级
	DisDressLv int32 // 允许距离穿戴等级差值
	Quality    int32 // 品质
	FiveElme   int32 // 属性元素
	Pos        int32 // 指定部位
	SuitId     int32 // 指定套装
	SubType    int32 // 装备子类型
}

type EquipResult struct {
	Equip *MazeGameEquip.MazeEquipInfo
	Pass  bool   // 是否满足条件
	Desc  string // 结果描述
}

func CheckEquipParam(p *EquipParam) error {
	if p.DressLv <= 0 || p.DressLv > 0 && p.DressLv > 200 {
		return fmt.Errorf("穿戴等级无效(0~200)")
	}
	if p.DisDressLv > p.DressLv {
		return fmt.Errorf("穿戴等级下限差值无效")
	}
	if p.Quality < 0 || p.Quality > 0 && p.Quality > 6 {
		return fmt.Errorf("装备品质无效(1~6)")
	}
	if p.FiveElme < 0 || p.FiveElme > 0 && p.FiveElme > 5 {
		return fmt.Errorf("五行相性无效(1~5)")
	}
	if p.Pos < 0 || p.Pos > 0 && p.Pos > 6 {
		return fmt.Errorf("装备位无效(1~6)")
	}
	if p.SuitId < 0 || p.SuitId > 0 && p.SuitId > 5 {
		return fmt.Errorf("套装Id无效(1~5)")
	}
	return nil
}

func AddEquipByCond(logger fklog.FKLogI, userId uint64, cond EquipParam) (result []*EquipResult, err error) {
	equipIds, findAll := FindEquipIdsByCond(logger, cond)
	if !findAll {
		err = fmt.Errorf("按条件未找到装备配置，请检查参数 已找到%d条配置", len(equipIds))
		return
	}
	equipMap, err := AddCondEquipToBag(logger, userId, equipIds, cond.SuitId, cond.SubType)
	if err != nil {
		return
	}
	result = CheckAddResult(cond, equipMap)
	return
}

// 根据指定条件查询装备ID
func FindEquipIdsByCond(logger fklog.FKLogI, cond EquipParam) (equipIds map[int32]int32, findAll bool) {
	equipIds = make(map[int32]int32)
	allRow := GMazeEquipInfoV8Cfg.GetAllMazeEquipInfoV8Config()
	posCnt := len(GMazeEquipPosRankV8Cfg.GetAll())
	for i := 0; i < len(allRow); i++ {
		if cond.DressLv > 0 {
			if allRow[i].Level <= cond.DressLv && allRow[i].Level >= cond.DressLv-cond.DisDressLv {
			} else {
				continue
			}
		}

		if cond.Quality > 0 && allRow[i].Quality != cond.Quality {
			continue
		}

		// 指定部位不满足
		if cond.Pos > 0 && allRow[i].Pos != cond.Pos {
			continue
		}
		if cond.SuitId > 0 {
			if allRow[i].Quality == 8 && allRow[i].Suite_id[cond.SuitId] > 0 {
			} else {
				continue
			}
			// TODO 武器类型不匹配,后面改成随机的话，得背包实例化支持下
			if allRow[i].Pos == 1 && cond.SubType > 0 {
				score, ok := allRow[i].Sub_type_random[cond.SubType]
				if !ok || score <= 0 {
					continue
				}
			}
		}
		if _, ok := equipIds[allRow[i].Pos]; !ok {
			equipIds[allRow[i].Pos] = allRow[i].Equipment_id
		}
		if cond.SuitId > 0 {
			// 套装是8个
			if len(equipIds) >= 8 {
				findAll = true
				break
			}
		}
		if cond.Pos > 0 && len(equipIds) >= 1 {
			findAll = true
			break
		}
		if cond.Pos <= 0 && len(equipIds) == posCnt {
			findAll = true
			break
		}
	}
	return equipIds, findAll
}

func AddCondEquipToBag(logger fklog.FKLogI, userId uint64, equipIds map[int32]int32, suitId, subType int32) (equipInfos map[int64]*MazeGameEquip.MazeEquipInfo, err error) {
	tradeNo := uniqueid.GenUniqueIdUInt64()
	rqAdd := &MazeEquipSvr.SvrAddMazeEquipRQ{
		UserId:      proto.Uint64(userId),
		OpType:      proto.Int32(int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_INIT_EQUIP)),
		TradeNumber: proto.Uint64(tradeNo),
	}

	for pos, v := range equipIds {
		equipCond := &MazeEquipSvr.SvrEquipInfo{
			EquipId: proto.Int32(v)}
		if suitId > 0 {
			cond := &MazeEquipSvr.ConditionInfo{}
			cond.Id = proto.Int32(1)
			cond.Value = proto.Int64(1)
			cond1 := &MazeEquipSvr.ConditionInfo{}
			cond1.Id = proto.Int32(2)
			cond1.Value = proto.Int64(int64(suitId))
			equipCond.Conditions = append(equipCond.Conditions, cond, cond1)
		}
		if pos == 1 && subType > 0 {
			cond := &MazeEquipSvr.ConditionInfo{
				Id:    proto.Int32(4),
				Value: proto.Int64(int64(subType)),
			}
			equipCond.Conditions = append(equipCond.Conditions, cond)
		}
		rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
	}

	rsAdd := &MazeEquipSvr.SvrAddMazeEquipRS{}
	err = dollequipbagrpc.MazeBagAddRQ(logger, rqAdd, rsAdd)
	if err != nil {
		return
	}
	equipInfos = make(map[int64]*MazeGameEquip.MazeEquipInfo)
	if rsAdd.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		equips := rsAdd.GetEquipList()
		var guidList []int64

		for _, equip := range equips {
			guidList = append(guidList, equip.GetEquipGuid())
		}
		if len(guidList) > 0 {
			equipDatails, e := effectequip.BatchGetEffectEquipInfo(logger, userId, guidList...)
			if e == nil {
				for _, equipDetail := range equipDatails {
					equipCli, e1 := packtopb.EquipInfoToCliPB(logger, equipDetail)
					if e1 == nil {
						equipInfos[equipDetail.GetEquipGuid()] = equipCli
					}
				}
			}
		}
		return
	}
	err = errors.New("装备加背包失败")
	return
}

func CheckAddResult(cond EquipParam, equipInfos map[int64]*MazeGameEquip.MazeEquipInfo) []*EquipResult {
	var rs []*EquipResult
	for _, equip := range equipInfos {
		result := &EquipResult{}
		result.Equip = equip
		result.Pass = true
		result.Desc = "ok"
		rs = append(rs, result)

		if equip.GetEquipLevel() > cond.DressLv || equip.GetEquipLevel() < cond.DressLv-cond.DisDressLv {
			result.Pass = false
			result.Desc = packNoMatchDesc("穿戴等级不匹配", cond.DressLv, equip.GetEquipLevel())
			continue
		}
		if cond.Pos > 0 && equip.GetPos() != cond.Pos {
			result.Pass = false
			result.Desc = packNoMatchDesc("装备位不匹配", cond.Pos, equip.GetPos())
			continue
		}
		if cond.Quality > 0 && equip.GetEquipQuality() != cond.Quality {
			result.Pass = false
			result.Desc = packNoMatchDesc("装备品质不匹配", cond.Quality, equip.GetEquipQuality())
			continue
		}

	}
	return rs
}

func packNoMatchDesc(title string, need int32, has int32) string {
	return fmt.Sprintf("%s need %d has %d", title, need, has)
}
