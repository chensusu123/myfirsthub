package calequipsequence

//
//import (
//	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
//	"go.uber.org/zap"
//	"maze_game_server/common/errors"
//	"maze_game_server/common/function/randfuncs"
//	"maze_game_server/config/GMazeBarriesV8Cfg"
//	"maze_game_server/config/GMazeShopEquipListV8Cfg"
//	"maze_game_server/config/GMazeShopV8Cfg"
//	"maze_game_server/io/redis/mazeshopseqredis"
//	"maze_game_server/module/mazeuserinfo"
//	"maze_game_server/services/equipdropservice"
//)
//
//func GetNewEquip(ctx context.Context, userId uint64, barrierId int32, level int32, equipNum int32) (newEquip map[int32]int32, err error) {
//
//	if level <= 0 {
//		userInfo, err2 := mazeuserinfo.GetUserInfoV2(logger, userId)
//		if err2 != nil {
//			logger.CtxError(ctx,"GetNewEquip GetNewEquip fail", zap.Error(err2))
//			err = err2
//			return
//		}
//		level = int32(userInfo.Level)
//	}
//
//	newLevel := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(level, barrierId)
//
//	shopInfo, err := GetMazeShopInfo(logger, userId, int32(newLevel), barrierId)
//	if err != nil {
//		logger.CtxError(ctx,"GetNewEquip GetMazeShopInfo fail", zap.Error(err), zap.Any("level", newLevel))
//		return
//	}
//
//	newEquip, err = GetMazeSeqEquipId(logger, userId, int32(newLevel), barrierId, shopInfo, equipNum, 0)
//	if err != nil {
//		logger.CtxError(ctx,"GetNewEquip GetMazeSeqEquipId fail", zap.Error(err),
//			zap.Any("shopInfo", shopInfo))
//		return
//	}
//
//	err = mazeshopseqredis.SetMazeShopInfo(logger, userId, int32(newLevel), shopInfo)
//	if err != nil {
//		logger.CtxError(ctx,"GetNewEquip SetMazeShopInfo fail", zap.Error(err), zap.Any("level", newLevel), zap.Any("shopInfo", shopInfo))
//		return
//	}
//
//	return
//}
//
//func GetMazeShopInfo(ctx context.Context, userId uint64, mazeLevel int32, barrier int32) (mazeShopInfo *mazeshopseqredis.MazeShopInfo, err error) {
//	newLevel := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(mazeLevel, barrier)
//	mazeShopInfo, err = mazeshopseqredis.GetMazeShopInfo(logger, userId, newLevel)
//	if err != nil {
//		logger.CtxError(ctx,"GetMazeShopInfo fail", zap.Error(err), zap.Any("newLevel", newLevel))
//		return
//	}
//	cfg := GMazeShopV8Cfg.GetWithCtx(ctx,newLevel)
//	if cfg == nil {
//		logger.CtxError(ctx,"GetMazeShopInfo GMazeShopV8Cfg fail",
//			zap.Int32("newLevel", newLevel))
//		return mazeShopInfo, errors.New("获取配置失败")
//	}
//	if mazeShopInfo == nil {
//		mazeShopInfo, err = HandleMazeShopSeqInit(logger, cfg)
//		if err != nil || mazeShopInfo == nil {
//			logger.CtxError(ctx,"GetMazeShopInfo HandleMazeShopSeqInit fail",
//				zap.Int32("newLevel", newLevel),
//				zap.Error(err))
//			return mazeShopInfo, errors.New("获取配置失败")
//		}
//	}
//	return mazeShopInfo, nil
//}
//
//func HandleMazeShopSeqInit(ctx context.Context, cfg *GMazeShopV8Cfg.MazeShopV8ConfigRow) (mazeShopInfo *mazeshopseqredis.MazeShopInfo, err error) {
//	mazeShopInfo = &mazeshopseqredis.MazeShopInfo{}
//	if cfg == nil {
//		logger.CtxError(ctx,"HandleMazeShopSeqInit GMazeShopV8Cfg err")
//		return nil, errors.New("配置不存在")
//	}
//	seqWeight := make(map[int32]int32, 0)
//	for _, seqId := range cfg.Buy_list {
//		seqWeight[seqId] = 1000
//	}
//	seqId := randfuncs.RandByWeightV2(logger, seqWeight, false)
//	if seqId <= 0 {
//		logger.CtxError(ctx,"HandleMazeShopSeqInit GMazeShopV8Cfg err",
//			zap.Any("seqWeight", seqWeight))
//		return nil, errors.New("配置不存在")
//	}
//	mazeShopInfo.SeqId = seqId
//	mazeShopInfo.BackId = cfg.Buy_list_base2
//	mazeShopInfo.ShopSlotNum = make(map[int32]int32)
//	return mazeShopInfo, nil
//}
//
//func GetMazeSeqEquipId(ctx context.Context, userId uint64, mazeLevel int32, barrier int32, mazeShopInfo *mazeshopseqredis.MazeShopInfo, addCount, assignPos int32) (equipMap map[int32]int32, err error) {
//	newLevel := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(mazeLevel, barrier)
//	equipMap = make(map[int32]int32)
//	shopCfg := GMazeShopV8Cfg.GetWithCtx(ctx,newLevel)
//	if shopCfg == nil {
//		logger.CtxError(ctx,"GetMazeSeqEquipId GMazeShopV8Cfg fail",
//			zap.Int32("newLevel", newLevel))
//		return equipMap, errors.New("等级配置不存在")
//	}
//
//	var seqCfg *GMazeShopEquipListV8Cfg.MazeShopEquipListV8ConfigRow
//	var backCfg *GMazeShopEquipListV8Cfg.MazeShopEquipListV8ConfigRow
//	for _, cfg := range GMazeShopEquipListV8Cfg.GetAll() {
//		if (cfg.List_id == mazeShopInfo.SeqId) && (cfg.Pos_id == 0) {
//			seqCfg = cfg
//		}
//		if (cfg.List_id == mazeShopInfo.BackId) && (cfg.Pos_id == 0) {
//			backCfg = cfg
//		}
//	}
//	if (seqCfg == nil) || (backCfg == nil) {
//		logger.CtxError(ctx,"GetMazeShopSeq GMazeShopEquipListV8Cfg err",
//			zap.Any("mazeShopInfo", mazeShopInfo))
//		return equipMap, errors.New("配置不存在")
//	}
//	for i := int32(0); i < addCount; i++ {
//		var equipId int32
//		if mazeShopInfo.BackIndex > 0 || (int(mazeShopInfo.CurSeqIndex) >= len(seqCfg.Equip_id)) {
//			if int(mazeShopInfo.BackIndex) >= len(backCfg.Equip_id) {
//				mazeShopInfo.BackIndex = 0
//			}
//			mazeShopInfo.BackIndex++
//			mazeShopInfo.TotalCount++
//			equipId = backCfg.Equip_id[mazeShopInfo.BackIndex-1]
//		} else {
//			mazeShopInfo.CurSeqIndex++
//			mazeShopInfo.TotalCount++
//			equipId = seqCfg.Equip_id[mazeShopInfo.CurSeqIndex-1]
//		}
//		convEquipId, err := ConvEquipId(logger, equipId, assignPos)
//		if equipId != convEquipId {
//			logger.CtxInfo(ctx,"GetMazeSeqEquipId conv equipId",
//				zap.Int32("equipId", equipId),
//				zap.Int32("convEquipId", convEquipId),
//				zap.Any("assignPos", assignPos))
//		}
//		if err != nil || convEquipId == 0 {
//			logger.CtxError(ctx,"GetMazeShopSeq ConvEquipId err",
//				zap.Error(err),
//				zap.Any("equipId", equipId))
//			return equipMap, errors.New("配置不存在")
//		}
//		equipMap[convEquipId] += 1
//	}
//	logger.CtxInfo(ctx,"GetMazeSeqEquipId end",
//		zap.Int32("newLevel", newLevel),
//		zap.Int32("mazeLevel", mazeLevel),
//		zap.Int32("barrier", barrier),
//		zap.Any("mazeShopInfo", mazeShopInfo),
//		zap.Any("addCount", addCount),
//		zap.Any("assignPos", assignPos),
//		zap.Any("equipMap", equipMap))
//	return equipMap, nil
//}
//
//func ConvEquipId(ctx context.Context, equipId int32, assignPos int32) (int32, error) {
//	convEquipId := equipId
//	if assignPos > 0 {
//		if assignPos == 1 {
//			convEquipId = 427*100000 + (equipId/10000)%10*10000 + assignPos*1000 + equipId%1000
//		} else {
//			convEquipId = 426*100000 + (equipId/10000)%10*10000 + assignPos*1000 + equipId%1000
//		}
//	}
//
//	return convEquipId, nil
//}
//
//func GetMazeBarrierLv(level int32, barrier int32) int32 {
//	cfg := GMazeBarriesV8Cfg.GetWithCtx(ctx,barrier)
//	if cfg == nil {
//		return level
//	}
//	if level < cfg.Drop_equip_lv_min {
//		return cfg.Drop_equip_lv_min
//	}
//	if level > cfg.Drop_equip_lv_max {
//		return cfg.Drop_equip_lv_max
//	}
//	return level
//}
