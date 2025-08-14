package mail

import (
	"google.golang.org/protobuf/proto"
	"maze_game_server/model/mailmodel"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeMail"
)

func PbMailData(info *mailmodel.MailInfo) *MazeMail.MailData {
	attachments := make([]*MazeMail.Attachment, len(info.Attachments))
	for _, attachment := range info.Attachments {
		itemType := MazeCommon.MazeItemType_ENUM_GOODS
		if attachment.Extra == "equip" {
			itemType = MazeCommon.MazeItemType_ENUM_EQUIP
		}

		attachments = append(attachments, &MazeMail.Attachment{
			Item: &MazeCommon.MazeItem{
				ItemId: proto.Int32(attachment.ItemID),
				Count:  proto.Int64(attachment.Count),
			},
			ItemType: &itemType,
		})
	}

	return &MazeMail.MailData{
		MailId:      proto.Uint64(info.ID),
		Title:       proto.String(info.Title),
		Content:     proto.String(info.Content),
		Label:       proto.Int32(info.Label),
		SenderName:  proto.String(info.Sender),
		ReciverId:   proto.Uint64(info.ReciverID),
		IsRead:      proto.Bool(info.IsRead),
		IsGetAttach: proto.Bool(info.IsGetAttach),
		SendTime:    proto.Int64(info.SendTime),
		ExpireTime:  proto.Int64(info.ExpireTime),
		Attachments: attachments,
	}
}
