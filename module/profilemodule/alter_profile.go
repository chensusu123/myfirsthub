// @Author pangchenyang 2025/6/19 15:24:00
// @Desc: 
package profilemodule

import (
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/io/mysql/flowrecord"
	"time"
	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"errors"
	"fmt"
)

func UpdateUserProfile(logger fklog.FKLogI, data *UserProfile.UserProfile) (*UserProfile.UserProfile, error) {
	if data == nil || data.GetUserId() <= 0 {
		return nil, errors.New("无效的修改信息")
	}
	userID := data.GetUserId()
	var ret = &UserProfile.UserProfile{}

	dbProfile, err := userprofileredis.GetProfile(userID)
	if err != nil {
		return ret, err
	}

	var (
		chgDesc, newVal, oldVal string
	)

	// 修改用户资料
	alterRet := alterProfile(data, dbProfile, &chgDesc, &newVal, &oldVal)

	// 修改资料流水
	flowrecord.SaveAlterProfileRecord(logger, flowrecord.AlterProfileRecord{
		UserId:     userID,
		NewVal:     newVal,
		OldVal:     oldVal,
		ChgDesc:    chgDesc,
		CreateTime: time.Now().UnixMilli(), // ms
	})
	ret = alterRet
	return ret, nil
}

// 修改需要更改的资料，支持用户资料为空时的初始化
func alterProfile(alterProfile, dbProfile *UserProfile.UserProfile,
	chgDesc *string, newVal *string, oldVal *string) *UserProfile.UserProfile {
	ret := proto.Clone(dbProfile).(*UserProfile.UserProfile)
	if alterProfile.GetNickName() != "" && alterProfile.GetNickName() != ret.GetNickName() {
		ret.NickName = alterProfile.NickName
		*chgDesc += chgDescMap[chgName]
		*newVal += getAlterVal(chgName, alterProfile.GetNickName())
		*oldVal += getAlterVal(chgName, ret.GetNickName())
	}
	if alterProfile.GetAvatar() != "" {
		ret.Avatar = alterProfile.Avatar
		*chgDesc += chgDescMap[chgAvatar]
		*newVal += getAlterVal(chgAvatar, alterProfile.GetAvatar())
		*oldVal += getAlterVal(chgAvatar, ret.GetAvatar())
	}
	if alterProfile.GetSex() != 0 {
		ret.Sex = alterProfile.Sex
		*chgDesc += chgDescMap[chgSex]
		*newVal += getAlterVal(chgSex, alterProfile.GetSex())
		*oldVal += getAlterVal(chgSex, ret.GetSex())
	}
	return ret
}

const (
	chgName = 1 + iota
	chgAvatar
	chgSex
)

var chgDescMap = map[int32]string{
	1: "修改昵称",
	2: "修改头像",
	3: "修改性别",
}

var descMap = map[int32]string{
	1: "昵称",
	2: "头像",
	3: "性别",
}

func getAlterVal(chgType int32, val interface{}) string {
	return fmt.Sprintf("%s:%s", descMap[chgType], val)
}
