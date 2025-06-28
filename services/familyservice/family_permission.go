package familyservice

import (
	"errors"
	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

// 修改家族某玩家权限
func (r *service) OperatePrivilege(logger fklog.FKLogI, userID uint64, familyID int32,
	privilegeLevel int32, operateUsers []uint64) error {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("OperatePrivilege LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	// 检查身份
	userPrivilege, err := familyInfoModel.GetFamilyUserPrivilege(logger, familyID, userID)
	if userPrivilege == 0 || userPrivilege != int32(MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_LEADER) {
		logger.ErrorWF("OperatePrivilege CheckPrivilegeLevel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return errors.New("权限不足")
	}

	// todo 读表检查家族该权限是否过多

	// 修改权限
	err = familyInfoModel.UpdatePrivilegeLevel(logger, familyID, operateUsers, privilegeLevel)
	if err != nil {
		logger.ErrorWF("OperatePrivilege UpdatePrivilegeLevel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	err = familyInfoModel.Save(logger, familyID)
	if err != nil {
		logger.ErrorWF("OperatePrivilege familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	return nil
}
