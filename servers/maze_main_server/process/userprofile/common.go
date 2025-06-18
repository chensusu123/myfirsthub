// @Author pangchenyang 2025/6/17 11:53:00
// @Desc: 
package userprofile

import (
	"fmt"
)

// 类型转换 mysql结构转为pb结构
// func convertToProfile(dbProfile *userprofilemysql.UserProfile) *UserProfile.UserProfile {
// 	return &UserProfile.UserProfile{
// 		UserId:   proto.Uint64(dbProfile.UserID),
// 		NickName: proto.String(dbProfile.NickName),
// 		Avatar:   proto.String(dbProfile.Avatar),
// 		Sex:      proto.Int32(int32(dbProfile.Sex)),
// 	}
// }
//
// // 类型转换 pb结构转为mysql结构
// func convertToDbProfile(profile *UserProfile.UserProfile, creatTime, updateTime time.Time) *userprofilemysql.UserProfile {
// 	if profile == nil {
// 		return nil
// 	}
// 	return &userprofilemysql.UserProfile{
// 		UserID:    profile.GetUserId(),
// 		NickName:  profile.GetNickName(),
// 		Avatar:    profile.GetAvatar(),
// 		Sex:       uint8(profile.GetSex()),
// 		CreatedAt: creatTime,
// 		UpdatedAt: updateTime,
// 	}
// }

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
