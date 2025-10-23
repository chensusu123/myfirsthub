package mailservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/pb/common/MazeCommon"

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

	// 从配置中获取邮件配置
	mailCfg, err := s.loadMailConfig(configId)
	if err != nil {
		logger.CtxError(ctx, "BuildMailFromConfig loadMailConfig failed",
			zap.Int32("configId", configId),
			zap.Error(err))
		return nil, err
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
		ID:              uint64(time.Now().UnixNano()),
		Title:           mailCfg["Mail_text"].(string),
		Content:         mailCfg["Mail_text"].(string),
		Label:           int32(constdef.MailLabelSystem),
		MailType:        s.parseMailType(mailCfg["Mail_type"].(string)),
		MailSubType:     s.parseMailSubType(mailCfg["Mail_sub_type"].(string)),
		MailContentType: int32(constdef.CONTENT_TEXT_AND_AWARD),
		Sender:          "系统",
		ReciverID:       userId,
		IsRead:          false,
		IsGetAttach:     false,
		SendTime:        time.Now().Unix(),
		ExpireTime:      s.calculateExpireTime(int32(mailCfg["Expire_hours"].(int))),
		ExtType:         1,
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
	if mailCfgMap, ok := mailCfg.(map[string]interface{}); ok && mailCfgMap["Mail_text_start"] != "" {
		content := s.replaceWildcards(mailCfgMap["Mail_text_start"].(string), []int32{}, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
			Content: content,
		})
	}

	// 构建邮件内容
	if mailCfgMap, ok := mailCfg.(map[string]interface{}); ok && mailCfgMap["Mail_text"] != "" {
		content := s.replaceWildcards(mailCfgMap["Mail_text"].(string), []int32{}, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
			Content: content,
		})
	}

	// 构建邮件结尾
	if mailCfgMap, ok := mailCfg.(map[string]interface{}); ok && mailCfgMap["Mail_text_end"] != "" {
		content := s.replaceWildcards(mailCfgMap["Mail_text_end"].(string), []int32{}, wildcardData)
		panelInfo.Contents = append(panelInfo.Contents, &structsdef.MailPanelContent{
			Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
			Content: content,
		})
	}

	// 构建用户信息
	if mailCfgMap, ok := mailCfg.(map[string]interface{}); ok {
		if startWildcards, exists := mailCfgMap["Text_start_wildcard"]; exists {
			if wildcardIds, ok := startWildcards.([]int32); ok {
				for _, wildcardId := range wildcardIds {
					userInfo := s.buildUserInfoFromWildcard(wildcardId, wildcardData)
					if userInfo != nil {
						panelInfo.User = append(panelInfo.User, userInfo)
					}
				}
			}
		}

		// 构建物品信息
		if contentWildcards, exists := mailCfgMap["Text_wildcard"]; exists {
			if wildcardIds, ok := contentWildcards.([]int32); ok {
				for _, wildcardId := range wildcardIds {
					goodsInfo := s.buildGoodsInfoFromWildcard(wildcardId, wildcardData)
					if goodsInfo != nil {
						panelInfo.Goods = append(panelInfo.Goods, goodsInfo)
					}
				}
			}
		}
	}

	return panelInfo, nil
}

// replaceWildcards 替换通配符
func (s *MailConfigService) replaceWildcards(text string, wildcardIds []int32, wildcardData map[string]interface{}) string {
	result := text

	for _, wildcardId := range wildcardIds {
		// 通配符配置加载
		wildcardCfg := s.getWildcardConfig(wildcardId)
		if wildcardCfg == nil {
			continue
		}

		// 替换通配符
		wildcardKey := fmt.Sprintf("{%d}", wildcardId)
		wildcardValue := s.getWildcardValue(wildcardId, wildcardData)
		result = strings.ReplaceAll(result, wildcardKey, wildcardValue)
	}

	// 替换自定义通配符数据
	for key, value := range wildcardData {
		wildcardKey := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, wildcardKey, fmt.Sprintf("%v", value))
	}

	return result
}

