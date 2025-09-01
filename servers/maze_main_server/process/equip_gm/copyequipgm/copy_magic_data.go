package copyequipgm

// func CopyMagicData(ctx context.Context, srcUserId uint64, dstUsers []uint64, param copyinterface.CopyParam) error {

// 	var err error
// 	var magicInfo *DollEquipMagicCache.EquipMagicDb
// 	var curUser uint64
// 	defer func() {
// 		if err != nil {
// 			logger.CtxError(ctx,"CopyMagicData fail", zap.Error(err),
// 				zap.Uint64("src", srcUserId),
// 				zap.Uint64("dst", curUser))
// 		} else {
// 			logger.CtxInfo(ctx,"CopyMagicData succ",
// 				zap.Uint64("src", srcUserId),
// 				zap.Int("dstLen", len(dstUsers)))
// 		}
// 	}()

// 	magicInfo, err = dollassembleredis.GetEquipMagicDb(logger, srcUserId)
// 	if err != nil {
// 		return err
// 	}

// 	for _, dstId := range dstUsers {
// 		curUser = dstId

// 		// 获取数据
// 		nextRow := GDollEquipMagicStrengthV8Cfg.GetDollEquipMagicStrengthV8Config(excelutil.GetEquipMagicStrengthenKey(magicInfo.GetLevel()))
// 		if nextRow == nil {
// 			return errors.New("未找到法宝强化配置")
// 		}

// 		// 拆分属性
// 		forceAttrs, others, err := attributeexcel.SplitForceAndOther(nextRow.Add_attr)
// 		if err != nil {
// 			logger.CtxError(ctx,"CopyMagicData SplitForceAndOther fail",
// 				zap.Int32("level", magicInfo.GetLevel()),
// 				zap.Int32("order", nextRow.Order),
// 				zap.Any("attrs", nextRow))
// 			return err
// 		}

// 		// 写法宝数据
// 		err = dollassembleredis.SetEquipMagicDb(logger, dstId, magicInfo)
// 		if err != nil {
// 			return err
// 		}

// 		// 计算非武力值buff加成  后续有穿戴数据拷贝通知计算 这里不通知
// 		attrDb := calcassembleattr.PackMapAttrAll(others, nextRow.Show_attr)
// 		e := dollattrredis.SaveDollAttr(logger, uint64(dstId), constdef.DollAttrSrcEquipMagicStrengthen, attrDb)
// 		if e != nil {
// 			logger.CtxError(ctx,"CopyMagicData SaveDollAttr err", zap.Error(e), zap.Any("attrDb", attrDb))
// 		}

// 		// 计算武力值   后续有穿戴数据拷贝通知计算 这里不通知
// 		err = dollassembleattrredis.SetDollAssembleAttr(logger, uint64(dstId), constdef.EquipMagicStrengthenForce, forceAttrs)
// 		if err != nil {
// 			logger.CtxError(ctx,"CopyMagicData SetDollAssembleAttr err", zap.Error(err), zap.Any("forceAttrs", forceAttrs))
// 		}

// 		logger.CtxInfo(ctx,"CopyMagicData user succ",
// 			zap.Uint64("src", srcUserId),
// 			zap.Int("dst", int(dstId)))
// 	}
// 	return nil
// }
