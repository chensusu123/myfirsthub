package mailservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"

	// "maze_game_server/config/GMailInfoCfg"
	// "maze_game_server/config/GWildcardCfg"
	"maze_game_server/model/mailmodel"
)

// MailConfigService 邮件配置服务
type MailConfigService struct {
	// mailInfoConfig *GMailInfoCfg.MailInfoConfig
	// wildcardConfig *GWildcardCfg.WildcardConfig
}

// NewMailConfigService 创建邮件配置服务
func NewMailConfigService() *MailConfigService {
	return &MailConfigService{
		// mailInfoConfig: GMailInfoCfg.GetMailInfoConfig(),
		// wildcardConfig: GWildcardCfg.GetWildcardConfig(),
	}
}

// BuildMailFromConfig 根据配置ID构建邮件
func (s *MailConfigService) BuildMailFromConfig(ctx context.Context, configId int32, userId uint64, wildcardData map[string]interface{}) (*mailmodel.MailInfo, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取邮件配置
	mailCfg := s.mailInfoConfig.GetMailInfoConfig(configId)
	if mailCfg == nil {
		logger.CtxError(ctx, "BuildMailFromConfig config not found", zap.Int32("configId", configId))
		return nil, fmt.Errorf("mail config not found: %d", configId)
	}

	// 构建富文本面板信息
	panelInfo, err := s.buildPanelInfoFromConfig(mailCfg, wildcardData)
	if err != nil {
		logger.CtxError(ctx, "BuildMailFromConfig buildPanelInfoFromConfig failed", zap.Error(err))
		return nil, err
	}

	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 构建邮件信息
	mailInfo := &mailmodel.MailInfo{
		ID:              uint64(ctx.Value("mail_id").(int64)),
		Title:           mailCfg.Mail_title,
		Content:         mailCfg.Mail_text,
		Label:           int32(constdef.MailLabelSystem),
		MailType:        s.parseMailType(mailCfg.Mail_type),
		MailSubType:     s.parseMailSubType(mailCfg.Type),
		MailContentType: int32(constdef.CONTENT_TEXT_AND_AWARD),
		Sender:          "系统",
		ReciverID:       userId,
		IsRead:          false,
		IsGetAttach:     false,
		SendTime:        ctx.Value("send_time").(int64),
		ExpireTime:      s.calculateExpireTime(mailCfg.Lose_time),
		ExtType:         mailCfg.Service_type,
		PanelInfo:       string(panelInfoJSON),
		PagePopUp:       true,
		AllSupportPop:   true,
	}

	logger.CtxInfo(ctx, "BuildMailFromConfig success",
		zap.Int32("configId", configId),
		zap.Uint64("userId", userId))

	return mailInfo, nil
}

// buildPanelInfoFromConfig 根据配置构建面板信息
func (s *MailConfigService) buildPanelInfoFromConfig(mailCfg interface{}, wildcardData map[string]interface{}) (*structsdef.MailPanelInfo, error) {
	panelInfo := &structsdef.MailPanelInfo{
		Contents: make([]*structsdef.MailPanelContent, 0),
		User:     make([]*structsdef.MailUserInfo, 0),
		Goods:    make([]*structsdef.MailGoodsInfo, 0),
		Value:    make([]*structsdef.MailValueInfo, 0),
	}

	// 构建邮件开头
	if mailCfg.Mail_text_start != "" {
		content := s.replaceWildcards(mailCfg.Mail_text_start, mailCfg.Text_start_wildcard, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
			Content: content,
		})
	}

	// 构建邮件内容
	if mailCfg.Mail_text != "" {
		content := s.replaceWildcards(mailCfg.Mail_text, mailCfg.Text_wildcard, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
			Content: content,
		})
	}

	// 构建邮件结尾
	if mailCfg.Mail_text_end != "" {
		content := s.replaceWildcards(mailCfg.Mail_text_end, mailCfg.Text_end_wildcard, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
			Content: content,
		})
	}

	// 构建用户信息
	for _, wildcardId := range mailCfg.Text_start_wildcard {
		userInfo := s.buildUserInfoFromWildcard(wildcardId, wildcardData)
		if userInfo != nil {
			panelInfo.User = append(panelInfo.User, userInfo)
		}
	}

	// 构建物品信息
	for _, wildcardId := range mailCfg.Text_wildcard {
		goodsInfo := s.buildGoodsInfoFromWildcard(wildcardId, wildcardData)
		if goodsInfo != nil {
			panelInfo.Goods = append(panelInfo.Goods, goodsInfo)
		}
	}

	return panelInfo, nil
}

