package mail

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/mailmodel"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeMail"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/mailservice"
)

// 邮件附件领取
func (g *Mail) OnMazeGetMailAttachmentsRQ_10630_10631(s *session.Session, req *MazeMail.MazeGetMailAttachmentsRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeGetMailAttachmentsRQ")()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeMail.MazeGetMailAttachmentsRS{}

	logger.InfoWF("OnMazeGetMailAttachmentsRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeGetMailAttachmentsRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeGetMailAttachmentsRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	mails := make([]*mailmodel.MailInfo, 0)
	equipMap := make(map[int32]int32) //装备列表
	itemMap := make(map[int32]int64)  //道具列表

	if req.GetIsAll() {
		mailList, attachements, err := mailservice.GlobalMailService.GetAllMailAttachment(logger, userId, req.GetLabel())
		if err != nil {
			logger.ErrorWF("OnMazeGetMailAttachmentsRQ GetAllMailAttachment fail", zap.Error(err), zap.Int32("label", req.GetLabel()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return err
		}

		mails = append(mails, mailList...)

		for _, attach := range attachements {
			if attach.Extra == "equip" {
				equipMap[attach.ItemID] += int32(attach.Count)
			} else {
				itemMap[attach.ItemID] += attach.Count
			}
		}

	} else {
		mailInfo, attachements, err := mailservice.GlobalMailService.GetMailAttachment(logger, userId, req.GetMailId(), req.GetLabel())
		if err != nil {
			logger.ErrorWF("OnMazeGetMailAttachmentsRQ GetMailAttachment fail", zap.Error(err), zap.Uint64("mId", req.GetMailId()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return err
		}
		mails = append(mails, mailInfo)

		for _, attach := range attachements {
			if attach.Extra == "equip" {
				equipMap[attach.ItemID] += int32(attach.Count)
			} else {
				itemMap[attach.ItemID] += attach.Count
			}
		}
	}

	tradeNo := gentradeno.GetTradeNum()
	if len(equipMap) > 0 {
		//TODO 差一个邮件领取枚举，看业务是否需要邮件支持发装备附件
		_, err := addequip.AddEquipToBag(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SWEEP_AWARD), tradeNo, equipMap)
		if err != nil {
			logger.ErrorWF("OnMazeGetMailAttachmentsRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equipMap))
		}
	}

	// 处理需要加入背包的道具
	if len(itemMap) > 0 {
		realAddItemList := make([]*MazeCommon.MazeItem, 0)
		for k, v := range itemMap {
			realAddItemList = append(realAddItemList, &MazeCommon.MazeItem{
				ItemId: proto.Int32(k),
				Count:  proto.Int64(v),
			})
		}

		errInfo := gentradeno.AddItemEx(logger, userId, 695, tradeNo, req.GetHeader(), realAddItemList...)
		if errInfo != nil {
			logger.ErrorWF("OnMazeGetMailAttachmentsRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", realAddItemList))
		}
	}

	err = mailservice.GlobalMailService.GetMailAttachmentAfter(logger, userId, mails)
	if err != nil {
		logger.ErrorWF("OnMazeGetMailAttachmentsRQ GetMailAttachmentAfter fail", zap.Error(err))
		return err
	}

	list, err := mailservice.GlobalMailService.GetMailListByLabel(logger, userId, req.GetLabel(), 0, 30)
	if err != nil {
		logger.ErrorWF("OnMazeGetMailAttachmentsRQ GetMailList fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	for _, info := range list {
		res.MailList = append(res.MailList, PbMailData(info))
	}
	return nil
}
