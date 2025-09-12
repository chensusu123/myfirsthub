package familyservice

import (
	"context"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"
	"maze_game_server/services/groupservice"

	"maze_game_server/app"

	"maze_game_server/common/constdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// CreateFamily 创建家族
func (r *service) CreateFamily(ctx context.Context, userID uint64, allianceID int32, familyName string,
	familySetting int32, userInfo familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 创建新家族
	familyInfoModel, err := familymodel.CreateFamilyInfoModel(ctx, familyName, familySetting)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily NewFamilyInfoModel failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting))
		return nil, err
	}
	// 加入成员
	familyInfoModel.AddMember(ctx, userInfo)

	//创建家族群聊
	groupInfo, err := groupservice.Default.CreateGroup(ctx, app.Maze, userID, constdef.GroupTypeFamily, make([]uint64, 0))
	if err != nil {
		logger.CtxError(ctx, "OnCreateFamilyRQ CreateGroup error", zap.Error(err))
		return nil, err
	}
	familyInfoModel.SetFamilyGroupID(ctx, groupInfo.ID)

	// //加入联盟群聊
	// err = grouppkg.InviteMember(ctx, app.Maze.ID(), int64(groupInfo.ID), userID)
	// if err != nil {
	// 	logger.CtxError(ctx, "CreateFamily InviteMember alliance err",
	// 		zap.Error(err))
	// 	return nil, err
	// }

	// 保存家族信息
	err = familyInfoModel.Save(ctx, familyInfoModel.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily model.Save failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting),
			zap.Error(err))
		return nil, err
	}

	// 把家族id加入家族列表
	familyListModel, err := familymodel.LoadFamilyListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily LoadFamilyListModel err",
			zap.Error(err))
		return nil, err
	}
	familyListModel.AddFamily(ctx, familyInfoModel.FamilyID)
	err = familyListModel.Save(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily familyListModel.Save err",
			zap.Error(err))
		return nil, err
	}

	// 设置家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, familyInfoModel.FamilyID)
	err = familyToAllianceModel.SetUserAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily familyToAllianceModel.SetUserAlliance err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	// 联盟中增加该家族id
	allianceModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID),
			zap.Error(err))
		return nil, err
	}

	err = allianceModel.AddFamilyID(ctx, familyInfoModel.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily allianceModel.AddFamilyID err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	return familyInfoModel, nil
}

// DeductCreateFamilyCost 创建家族扣物品
func (r *service) DeductCreateFamilyCost(ctx context.Context, cost map[int32]int64) error {
	// todo 待补充，读表

	return nil
}

// 查询创建家族消耗
func (r *service) QueryCreateFamilyCost(ctx context.Context) (map[int32]int64, error) {
	// todo 待补充，读表

	return nil, nil
}

// CheckCreateFamilyCost 检查创建家族消耗
func (r *service) CheckCreateFamilyCost(ctx context.Context) (bool, error) {
	// todo 待补充，读表

	return true, nil
}
