package familyservice

import (
	"context"
	"errors"

	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// OnKickFamilyRQ 踢出家族
func (r *service) KickFamily(ctx context.Context, userID uint64, familyID int32,
	kickUsers []uint64,
) error {
	logger := fklog.ContextAppLogger(ctx)
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "KickFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	// 检查用户权限
	userPrivilege, err := familyInfoModel.GetFamilyUserPrivilege(ctx, familyID, userID)
	if userPrivilege == 0 || userPrivilege != int32(MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_LEADER) {
		logger.CtxError(ctx, "KickFamily CheckPrivilegeLevel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("权限不足")
	}
	// 族长不能被踢
	if familyInfoModel.CheckHaveLeader(ctx, kickUsers) {
		logger.CtxError(ctx, "KickFamily CheckHaveLeader err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("族长不能被踢")
	}
	familyInfoModel.RemMember(ctx, kickUsers)
	return familyInfoModel.Save(ctx, familyID)
}

// SendKickFamilyID 踢出家族id包
func (r *service) SendKickFamilyID(ctx context.Context, userID uint64, familyID int32,
	kickUsers []uint64,
) error {
	logger := fklog.ContextAppLogger(ctx)
	familymodel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "SendKickFamilyID LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.KickFamilyID{
		FamilyId:       proto.Int32(familyID),
		MazeFamilyInfo: familymodel.DataToFamilyInfoPb(ctx),
		MemberList:     familymodel.DataToFamilyMembersPb(ctx),
	}

	familyMembers := familymodel.DataToFamilyMembersPb(ctx)
	for _, v := range familyMembers {
		online.ClusterPush(ctx, v.GetUserId(), 0, pack)
	}
	return nil
}

// ExitFamily 退出家族
func (r *service) ExitFamily(ctx context.Context, userID uint64, familyID int32) error {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	logger := fklog.ContextAppLogger(ctx)
	if err != nil {
		logger.CtxError(ctx, "ExitFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	if familyInfoModel.CheckHaveLeader(ctx, []uint64{userID}) {
		logger.CtxError(ctx, "ExitFamily CheckHaveLeader err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("族长不能退出")
	}
	familyInfoModel.RemMember(ctx, []uint64{userID})
	return familyInfoModel.Save(ctx, familyID)
}

// SendKickFamilyID 踢出家族id包
func (r *service) SendExitFamilyID(ctx context.Context, userID uint64, familyID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	familymodel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "SendExitFamilyID LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.ExitFamilyID{
		FamilyId:       proto.Int32(familyID),
		MazeFamilyInfo: familymodel.DataToFamilyInfoPb(ctx),
		MemberList:     familymodel.DataToFamilyMembersPb(ctx),
	}

	familyMembers := familymodel.DataToFamilyMembersPb(ctx)
	for _, v := range familyMembers {
		online.ClusterPush(ctx, v.GetUserId(), 0, pack)
	}
	return nil
}