// getWildcardConfig 获取通配符配置
func (s *MailConfigService) getWildcardConfig(wildcardId int32) map[string]interface{} {
	// 尝试从配置系统加载
	// if s.wildcardConfig != nil {
	// 	config := s.wildcardConfig.GetWildcardConfig(wildcardId)
	// 	if config != nil {
	// 		return config
	// 	}
	// }

	// 如果配置系统不可用，使用默认配置
	defaultWildcardConfigs := map[int32]map[string]interface{}{
		1: { // 用户名称通配符
			"Display_mode":  1,
			"Wildcard_text": "用户",
			"Text_col":      "#000000",
			"Text_font":     "Arial",
			"Text_size":     12,
			"Text_weight":   400,
			"Icon_w":        32,
			"Icon_h":        32,
		},
		2: { // 物品名称通配符
			"Display_mode":  1,
			"Wildcard_text": "物品",
			"Text_col":      "#FFD700",
			"Text_font":     "Arial",
			"Text_size":     14,
			"Text_weight":   600,
			"Icon_w":        48,
			"Icon_h":        48,
		},
		3: { // 数值通配符
			"Display_mode":  1,
			"Wildcard_text": "数值",
			"Text_col":      "#00FF00",
			"Text_font":     "Arial",
			"Text_size":     16,
			"Text_weight":   700,
			"Icon_w":        24,
			"Icon_h":        24,
		},
		4: { // 时间通配符
			"Display_mode":  1,
			"Wildcard_text": "时间",
			"Text_col":      "#FF6B6B",
			"Text_font":     "Arial",
			"Text_size":     12,
			"Text_weight":   400,
			"Icon_w":        28,
			"Icon_h":        28,
		},
	}

	if config, exists := defaultWildcardConfigs[wildcardId]; exists {
		return config
	}

	// 如果找不到配置，返回默认配置
	return defaultWildcardConfigs[1]
}

// getWildcardValue 获取通配符值
func (s *MailConfigService) getWildcardValue(wildcardId int32, wildcardData map[string]interface{}) string {
	// 从通配符数据中获取值
	if value, exists := wildcardData[fmt.Sprintf("wildcard_%d", wildcardId)]; exists {
		return fmt.Sprintf("%v", value)
	}

	// 返回默认值
	return "通配符值"
}

// buildUserInfoFromWildcard 根据通配符构建用户信息
func (s *MailConfigService) buildUserInfoFromWildcard(wildcardId int32, wildcardData map[string]interface{}) *structsdef.MailUserInfo {
	// 通配符配置
	wildcardCfg := s.getWildcardConfig(wildcardId)
	if wildcardCfg == nil {
		return nil
	}

	userInfo := &structsdef.MailUserInfo{
		Key:  fmt.Sprintf("{%d}", wildcardId),
		Type: 0,
		Style: &structsdef.MailStyleInfo{
			Color:    wildcardCfg["Text_col"].(string),
			Font:     wildcardCfg["Text_font"].(string),
			FontS:    uint32(wildcardCfg["Text_size"].(int)),
			FontW:    uint32(wildcardCfg["Text_weight"].(int)),
			IconW:    uint32(wildcardCfg["Icon_w"].(int)),
			IconH:    uint32(wildcardCfg["Icon_h"].(int)),
			ShowName: wildcardCfg["Display_mode"].(int) != 2, // 不是只显示头像
			ShowIcon: wildcardCfg["Display_mode"].(int) != 1, // 不是只显示昵称
		},
	}

	// 从通配符数据中获取用户信息
	if userId, ok := wildcardData["user_id"]; ok {
		// 根据用户ID获取用户详细信息
		userInfo.User = s.getUserInfo(userId.(uint64))
	}

	return userInfo
}

// buildGoodsInfoFromWildcard 根据通配符构建物品信息
func (s *MailConfigService) buildGoodsInfoFromWildcard(wildcardId int32, wildcardData map[string]interface{}) *structsdef.MailGoodsInfo {
	// 通配符配置
	wildcardCfg := s.getWildcardConfig(wildcardId)
	if wildcardCfg == nil {
		return nil
	}

	goodsInfo := &structsdef.MailGoodsInfo{
		Key: fmt.Sprintf("{%d}", wildcardId),
		Style: &structsdef.MailStyleInfo{
			Color:    wildcardCfg["Text_col"].(string),
			Font:     wildcardCfg["Text_font"].(string),
			FontS:    uint32(wildcardCfg["Text_size"].(int)),
			FontW:    uint32(wildcardCfg["Text_weight"].(int)),
			IconW:    uint32(wildcardCfg["Icon_w"].(int)),
			IconH:    uint32(wildcardCfg["Icon_h"].(int)),
			ShowName: wildcardCfg["Display_mode"].(int) != 2,
			ShowIcon: wildcardCfg["Display_mode"].(int) != 1,
		},
	}

	// 从通配符数据中获取物品信息
	if itemId, ok := wildcardData["item_id"]; ok {
		// 根据物品ID获取物品详细信息
		goodsInfo.Item = s.getItemInfo(itemId.(int32))
	}

	return goodsInfo
}

