package mailservice

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"
	"maze_game_server/pb/common/MazeCommon"
)

// MailAdvancedService 高级邮件功能服务
type MailAdvancedService struct {
	mailService MailService
}

// NewMailAdvancedService 创建高级邮件功能服务
func NewMailAdvancedService() *MailAdvancedService {
	return &MailAdvancedService{
		mailService: GlobalMailService,
	}
}

// VoteMailRequest 投票邮件请求
type VoteMailRequest struct {
	Title       string   `json:"title"`        // 邮件标题
	Content     string   `json:"content"`      // 邮件内容
	VoteId      int32    `json:"vote_id"`      // 投票ID
	VoteTitle   string   `json:"vote_title"`   // 投票标题
	Options     []string `json:"options"`      // 选项列表
	EndTime     int64    `json:"end_time"`     // 结束时间
	MaxChoices  int32    `json:"max_choices"`  // 最大选择数
	ReceiverIds []uint64 `json:"receiver_ids"` // 接收者ID列表
}

// BattleReportMailRequest 战报邮件请求
type BattleReportMailRequest struct {
	Title          string `json:"title"`            // 邮件标题
	Content        string `json:"content"`          // 邮件内容
	BattleRecordId int64  `json:"battle_record_id"` // 战报ID
	BattleData     string `json:"battle_data"`      // 战报数据
	AttackerId     uint64 `json:"attacker_id"`      // 攻击者ID
	DefenderId     uint64 `json:"defender_id"`      // 防御者ID
	BattleResult   int32  `json:"battle_result"`    // 战斗结果
	ReceiverId     uint64 `json:"receiver_id"`      // 接收者ID
}

// GiftBoxMailRequest 礼盒邮件请求
type GiftBoxMailRequest struct {
	Title       string                 `json:"title"`        // 邮件标题
	Content     string                 `json:"content"`      // 邮件内容
	BoxType     int32                  `json:"box_type"`     // 宝箱类型
	BoxId       int32                  `json:"box_id"`       // 宝箱ID
	ItemList    []*MazeCommon.MazeItem `json:"item_list"`    // 物品列表
	ReceiverIds []uint64               `json:"receiver_ids"` // 接收者ID列表
}

// LinkMailRequest 跳转邮件请求
type LinkMailRequest struct {
	Title       string   `json:"title"`        // 邮件标题
	Content     string   `json:"content"`      // 邮件内容
	LinkType    int32    `json:"link_type"`    // 链接类型
	LinkParam   string   `json:"link_param"`   // 跳转参数
	LinkUrl     string   `json:"link_url"`     // 链接URL
	ReceiverIds []uint64 `json:"receiver_ids"` // 接收者ID列表
}

// SendVoteMail 发送投票邮件
func (s *MailAdvancedService) SendVoteMail(ctx context.Context, req *VoteMailRequest) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建投票信息
	voteInfo := &structsdef.VoteInfo{
		VoteId:     req.VoteId,
		Title:      req.VoteTitle,
		Options:    req.Options,
		EndTime:    req.EndTime,
		MaxChoices: req.MaxChoices,
	}

	voteInfoJSON, err := json.Marshal(voteInfo)
	if err != nil {
		logger.CtxError(ctx, "SendVoteMail marshal vote info failed", zap.Error(err))
		return err
	}

	// 构建富文本面板信息
	panelInfo := s.buildVoteMailPanel(req)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 批量发送邮件
	for _, receiverId := range req.ReceiverIds {

		logger := fklog.ContextAppLogger(ctx)
		err := s.mailService.SendMail(logger, req.Title, req.Content, "系统",
			int32(constdef.MailLabelSystem), receiverId, nil, req.EndTime)

		if err != nil {
			logger.CtxError(ctx, "SendVoteMail SendMail failed",
				zap.Uint64("receiverId", receiverId),
				zap.Error(err))
		}
	}

	logger.CtxInfo(ctx, "SendVoteMail completed",
		zap.Int32("voteId", req.VoteId),
		zap.Int("receiverCount", len(req.ReceiverIds)))

	return nil
}

// SendBattleReportMail 发送战报邮件
func (s *MailAdvancedService) SendBattleReportMail(ctx context.Context, req *BattleReportMailRequest) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建富文本面板信息
	panelInfo := s.buildBattleReportMailPanel(req)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 创建邮件信息
	mailInfo := &mailmodel.MailInfo{
		ID:               uint64(time.Now().UnixNano()),
		Title:            req.Title,
		Content:          req.Content,
		Label:            int32(constdef.MailLabelSystem),
		MailType:         int32(constdef.BATTLE),
		MailSubType:      int32(constdef.DOLL_BATTLE_ROB_SUCCESS_MAIL),
		MailContentType:  int32(constdef.CONTENT_DOLL_BATTLE_REPORT),
		Sender:           "系统",
		ReciverID:        req.ReceiverId,
		IsRead:           false,
		IsGetAttach:      false,
		SendTime:         time.Now().Unix(),
		ExpireTime:       time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天后过期
		BattleRecordBody: req.BattleData,
		PanelInfo:        string(panelInfoJSON),
		PagePopUp:        true,
		AllSupportPop:    true,
	}

	// 发送邮件
	err := s.mailService.SendMail(ctx, req.Title, req.Content, "系统",
		int32(constdef.MailLabelSystem), req.ReceiverId, nil, mailInfo.ExpireTime)

	if err != nil {
		logger.CtxError(ctx, "SendBattleReportMail SendMail failed",
			zap.Uint64("receiverId", req.ReceiverId),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "SendBattleReportMail success",
		zap.Int64("battleRecordId", req.BattleRecordId),
		zap.Uint64("receiverId", req.ReceiverId))

	return nil
}

