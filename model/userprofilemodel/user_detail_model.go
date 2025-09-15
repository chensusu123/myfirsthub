package userprofilemodel

//用于当前用户查询其他人的详细信息
type UserDetailModel struct {
	UserID    int64  `json:"user_id"`
	NickName  string `json:"nick_name"`
	Sex       int32  `json:"sex"`
	IconToken string `json:"icon_token"`
	IsBlack   bool   `json:"is_black"`
	IsFriend  bool   `json:"is_friend"`
	ShowID    int32  `json:"show_id"`
}
