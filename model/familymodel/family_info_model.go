package familymodel

import (
	"context"
	"errors"
	"maze_game_server/io/redis/familyredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/MazeFamily"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/*
  optional int32 family_id = 1; // 家族id
  optional string family_name = 2; // 家族名字
  optional int32 family_level = 3; // 家族等级
  optional FamilySetting family_setting = 4; // 家族设置信息
  optional int32 member_count = 5; // 家族成员数量
  optional int32 member_count_limit = 6; // 家族成员数量上限
*/

var memberCountLimit int32 = 50 // 初步设置家族人数上限为固定值50

const (
	FamilyJoinTypeNone   = 1
	FamilyJoinTypeQRCode = 2
)

const (
	FamilyPrivilegeLevelLeader     = 10 // 族长
	FamilyPrivilegeLevelViceLeader = 20 // 副族长
	FamilyPrivilegeLevelManager    = 30 // 管理员
	FamilyPrivilegeLevelMember     = 40 // 普通成员
)

// 家族成员信息
type FamilyMember struct {
	UserID         uint64 `json:"user_id,omitempty"`
	NickName       string `json:"nick_name,omitempty"`
	Sex            int32  `json:"sex,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	Level          int32  `json:"level,omitempty"` // 这个字段可以先忽略，加pb的时候先定义出来了，但是暂时没用
	PrivilegeLevel int32  `json:"privilege_level,omitempty"`
	JoinTime       int64  `json:"join_time,omitempty"` // 加入家族时间
}

// 家族设置信息
type FamilySetting struct {
	FamilyJoinType int32 `json:"family_join_type,omitempty"`
}

// 家族信息模型
type FamilyInfoModel struct {
	FamilyID         int32           `json:"family_id,omitempty"`
	FamilyName       string          `json:"family_name,omitempty"`
	FamilyLevel      int32           `json:"family_level,omitempty"`
	FamilySetting    FamilySetting   `json:"family_setting,omitempty"`
	MemberCount      int32           `json:"member_count,omitempty"`
	MemberCountLimit int32           `json:"member_count_limit,omitempty"`
	FamilyMembers    []*FamilyMember `json:"family_members,omitempty"`
	FamilyApplyUsers []*FamilyMember `json:"family_apply_users,omitempty"`
	FamilyGroupID    int32           `json:"group_id,omitempty"`
}

type FamilysInfoModel []*FamilyInfoModel

func LoadFamilyInfoModel(ctx context.Context, familyID int32) (r *FamilyInfoModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	r = &FamilyInfoModel{}
	if err = r.load(ctx, familyID); err != nil {
		logger.CtxError(ctx, "LoadFamilyModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return
}

func LoadFamilyListInfoModel(ctx context.Context, familyIDs []int32) (r FamilysInfoModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	result, err := familyredis.BatchGetFamilyInfo(familyIDs)
	if err != nil {
		logger.CtxError(ctx, "BatchGetFamilyInfo err",
			zap.Int32s("familyIDs", familyIDs), zap.Error(err))
		return nil, err
	}
	for _, v := range result {
		if v == nil {
			continue
		}
		familyInfo := &FamilyInfoModel{}
		err = serialize.Unmarshal(v, familyInfo)
		if err != nil {
			logger.CtxError(ctx, "Unmarshal err",
				zap.Int32("familyID", familyInfo.FamilyID), zap.Error(err))
			continue
		}
		r = append(r, familyInfo)
	}

	return
}

// NewFamilyInfoModel 新创建家族
func NewFamilyInfoModel(ctx context.Context, familyName string, joinType int32) (*FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	if joinType == 0 && familyName == "" {
		return nil, errors.New("familyName or familySetting is nil")
	}
	familyID, err := familyredis.CreateFamilyId()
	if err != nil {
		logger.CtxError(ctx, "CreateFamilyId err",
			zap.String("familyName", familyName), zap.Any("joinType", joinType),
			zap.Error(err))
		return nil, err
	}
	r := &FamilyInfoModel{
		FamilyName: familyName,
		FamilyID:   familyID,
		FamilySetting: FamilySetting{
			FamilyJoinType: joinType,
		},
		MemberCountLimit: memberCountLimit,
	}
	return r, nil
}

func (r *FamilyInfoModel) AddMember(ctx context.Context, userInfo FamilyMember) {
	r.FamilyMembers = append(r.FamilyMembers, &FamilyMember{
		UserID:         userInfo.UserID,
		NickName:       userInfo.NickName,
		Avatar:         userInfo.Avatar,
		Sex:            userInfo.Sex,
		Level:          userInfo.Level,
		PrivilegeLevel: userInfo.PrivilegeLevel,
		JoinTime:       time.Now().Unix(),
	})
	r.MemberCount = int32(len(r.FamilyMembers))
}

// CheckRemMember 检查成员退出
func (r *FamilyInfoModel) CheckRemMember(ctx context.Context) error {
	if len(r.FamilyMembers) > 0 && len(r.FamilyApplyUsers) == 1 {
		return errors.New("rem member limit, Keep at least one person in the family")
	}
	return nil
}

// RemMember 批量移除成员
func (r *FamilyInfoModel) RemMember(ctx context.Context, userIDs []uint64) {
	newMembers := make([]*FamilyMember, 0, len(r.FamilyMembers))
	mapUsers := r.mapUserList(userIDs)
	for _, member := range r.FamilyMembers {
		if _, ok := mapUsers[member.UserID]; !ok {
			newMembers = append(newMembers, member)
		}
	}
	r.FamilyMembers = newMembers
	r.MemberCount = int32(len(r.FamilyMembers))
}

func (r *FamilyInfoModel) AddApplyUser(ctx context.Context, userInfo FamilyMember) {
	r.FamilyApplyUsers = append(r.FamilyApplyUsers, &FamilyMember{
		UserID:         userInfo.UserID,
		NickName:       userInfo.NickName,
		Avatar:         userInfo.Avatar,
		Sex:            userInfo.Sex,
		Level:          userInfo.Level,
		PrivilegeLevel: userInfo.PrivilegeLevel,
	})
}

// RemApplyUser 去掉申请者
func (r *FamilyInfoModel) RemApplyUser(ctx context.Context, userInfo FamilyMember) {
	newApplyUsers := make([]*FamilyMember, 0, len(r.FamilyApplyUsers))
	mapUsers := r.mapUserList([]uint64{userInfo.UserID})
	for _, applyUser := range r.FamilyApplyUsers {
		if _, ok := mapUsers[applyUser.UserID]; !ok {
			newApplyUsers = append(newApplyUsers, applyUser)
		}
	}
	r.FamilyApplyUsers = newApplyUsers
}

func (r *FamilyInfoModel) load(ctx context.Context, familyID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := familyredis.GetFamilyInfo(familyID)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	return
}

func (r *FamilyInfoModel) Save(ctx context.Context, familyID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.CtxError(ctx, "save Marshal err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	err = familyredis.SetFamilyInfo(familyID, value)
	if err != nil {
		logger.CtxError(ctx, "save SetFamilyInfo err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	return
}

// Delete 解散家族
func (r *FamilyInfoModel) Delete(ctx context.Context, familyID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	err = familyredis.DelFamilyInfo(familyID)
	if err != nil {
		logger.CtxError(ctx, "delete DelFamilyInfo err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	// 删除家族所有成员的绑定家族关系
	for _, member := range r.FamilyMembers {
		err = familyredis.DelUserFamilyID(member.UserID)
		if err != nil {
			logger.CtxError(ctx, "delete DelUserFamilyID err",
				zap.Uint64("userID", member.UserID), zap.Error(err))
		}
	}

	return
}

// DataToFamilyInfoPb 把家族信息填进pb
func (r *FamilyInfoModel) DataToFamilyInfoPb(ctx context.Context) *MazeFamily.FamilyInfo {
	return &MazeFamily.FamilyInfo{
		FamilyId:         proto.Int32(r.FamilyID),
		FamilyName:       proto.String(r.FamilyName),
		FamilyLevel:      proto.Int32(r.FamilyLevel),
		MemberCount:      proto.Int32(r.MemberCount),
		MemberCountLimit: proto.Int32(r.MemberCountLimit),
		FamilySetting: &MazeFamily.FamilySetting{
			FamilyJoinType: MazeFamily.FamilyJoinType(r.FamilySetting.FamilyJoinType).Enum(),
		},
		GroupId: proto.Int32(r.FamilyGroupID),
	}
}

// DataToFamilyInfoPb 把家族列表信息填进pb
func (r FamilysInfoModel) DataToFamilyListInfoPb(ctx context.Context) []*MazeFamily.FamilyInfo {
	ret := make([]*MazeFamily.FamilyInfo, 0)
	for _, v := range r {
		ret = append(ret, v.DataToFamilyInfoPb(ctx))
	}
	return ret
}

// DataToFamilyMembersPb 取出所有家族成员信息
func (r *FamilyInfoModel) DataToFamilyMembersPb(ctx context.Context) []*MazeFamily.FamilyMemberInfo {
	members := make([]*MazeFamily.FamilyMemberInfo, 0)
	for _, member := range r.FamilyMembers {
		members = append(members, &MazeFamily.FamilyMemberInfo{
			UserId:         proto.Uint64(member.UserID),
			NickName:       proto.String(member.NickName),
			Sex:            proto.Int32(member.Sex),
			Avatar:         proto.String(member.Avatar),
			Level:          proto.Int32(member.Level),
			PrivilegeLevel: MazeFamily.PrivilegeLevel(member.PrivilegeLevel).Enum(),
		})
	}
	return members
}

// DataToApplyUsersPb 取出所有申请用户信息
func (r *FamilyInfoModel) DataToApplyUsersPb(ctx context.Context) []*MazeFamily.FamilyMemberInfo {
	members := make([]*MazeFamily.FamilyMemberInfo, 0)
	for _, applyUser := range r.FamilyApplyUsers {
		members = append(members, &MazeFamily.FamilyMemberInfo{
			UserId:         proto.Uint64(applyUser.UserID),
			NickName:       proto.String(applyUser.NickName),
			Sex:            proto.Int32(applyUser.Sex),
			Avatar:         proto.String(applyUser.Avatar),
			Level:          proto.Int32(applyUser.Level),
			PrivilegeLevel: MazeFamily.PrivilegeLevel(applyUser.PrivilegeLevel).Enum(),
		})
	}
	return members
}

// DataToFamilyUserPb 取出指定家族用户信息
func (r *FamilyInfoModel) DataToFamilyUserPb(ctx context.Context, userID uint64) (ret *MazeFamily.FamilyMemberInfo) {
	for _, member := range r.FamilyMembers {
		if member.UserID != userID {
			continue
		}
		ret = &MazeFamily.FamilyMemberInfo{
			UserId:         proto.Uint64(member.UserID),
			NickName:       proto.String(member.NickName),
			Sex:            proto.Int32(member.Sex),
			Avatar:         proto.String(member.Avatar),
			Level:          proto.Int32(member.Level),
			PrivilegeLevel: MazeFamily.PrivilegeLevel(member.PrivilegeLevel).Enum(),
		}
	}
	return
}

// 家族升级
func (r *FamilyInfoModel) Upgrade(ctx context.Context, familyID int32, targetLevel int32) (err error) {
	if targetLevel <= r.FamilyLevel {
		return errors.New("targetLevel <= r.FamilyLevel")
	}
	r.FamilyLevel = targetLevel
	return
}

// UpdatePrivilegeLevel 更新家族权限
func (r *FamilyInfoModel) UpdatePrivilegeLevel(ctx context.Context, familyID int32, users []uint64, level int32) (err error) {
	userMap := r.mapFamilyMembers()
	for _, userID := range users {
		if _, ok := userMap[userID]; ok {
			userMap[userID].PrivilegeLevel = level
		}
	}
	return
}

func (r *FamilyInfoModel) mapFamilyMembers() map[uint64]*FamilyMember {
	userMap := make(map[uint64]*FamilyMember)
	for _, v := range r.FamilyMembers {
		userMap[v.UserID] = v
	}
	return userMap
}

// 去重
func (r *FamilyInfoModel) mapUserList(userList []uint64) map[uint64]struct{} {
	userMap := make(map[uint64]struct{})
	for _, v := range userList {
		userMap[v] = struct{}{}
	}
	return userMap
}

// 返回用户家族权限
func (r *FamilyInfoModel) GetFamilyUserPrivilege(ctx context.Context, familyID int32, userID uint64) (int32, error) {
	for _, v := range r.FamilyMembers {
		if v.UserID == userID {
			return v.PrivilegeLevel, nil
		}
	}
	return 0, errors.New("not found user")
}

// 检查用户是否为族长
func (r *FamilyInfoModel) CheckHaveLeader(ctx context.Context, userIDs []uint64) bool {
	userMap := r.mapUserList(userIDs)
	for _, v := range r.FamilyMembers {
		if _, ok := userMap[v.UserID]; ok {
			if v.PrivilegeLevel == int32(MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_LEADER) {
				return true
			}
		}
	}
	return false
}

// 校验用户是否在家族
func (r *FamilyInfoModel) CheckUserInFamily(ctx context.Context, userID uint64) bool {
	logger := fklog.ContextAppLogger(ctx)
	familyID, err := familyredis.GetUserFamilyID(userID)
	if err != nil {
		logger.CtxError(ctx, "CheckUserInFamily GetUserFamilyID err",
			zap.Uint64("userID", userID), zap.Error(err))
		return false
	}
	return familyID == r.FamilyID
}

// 检查家族人数
func (r *FamilyInfoModel) CheckFamilyMemberCount(ctx context.Context) error {
	if len(r.FamilyMembers) >= int(r.MemberCountLimit) {
		return errors.New("家族人数已满")
	}
	return nil
}

// 检查用户是否在家族的请求列表中
func (r *FamilyInfoModel) CheckUserInApplyList(ctx context.Context, userID uint64) bool {
	for _, v := range r.FamilyApplyUsers {
		if v.UserID == userID {
			return true
		}
	}
	return false
}

// 设置家族群组ID
func (r *FamilyInfoModel) SetFamilyGroupID(ctx context.Context, groupID int32) error {
	if r.FamilyGroupID != 0 {
		return errors.New("family group id has set")
	}
	r.FamilyGroupID = groupID
	return nil
}
