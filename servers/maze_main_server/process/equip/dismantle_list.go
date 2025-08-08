package equip

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/config/GDollEquipInfoV8Cfg"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"sort"
)

var MaxEquipNum int32 // 分解列表可显示最大装备数

func init() {
	param.Int32P(&MaxEquipNum, "list:max:equip:num", 1000, "分解列表最大装备数")
}

func (e *Equip) OnMazeEquipDismantleListRQ_10616_10617(s *session.Session, req *MazeGameEquip.MazeEquipDismantleRQ) (err error) {
	defer fkprometheus.DebugPMT("OnMazeEquipDismantleListRQ")()
	res := &MazeGameEquip.MazeEquipDismantleRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	userCtx := log.Clone("Equip", uint64(s.UID()), 0)
	shardingID := uint64(s.UID())

	defer func() {
		userCtx.InfoWF("OnMazeEquipDismantleListRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnMazeEquipDismantleListRQ with", zap.Any("req", req))

	all := GMazeEquipPosRankV8Cfg.GetAll()
	if len(all) == 0 {
		userCtx.ErrorWF("OnMazeEquipDismantleListRQ GDollEquipPosRankV8Cfg empty")
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}

	posRankMap := make(map[int32]int32)
	for _, v := range all {
		posRankMap[v.Pos_id] = v.Rank
	}

	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnMazeEquipDismantleListRQ GetAllEquipInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 获取身上的装备信息
	assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnDollEquipSaleSelectRQ GetAllDollAssembleSuit fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	// 身上的装备
	assembleGuid := make(map[int64]struct{})
	if len(assembleInfoMap) > 0 {
		for _, equipList := range assembleInfoMap {
			if len(equipList) > 0 {
				for _, v := range equipList {
					assembleGuid[v.GetEquipGuid()] = struct{}{}
				}
			}
		}
	}

	// getForceEquip := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	//equipMain := make(map[int64]int64)
	bagEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equip := range equipInfoMap {
		if equip.GetLock() > 0 { // 上锁的装备不能分解
			continue
		}

		if _, ok := assembleGuid[equip.GetEquipGuid()]; ok { // 身上的装备不能分解
			continue
		}

		// 临时装备过滤掉
		// if equip.GetEnterTime() > 0 {
		// 	continue
		// }

		equipCfg := GDollEquipInfoV8Cfg.Get(equip.GetEquipId())
		if equipCfg == nil {
			userCtx.ErrorWF("OnDollEquipDismantleListRQ getEquipInfoCfg fail", zap.Any("equipId", equip.GetEquipId()))
			res.ErrInfo = errors.MODULE_ERROR.Wrap("配置数据错误")
			return
		}

		// 未解封的装备 需要转一下
		newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(userCtx, equip)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			ctx.ErrorWF("OnDollEquipDismantleListRQ ConvertIdentifyEquipDb fail", zap.Error(err), zap.Any("equip", equip))
			return err
		}

		equipInfoPb, err := packtopb.EquipSimplifyToCliPB(userCtx, newEquipInfo)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			ctx.ErrorWF("OnDollEquipDismantleListRQ EquipSimplifyToCliPB fail", zap.Error(err))
			return err
		}

		res.EquipList = append(res.EquipList, equipInfoPb)

		bagEquips = append(bagEquips, newEquipInfo)

		if len(res.EquipList) >= int(MaxEquipNum) {
			break
		}
	}

	err = bagequipstrength.BagEquipStrengthCalc(userCtx, shardingID, bagEquips, res.EquipList, false, "doll_equip_dismantle_server")
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnDollEquipDismantleListRQ BagEquipStrengthCalc fail", zap.Error(err))
		return err
	}

	sort.Slice(res.EquipList, func(i, j int) bool {

		if res.EquipList[i].GetEquipLevel() != res.EquipList[j].GetEquipLevel() {
			return res.EquipList[i].GetEquipLevel() < res.EquipList[j].GetEquipLevel()
		}

		if res.EquipList[i].GetEquipQuality() != res.EquipList[j].GetEquipQuality() {
			return res.EquipList[i].GetEquipQuality() < res.EquipList[j].GetEquipQuality()
		}

		if res.EquipList[i].GetPos() != res.EquipList[j].GetPos() {
			return posRankMap[res.EquipList[i].GetPos()] < posRankMap[res.EquipList[j].GetPos()]
		}

		if res.EquipList[i].GetForceValue() != res.EquipList[j].GetForceValue() {
			return res.EquipList[i].GetForceValue() < res.EquipList[j].GetForceValue()
		}

		if res.EquipList[i].GetStrengthScore() != res.EquipList[j].GetStrengthScore() {
			return res.EquipList[i].GetStrengthScore() < res.EquipList[j].GetStrengthScore()
		}

		var baseI, baseJ int64
		// attrI := res.EquipList[i].GetMainAttrs()
		// if attrI.GetAttrType() == int32(DollEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_MAIN) {
		// 	if attrI.GetAttrInfo() != nil {
		// 		baseI = attrI.GetAttrInfo().GetValue()
		// 	}
		// }
		// attrJ := res.EquipList[j].GetMainAttrs()
		// if attrJ.GetAttrType() == int32(DollEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_MAIN) {
		// 	if attrJ.GetAttrInfo() != nil {
		// 		baseJ = attrJ.GetAttrInfo().GetValue()
		// 	}
		// }

		return baseI < baseJ
		// return equipMain[res.EquipList[i].GetEquipGuid()] < equipMain[res.EquipList[j].GetEquipGuid()]
	})

	return
}