// replaceWildcards 替换通配符
func (s *MailConfigService) replaceWildcards(text string, wildcardIds []int32, wildcardData map[string]interface{}) string {
	result := text

	for _, wildcardId := range wildcardIds {
		wildcardCfg := s.wildcardConfig.GetWildcardConfig(wildcardId)
		if wildcardCfg == nil {
			continue
		}

		// 替换通配符
		wildcardKey := fmt.Sprintf("{%d}", wildcardId)
		if wildcardCfg.Wildcard_text != "" {
			result = strings.ReplaceAll(result, wildcardKey, wildcardCfg.Wildcard_text)
		}
	}

	// 替换自定义通配符数据
	for key, value := range wildcardData {
		wildcardKey := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, wildcardKey, fmt.Sprintf("%v", value))
	}

	return result
}

// buildUserInfoFromWildcard 根据通配符构建用户信息
func (s *MailConfigService) buildUserInfoFromWildcard(wildcardId int32, wildcardData map[string]interface{}) *structsdef.MailUserInfo {
	wildcardCfg := s.wildcardConfig.GetWildcardConfig(wildcardId)
	if wildcardCfg == nil {
		return nil
	}

	userInfo := &structsdef.MailUserInfo{
		Key:  fmt.Sprintf("{%d}", wildcardId),
		Type: 0,
		Style: &structsdef.MailStyleInfo{
			Color:    wildcardCfg.Text_col,
			Font:     wildcardCfg.Text_font,
			FontS:    uint32(wildcardCfg.Text_size),
			FontW:    uint32(wildcardCfg.Text_weight),
			IconW:    uint32(wildcardCfg.Icon_w),
			IconH:    uint32(wildcardCfg.Icon_h),
			ShowName: wildcardCfg.Display_mode != 2, // 不是只显示头像
			ShowIcon: wildcardCfg.Display_mode != 1, // 不是只显示昵称
		},
	}

	// 从通配符数据中获取用户信息
	if userId, ok := wildcardData["user_id"]; ok {
		// TODO: 根据用户ID获取用户详细信息
		// userInfo.User = getUserInfo(userId.(uint64))
	}

	return userInfo
}

// buildGoodsInfoFromWildcard 根据通配符构建物品信息
func (s *MailConfigService) buildGoodsInfoFromWildcard(wildcardId int32, wildcardData map[string]interface{}) *structsdef.MailGoodsInfo {
	wildcardCfg := s.wildcardConfig.GetWildcardConfig(wildcardId)
	if wildcardCfg == nil {
		return nil
	}

	goodsInfo := &structsdef.MailGoodsInfo{
		Key: fmt.Sprintf("{%d}", wildcardId),
		Style: &structsdef.MailStyleInfo{
			Color:    wildcardCfg.Text_col,
			Font:     wildcardCfg.Text_font,
			FontS:    uint32(wildcardCfg.Text_size),
			FontW:    uint32(wildcardCfg.Text_weight),
			IconW:    uint32(wildcardCfg.Icon_w),
			IconH:    uint32(wildcardCfg.Icon_h),
			ShowName: wildcardCfg.Display_mode != 2,
			ShowIcon: wildcardCfg.Display_mode != 1,
		},
	}

	// 从通配符数据中获取物品信息
	if itemId, ok := wildcardData["item_id"]; ok {
		// TODO: 根据物品ID获取物品详细信息
		// goodsInfo.Item = getItemInfo(itemId.(int32))
	}

	return goodsInfo
}

// parseMailType 解析邮件类型
func (s *MailConfigService) parseMailType(mailType string) int32 {
	switch mailType {
	case "SYSTEM":
		return int32(constdef.SYSTEM)
	case "BATTLE":
		return int32(constdef.BATTLE)
	case "TEAM":
		return int32(constdef.TEAM)
	case "FAMILY":
		return int32(constdef.SWITCH_FAMILY)
	default:
		return int32(constdef.MAIL_COMMON)
	}
}

// parseMailSubType 解析邮件子类型
func (s *MailConfigService) parseMailSubType(typeDesc string) int32 {
	switch typeDesc {
	case "AWARD":
		return int32(constdef.AWARD)
	case "TEAM_INVITAE":
		return int32(constdef.TEAM_INVITAE)
	case "TEAM_JOIN":
		return int32(constdef.TEAM_JOIN)
	case "TEAM_QUIT":
		return int32(constdef.TEAM_QUIT)
	default:
		return int32(constdef.AWARD)
	}
}

// calculateExpireTime 计算过期时间
func (s *MailConfigService) calculateExpireTime(loseTime int32) int64 {
	if loseTime <= 0 {
		return 0 // 永不过期
	}

	// loseTime 是秒数，转换为时间戳
	return int64(loseTime)
}

// 全局邮件配置服务实例
var GlobalMailConfigService *MailConfigService

func init() {
	GlobalMailConfigService = NewMailConfigService()
}
