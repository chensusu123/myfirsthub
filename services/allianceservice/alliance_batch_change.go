package allianceservice

import (
	"context"
	"errors"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"
	"maze_game_server/services/familyservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) BatchChangeAlliance(ctx context.Context, userIDs []uint64, allianceID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "BatchChangeAlliance Start",
		zap.Any("userIDs", userIDs),
		zap.Any("allianceID", allianceID),
	)
	defer func() {
		logger.CtxInfo(ctx, "BatchChangeAlliance End",
			zap.Any("userIDs", userIDs),
			zap.Any("allianceID", allianceID),
		)
	}()

	// 获取联盟信息
	allianceInfo, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "BatchChangeAlliance Fail",
			zap.Error(err),
			zap.Any("userIDs", userIDs),
			zap.Any("alliance", allianceID),
		)
		return err
	}

	if allianceInfo.AllianceID == 0 {
		logger.CtxError(ctx, "BatchChangeAlliance Fail",
			zap.Error(err),
			zap.Any("userIDs", userIDs),
			zap.Any("alliance", allianceID),
		)
		return errors.New("联盟信息不存在")
	}

	defaultFamilyID := allianceInfo.DefaultFamilyID
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, defaultFamilyID)
	if err != nil {
		logger.CtxError(ctx, "BatchChangeAlliance LoadFamilyInfoModel",
			zap.Any("userIDs", userIDs),
			zap.Any("allianceID", allianceID),
			zap.Any("defaultFamilyID", defaultFamilyID),
			zap.Error(err),
		)
		return err
	}

	for _, userID := range userIDs {
		// 先获取当前用户家族
		nowFamilyID, err := familyservice.GlobalFamilyService.GetUserFamily(ctx, userID)
		if err != nil {
			logger.CtxError(ctx, "BatchChangeAlliance GetUserFamily Fail",
				zap.Any("userID", userID),
				zap.Any("nowFamilyID", nowFamilyID),
			)
			return err
		}
		if nowFamilyID != 0 {
			// 让当前用户退出原来家族
			err = familyservice.GlobalFamilyService.ExitFamily(ctx, userID, nowFamilyID)
			if err != nil {
				logger.CtxError(ctx, "BatchChangeAlliance ExitFamily Fail",
					zap.Any("userID", userID),
					zap.Any("nowFamilyID", nowFamilyID),
				)
				return err
			}
		}

		familyInfoModel.AddApplyUser(ctx, familymodel.FamilyMember{
			UserID:         userID,
			NickName:       "",
			Sex:            1,
			Avatar:         "",
			Level:          1,
			PrivilegeLevel: familymodel.FamilyPrivilegeLevelMember,
			JoinTime:       time.Now().UnixMilli(),
		})
		err = familyservice.GlobalFamilyService.SetUserFamily(ctx, userID, defaultFamilyID)
		if err != nil {
			logger.CtxError(ctx, "BatchChangeAlliance SetUserFamily Fail",
				zap.Any("userID", userID),
				zap.Any("allianceID", allianceID),
			)
			return err
		}
	}

	return familyInfoModel.Save(ctx, familyInfoModel.FamilyID)
}