// loadMailConfig 从配置中加载邮件配置
func (s *MailConfigService) loadMailConfig(configId int32) (map[string]interface{}, error) {
	// 尝试从配置系统加载
	// if s.mailInfoConfig != nil {
	// 	config := s.mailInfoConfig.GetMailInfoConfig(configId)
	// 	if config != nil {
	// 		return config, nil
	// 	}
	// }

	// 如果配置系统不可用，使用默认配置
	defaultConfigs := map[int32]map[string]interface{}{
		1: { // 系统邮件
			"Mail_text_start":     "亲爱的玩家",
			"Text_start_wildcard": "{{user_name}}",
			"Mail_text":           "您有一封新邮件",
			"Text_wildcard":       "{{content}}",
			"Mail_text_end":       "祝您游戏愉快！",
			"Text_end_wildcard":   "{{sender_name}}",
			"Mail_type":           "SYSTEM",
			"Mail_sub_type":       "AWARD",
			"Expire_hours":        24,
		},
		2: { // 奖励邮件
			"Mail_text_start":     "恭喜您获得奖励！",
			"Text_start_wildcard": "{{user_name}}",
			"Mail_text":           "您获得了以下奖励",
			"Text_wildcard":       "{{reward_list}}",
			"Mail_text_end":       "请及时领取！",
			"Text_end_wildcard":   "{{sender_name}}",
			"Mail_type":           "SYSTEM",
			"Mail_sub_type":       "AWARD",
			"Expire_hours":        72,
		},
		3: { // 活动邮件
			"Mail_text_start":     "活动通知",
			"Text_start_wildcard": "{{user_name}}",
			"Mail_text":           "新的活动开始了！",
			"Text_wildcard":       "{{activity_info}}",
			"Mail_text_end":       "快来参与吧！",
			"Text_end_wildcard":   "{{sender_name}}",
			"Mail_type":           "ACTIVITY",
			"Mail_sub_type":       "NOTICE",
			"Expire_hours":        168, // 7天
		},
	}

	if config, exists := defaultConfigs[configId]; exists {
		return config, nil
	}

	// 如果找不到配置，返回默认配置
	return defaultConfigs[1], nil
}

// getUserInfo 获取用户信息
func (s *MailConfigService) getUserInfo(userId uint64) *structsdef.SimpleUserInfo {
	// 临时返回默认用户信息
	return &structsdef.SimpleUserInfo{
		UserId:    userId,
		Nickname:  fmt.Sprintf("用户%d", userId),
		IconToken: 1,
		UserSex:   1,
	}
}

// getItemInfo 获取物品信息
func (s *MailConfigService) getItemInfo(itemId int32) *MazeCommon.MazeItem {
	// 临时返回默认物品信息
	count := int64(1)
	return &MazeCommon.MazeItem{
		ItemId: &itemId,
		Count:  &count,
	}
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

// SendMailWithConfig 根据配置发送邮件
func (s *MailConfigService) SendMailWithConfig(ctx context.Context, configId int32, userId uint64, wildcardData map[string]interface{}) error {
	logger := fklog.ContextAppLogger(ctx)

	// 根据配置构建邮件
	mailInfo, err := s.BuildMailFromConfig(ctx, configId, userId, wildcardData)
	if err != nil {
		logger.CtxError(ctx, "SendMailWithConfig BuildMailFromConfig failed",
			zap.Int32("configId", configId),
			zap.Uint64("userId", userId),
			zap.Error(err))
		return err
	}

	// 发送邮件 - 使用全局邮件服务
	err = GlobalMailService.SendMail(logger, mailInfo.Title, mailInfo.Content, mailInfo.Sender,
		mailInfo.Label, mailInfo.ReciverID, mailInfo.Attachments, mailInfo.ExpireTime)

	if err != nil {
		logger.CtxError(ctx, "SendMailWithConfig SendMail failed",
			zap.Int32("configId", configId),
			zap.Uint64("userId", userId),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "SendMailWithConfig success",
		zap.Int32("configId", configId),
		zap.Uint64("userId", userId),
		zap.String("title", mailInfo.Title))

	return nil
}

// 全局邮件配置服务实例
var GlobalMailConfigService *MailConfigService

func init() {
	GlobalMailConfigService = NewMailConfigService()
}
