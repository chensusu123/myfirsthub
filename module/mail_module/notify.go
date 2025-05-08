/*
 * Created: 2020-06-24 17:27 +0800
 *
 * Modified: 2020-07-06 10:44 +0800
 *
 * Description: 奖励发送离线推送
 *
 * Author: libinbin
 */
package mail_module

import (
	"context"

	"gitlab.ifreetalk.com/plate/excel/auto/GMailPushTextCfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/PushNotifyMsgQueue"
	"gitlab.ifreetalk.com/plate/protodef/NotifyMsgDef"
)

// 根据用户离线10小时进行推送（时间待定，可配的）
// 当满足离线时间时判断是否在大的CD内（2小时），如果在CD中将不推送
// customText 自定义推送消息，如果为空，使用mail_push_text表获取
func OfflineNotify(ctx context.Context, logger fklog.FKLogI, userId uint64, extType int32, customText string) int {
	pushType := int32(NotifyMsgDef.LocalPushType_LocalPushTypeMailAward)
	if customText == "" {
		if extType == 0 {
			return 0
		}

		cfg := GMailPushTextCfg.GetMailPushTextConfig(extType)
		if cfg == nil {
			return 0
		}
		if len(cfg.Text) < 1 {
			return 0
		}
		customText = cfg.Text

		pushType = cfg.Push_type
		if pushType <= 0 {
			return 0
		}

	} else {
		if extType != 0 {
			cfg := GMailPushTextCfg.GetMailPushTextConfig(extType)
			if cfg != nil && cfg.Push_type > 0 {
				pushType = cfg.Push_type
			}
		}
	}

	msg := &NotifyMsgDef.PushNotifyMsgID{}
	msg.ServerName = proto.String("mail_server")
	msg.Type = proto.Int32(pushType)
	msg.UserId = proto.Uint64(userId)
	msg.Content = proto.String(customText)
	msg.ExtraInfo = proto.String("{\"jump_to\":1001}")
	PushNotifyMsgQueue.PushMsg(fkserver.NewUserContext(ctx, userId, logger), msg)
	return 1
}
