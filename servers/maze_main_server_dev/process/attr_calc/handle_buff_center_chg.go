/*
 * @Author: majian
 * @Date: 2024-07-10 11:36:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-25 22:08:03
 */
package attr_calc

// func MazeBuffKafkaConsumer(c context.Context, logger fklog.FKLogI, index int, key, data []byte) error {
// 	defer fkprometheus.DebugPMT("MazeBuffKafkaConsumer")()
// 	msg := &BuffManager.BuffChangeNotify{}
// 	err := proto.Unmarshal(data, msg)
// 	if err != nil {
// 		logger.ErrorWF("MazeBuffKafkaConsumer Unmarshal fail", zap.Error(err), zap.String("msg", string(data)))
// 		return err
// 	}
// 	userId := msg.GetUserId()

// 	nLogger := logger.Clone("MazeBuffKafkaConsumer")
// 	nLogger.SetUid(userId)
// 	nLogger.SetLogId(time.Now().UnixNano())

// 	r2, bcSrc2 := HandleBCMazeAttrCalc(nLogger, userId, msg)
// 	if r2 {
// 		doMazeAttrCalcRetry(logger, BCtoAttrChgMsg(logger, userId, bcSrc2, int64(time.Now().UnixNano()/1000000)))
// 	}

// 	return err
// }

// func BCtoAttrChgMsg(logger fklog.FKLogI, userId uint64, subType int32, stamp int64) *structsdef.MazeCalcAttrNotifyMsg {
// 	bcMsg := &structsdef.MazeCalcAttrNotifyMsg{}
// 	bcMsg.FromServer = "buff中心服务"
// 	bcMsg.ChgType = constdef.MazeBuffCenter
// 	bcMsg.ChgDesc = fmt.Sprintf("oldBc_%d", subType)
// 	bcMsg.Stamp = stamp
// 	bcMsg.BuffSrc = constdef.MazeBuffSrcOldBC
// 	bcMsg.UserId = userId
// 	return bcMsg
// }

// func HandleBCMazeAttrCalc(logger fklog.FKLogI, userId uint64, msg *BuffManager.BuffChangeNotify) (retry bool, src int32) {
// 	var needCare bool
// 	var bcSrc int32

// 	for _, buff := range msg.GetChanges() {
// 		if commonlogic.IsMazeCalcAttr(buff.GetBuffId()) {
// 			needCare = true
// 			bcSrc = buff.GetSrcType()
// 			break
// 		}
// 	}

// 	logger.WarnWF("HandleBCMazeAttrCalc pop maze buff change", zap.Any("msg", msg), zap.Bool("care", needCare))

// 	if !needCare {
// 		return
// 	}

// 	needRetry, e := mazeattrlogic.RunDacFromBC(logger, userId, bcSrc)
// 	if e != nil && needRetry {
// 		return true, bcSrc
// 	}
// 	return
// }
