package dollassembleredis

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
)

var (
	gRedis = &fkredis.FkRedis{}
)

func init() {
	// 21637 maze:assemble:info:u:%llu 迷宫游戏装配数据存储
	fkconfig.RegisterNameNode("dollassembleredis", 21637, gRedis)
}

// 获取人偶装配信息
// func GetDollAssembleInfo(logger fklog.FKLogI, userId uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
// 	assembleInfo = new(MazeEquipCache.MazeAssembleDb)
// 	key := fmt.Sprintf("maze:assemble:info:u:%d",userId)
// 	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "hgetall", key))
// 	if err == redis.ErrNil {
// 		err = nil
// 		logger.InfoWF("GetDollAssembleInfo hgetall nil", zap.String("key", key))
// 		return
// 	}
// 	if err != nil {
// 		logger.ErrorWF("GetDollAssembleInfo hgetall fail", zap.Error(err), zap.String("key", key))
// 		return
// 	}

// 	sLen := len(res)
// 	for i := 0; i < sLen; i += 2 {
// 		ks := fkutil.ByteSliceToString(res[i])

// 		// 解析神器
// 		if ks == constdef.AssemblePrefixMagicWeapon {
// 			magicWeaponPb := &DollEquipCache.MagicWeaponDB{}
// 			e := proto.Unmarshal(res[i+1], magicWeaponPb)
// 			if e != nil {
// 				logger.ErrorWF("GetDollAssembleInfo Unmarshal MagicWeaponDB fail", zap.Error(e),
// 					zap.String("key", key))
// 				return nil, e
// 			}
// 			assembleInfo.MagicWeapon = magicWeaponPb
// 			continue
// 		}

// 		if ks == constdef.AssemblePrefixCurAssembleSuitIndex {
// 			v, e := fkutil.Bytes2Int64(res[i+1])
// 			if e != nil {
// 				err = e
// 				logger.ErrorWF("GetDollAssembleInfo Unmarshal suit index fail", zap.Error(err),
// 					zap.String("key", key))
// 				return nil, err
// 			}
// 			v32 := int32(v)
// 			assembleInfo.CurSuitIndex = proto.Int32(v32)
// 			continue
// 		}
// 		if ks == constdef.AssemblePrefixSwitchSuitTime {
// 			v, e := fkutil.Bytes2Int64(res[i+1])
// 			if e != nil {
// 				err = e
// 				logger.ErrorWF("GetDollAssembleInfo Unmarshal switch time fail", zap.Error(err),
// 					zap.String("key", key))
// 				return nil, err
// 			}
// 			assembleInfo.SwitchSuitTime = proto.Int64(v)
// 			continue
// 		}
// 		if ks == constdef.AssemblePrefixMount {
// 			v, e := fkutil.Bytes2Int64(res[i+1])
// 			if e != nil {
// 				err = e
// 				logger.ErrorWF("GetDollAssembleInfo Unmarshal mountId fail", zap.Error(err),
// 					zap.String("key", key))
// 				return nil, err
// 			}
// 			assembleInfo.MountId = proto.Int32(int32(v))
// 			continue
// 		}
// 		if pos := assemble.DecodeAssemblePosField(ks); pos > 0 {
// 			posInfo := &MazeEquipCache.MazeEquipSlotDb{}
// 			e := proto.Unmarshal(res[i+1], posInfo)
// 			if e != nil {
// 				err = e
// 				logger.ErrorWF("GetDollAssembleInfo Unmarshal equip pos fail", zap.Error(err),
// 					zap.String("key", key), zap.Int32("pos", pos))
// 				return nil, err
// 			}
// 			posTotalInfo := &MazeEquipCache.MazeEquipPosInfo{}
// 			posTotalInfo.EquipPos = posInfo
// 			assembleInfo.MazeEquips = append(assembleInfo.MazeEquips, posTotalInfo)
// 			continue
// 		}
// 	}
// 	logger.InfoWF("GetDollAssembleInfo hgetall succ", zap.Any("res", assembleInfo), zap.String("key", key))
// 	return assembleInfo, err
// }
