package familyservice

import (
	"errors"
	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// OnKickFamilyRQ 踢出家族
func (r *service) KickFamily(logger fklog.FKLogI, userID uint64, familyID int32,
	kickUsers []uint64) error {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("KickFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	// 检查用户权限
	userPrivilege, err := familyInfoModel.GetFamilyUserPrivilege(logger, familyID, userID)
	if userPrivilege == 0 || userPrivilege != int32(MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_LEADER) {
		logger.ErrorWF("KickFamily CheckPrivilegeLevel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("权限不足")
	}
	// 族长不能被踢
	if familyInfoModel.CheckHaveLeader(logger, kickUsers) {
		logger.ErrorWF("KickFamily CheckHaveLeader err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("族长不能被踢")
	}
	familyInfoModel.RemMember(logger, kickUsers)
	return familyInfoModel.Save(logger, familyID)
}

// SendKickFamilyID 踢出家族id包
func (r *service) SendKickFamilyID(logger fklog.FKLogI, userID uint64, familyID int32,
	kickUsers []uint64) error {
	familymodel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("SendKickFamilyID LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.KickFamilyID{
		FamilyId:       proto.Int32(familyID),
		MazeFamilyInfo: familymodel.DataToFamilyInfoPb(logger),
		MemberList:     familymodel.DataToFamilyMembersPb(logger),
	}
	// todo 包id统一换掉
	familyMembers := familymodel.DataToFamilyMembersPb(logger)
	for _, v := range familyMembers {
		online.Push(logger, v.GetUserId(), 0, pack)
	}
	return nil
}

// ExitFamily 退出家族
func (r *service) ExitFamily(logger fklog.FKLogI, userID uint64, familyID int32) error {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("ExitFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	if familyInfoModel.CheckHaveLeader(logger, []uint64{userID}) {
		logger.ErrorWF("ExitFamily CheckHaveLeader err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("族长不能退出")
	}
	familyInfoModel.RemMember(logger, []uint64{userID})
	return familyInfoModel.Save(logger, familyID)
}

// SendKickFamilyID 踢出家族id包
func (r *service) SendExitFamilyID(logger fklog.FKLogI, userID uint64, familyID int32) error {
	familymodel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("SendExitFamilyID LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.ExitFamilyID{
		FamilyId:       proto.Int32(familyID),
		MazeFamilyInfo: familymodel.DataToFamilyInfoPb(logger),
		MemberList:     familymodel.DataToFamilyMembersPb(logger),
	}
	// todo 包id统一换掉
	familyMembers := familymodel.DataToFamilyMembersPb(logger)
	for _, v := range familyMembers {
		online.Push(logger, v.GetUserId(), 0, pack)
	}
	return nil
}