// SendGiftBoxMail 发送礼盒邮件
func (s *MailAdvancedService) SendGiftBoxMail(ctx context.Context, req *GiftBoxMailRequest) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建宝箱信息
	var boxInfo interface{}
	if req.BoxType == 1 {
		// 随机宝箱
		boxInfo = &structsdef.RandomBoxInfo{
			BoxId:    req.BoxId,
			ItemList: req.ItemList,
		}
	} else if req.BoxType == 2 {
		// 贸易宝箱
		boxInfo = &structsdef.TradeBoxInfo{
			BoxId:    req.BoxId,
			ItemList: req.ItemList,
		}
	}

	boxInfoJSON, _ := json.Marshal(boxInfo)

	// 构建富文本面板信息
	panelInfo := s.buildGiftBoxMailPanel(req)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 批量发送邮件
	for _, receiverId := range req.ReceiverIds {
		mailInfo := &mailmodel.MailInfo{
			ID:              uint64(time.Now().UnixNano()),
			Title:           req.Title,
			Content:         req.Content,
			Label:           int32(constdef.MailLabelSystem),
			MailType:        int32(constdef.SYSTEM),
			MailSubType:     int32(constdef.AWARD),
			MailContentType: int32(constdef.CONTENT_GIFT_PACK),
			Sender:          "系统",
			ReciverID:       receiverId,
			IsRead:          false,
			IsGetAttach:     false,
			SendTime:        time.Now().Unix(),
			ExpireTime:      time.Now().Add(30 * 24 * time.Hour).Unix(), // 30天后过期
			PanelInfo:       string(panelInfoJSON),
			PagePopUp:       true,
			AllSupportPop:   true,
		}

		err := s.mailService.SendMail(ctx, req.Title, req.Content, "系统",
			int32(constdef.MailLabelSystem), receiverId, nil, mailInfo.ExpireTime)

		if err != nil {
			logger.CtxError(ctx, "SendGiftBoxMail SendMail failed",
				zap.Uint64("receiverId", receiverId),
				zap.Error(err))
		}
	}

	logger.CtxInfo(ctx, "SendGiftBoxMail completed",
		zap.Int32("boxType", req.BoxType),
		zap.Int32("boxId", req.BoxId),
		zap.Int("receiverCount", len(req.ReceiverIds)))

	return nil
}

// SendLinkMail 发送跳转邮件
func (s *MailAdvancedService) SendLinkMail(ctx context.Context, req *LinkMailRequest) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建跳转参数
	linkParam := &structsdef.LinkParam{
		LinkType: req.LinkType,
		Param:    req.LinkParam,
	}

	linkParamJSON, _ := json.Marshal(linkParam)

	// 构建富文本面板信息
	panelInfo := s.buildLinkMailPanel(req)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 批量发送邮件
	for _, receiverId := range req.ReceiverIds {
		mailInfo := &mailmodel.MailInfo{
			ID:              uint64(time.Now().UnixNano()),
			Title:           req.Title,
			Content:         req.Content,
			Label:           int32(constdef.MailLabelSystem),
			MailType:        int32(constdef.SYSTEM),
			MailSubType:     int32(constdef.AWARD),
			MailContentType: int32(constdef.CONTENT_CLIENT_OPERATA),
			Sender:          "系统",
			ReciverID:       receiverId,
			IsRead:          false,
			IsGetAttach:     false,
			SendTime:        time.Now().Unix(),
			ExpireTime:      time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天后过期
			LinkUrl:         req.LinkUrl,
			LinkParam:       string(linkParamJSON),
			PanelInfo:       string(panelInfoJSON),
			PagePopUp:       true,
			AllSupportPop:   true,
		}

		err := s.mailService.SendMail(ctx, req.Title, req.Content, "系统",
			int32(constdef.MailLabelSystem), receiverId, nil, mailInfo.ExpireTime)

		if err != nil {
			logger.CtxError(ctx, "SendLinkMail SendMail failed",
				zap.Uint64("receiverId", receiverId),
				zap.Error(err))
		}
	}

	logger.CtxInfo(ctx, "SendLinkMail completed",
		zap.Int32("linkType", req.LinkType),
		zap.Int("receiverCount", len(req.ReceiverIds)))

	return nil
}

// buildVoteMailPanel 构建投票邮件面板
func (s *MailAdvancedService) buildVoteMailPanel(req *VoteMailRequest) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: req.Title,
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: req.Content,
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统",
			},
		},
	}
}

// buildBattleReportMailPanel 构建战报邮件面板
func (s *MailAdvancedService) buildBattleReportMailPanel(req *BattleReportMailRequest) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: req.Title,
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: req.Content,
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统",
			},
		},
	}
}

// buildGiftBoxMailPanel 构建礼盒邮件面板
func (s *MailAdvancedService) buildGiftBoxMailPanel(req *GiftBoxMailRequest) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: req.Title,
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: req.Content,
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统",
			},
		},
	}
}

// buildLinkMailPanel 构建跳转邮件面板
func (s *MailAdvancedService) buildLinkMailPanel(req *LinkMailRequest) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: req.Title,
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: req.Content,
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统",
			},
		},
	}
}

// 全局高级邮件功能服务实例
var GlobalMailAdvancedService *MailAdvancedService

func init() {
	GlobalMailAdvancedService = NewMailAdvancedService()
}
