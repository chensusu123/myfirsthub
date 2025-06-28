package family

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/services/familyservice"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Family struct {
	component.Base
}

func NewFamily() *Family {
	return &Family{}
}

// OnGetCreateFamilyCostRQ 获取创建家族消耗
func (f *Family) OnGetCreateFamilyCostRQ_10582_10583(s *session.Session, req *MazeFamily.GetCreateFamilyCostRQ) (err error) {
	logger := log.Clone("OnGetCreateFamilyCostRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetCreateFamilyCostRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetCreateFamilyCostRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	familyCost, err := familyservice.GlobalFamilyService.QueryCreateFamilyCost(logger)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("查询创建家族消耗失败")
		return
	}

	// todo 根据costID 查询物品
	for costID, costNum := range familyCost {
		res.CostList = append(res.CostList, &MazeCommon.MazeItem{
			ItemId: proto.Int32(costID),
			Count:  proto.Int64(costNum),
		})
	}

	return
}

// OnCreateFamilyRQ 创建家族
func (f *Family) OnCreateFamilyRQ_10585_10586(s *session.Session, req *MazeFamily.CreateFamilyRQ) (err error) {
	logger := log.Clone("OnCreateFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.CreateFamilyRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnCreateFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyName() == "" {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("无效参数")
		return
	}
	if req.GetJoinType() != MazeFamily.FamilyJoinType_FAMILY_JOIN_TYPE_NONE &&
		req.GetJoinType() != MazeFamily.FamilyJoinType_FAMILY_JOIN_TYPE_QR_CODE {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("无效的加入方式")
		return
	}

	// 检查用户是否在某个家族中
	ok, err := familyservice.GlobalFamilyService.CheckUserHaveFamily(logger, uid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查用户是否在某个家族中失败")
		return
	}
	if ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("用户已经在某个家族中")
		return
	}

	// 检查玩家距离上次退出家族时间是否超过一天
	ok, err = familyservice.GlobalFamilyService.CheckUserLastLeaveFamilyTime(logger, uid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查玩家距离上次退出家族时间是否超过一天失败")
		return
	}
	if !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("玩家距离上次退出家族时间未超过一天")
		return
	}

	// 检查消耗
	ok, err = familyservice.GlobalFamilyService.CheckCreateFamilyCost(logger)
	if err != nil || !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查创建家族消耗失败")
		return
	}

	costList := make(map[int32]int64)
	for _, cost := range req.GetCostList() {
		costList[cost.GetItemId()] = cost.GetCount()
	}

	// 扣物品
	err = familyservice.GlobalFamilyService.DeductCreateFamilyCost(logger, costList)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
		return
	}
	// 创建家族
	familyInfo, err := familyservice.GlobalFamilyService.CreateFamily(logger, uid, req.GetAllianceId(), req.GetFamilyName(),
		int32(*req.JoinType.Enum()), familymodel.FamilyMember{
			UserID:         uid,
			NickName:       req.GetLeader().GetNickName(),
			Avatar:         req.GetLeader().GetAvatar(),
			Sex:            int32(req.GetLeader().GetSex()),
			Level:          int32(req.GetLeader().GetLevel()),
			PrivilegeLevel: int32(req.GetLeader().GetPrivilegeLevel()),
		})
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("创建家族失败")
		return
	}

	res.FamilyInfo = familyInfo.DataToFamilyInfoPb(logger)
	res.MemberList = familyInfo.DataToFamilyMembersPb(logger)

	// 设置玩家所在家族
	err = familyservice.GlobalFamilyService.SetUserFamily(logger, uid, res.FamilyInfo.GetFamilyId())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("设置玩家所在家族失败")
		return
	}

	return
}

// OnGetUpgradeFamilyCostRQ 获取升级家族消耗
func (f *Family) OnGetUpgradeFamilyCostRQ_10574_10575(s *session.Session, req *MazeFamily.GetUpgradeFamilyCostRQ) (err error) {
	logger := log.Clone("OnGetUpgradeFamilyCostRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetUpgradeFamilyCostRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId
	res.TargetLevel = req.TargetLevel

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetUpgradeFamilyCostRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 || req.GetTargetLevel() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 检查升级消耗
	familyCost, err := familyservice.GlobalFamilyService.QueryUpgradeFamilyCost(logger, req.GetFamilyId(), req.GetTargetLevel())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("查询升级家族消耗失败")
		return
	}

	for costID, costNum := range familyCost {
		res.CostList = append(res.CostList, &MazeCommon.MazeItem{
			ItemId: proto.Int32(costID),
			Count:  proto.Int64(costNum),
		})
	}

	return
}

// OnUpgradeFamilyRQ 升级家族
func (f *Family) OnUpgradeFamilyRQ_10576_10577(s *session.Session, req *MazeFamily.UpgradeFamilyRQ) (err error) {
	logger := log.Clone("OnUpgradeFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.UpgradeFamilyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnUpgradeFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 || req.GetTargetLevel() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// 检查升级消耗
	ok, err := familyservice.GlobalFamilyService.CheckUpgradeFamilyCost(logger, req.GetTargetLevel())
	if err != nil || !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查升级消耗失败")
		return
	}
	// 扣物品
	costList := make(map[int32]int64)
	for _, cost := range req.GetCostList() {
		costList[cost.GetItemId()] = cost.GetCount()
	}
	err = familyservice.GlobalFamilyService.DeductUpgradeFamilyCost(logger, costList)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
		return
	}
	// 升级家族
	familyInfo, err := familyservice.GlobalFamilyService.UpgradeFamily(logger, req.GetFamilyId(), req.GetTargetLevel())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("升级家族失败")
		return
	}

	res.FamilyInfo = familyInfo.DataToFamilyInfoPb(logger)

	_ = familyservice.GlobalFamilyService.SendUpgradeFamilyIDPack(logger, req.GetFamilyId())
	return
}

// OnGetFamilyListRQ 获取家族列表
func (f *Family) OnGetFamilyListRQ_10578_10579(s *session.Session, req *MazeFamily.GetFamilyListRQ) (err error) {
	logger := log.Clone("OnGetFamilyListRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetFamilyListRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetFamilyListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// todo 后续根据业务需求 可能需要给列表排序一下 比如说根据火热程度 等
	familysInfo, err := familyservice.GlobalFamilyService.GetFamilyList(logger)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取家族列表失败")
		return
	}

	for _, familyInfo := range familysInfo {
		res.FamilyInfo = append(res.FamilyInfo, familyInfo.DataToFamilyInfoPb(logger))
	}
	return
}

// OnGetFamilyInfoRQ 获取家族信息
func (f *Family) OnGetFamilyInfoRQ_10580_10581(s *session.Session, req *MazeFamily.GetFamilyInfoRQ) (err error) {
	logger := log.Clone("OnGetFamilyInfoRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetFamilyInfoRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetFamilyInfoRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	familyInfo, err := familyservice.GlobalFamilyService.GetFamilyInfo(logger, req.GetFamilyId())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取家族信息失败")
		return
	}

	res.MazeFamilyInfo = familyInfo.DataToFamilyInfoPb(logger)
	res.MemberList = familyInfo.DataToFamilyMembersPb(logger)
	res.ApplyList = familyInfo.DataToApplyUsersPb(logger)
	return
}

// OnApplyFamilyRQ 请求加入家族
func (f *Family) OnApplyFamilyRQ_10587_10588(s *session.Session, req *MazeFamily.ApplyFamilyRQ) (err error) {
	logger := log.Clone("OnApplyFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.ApplyFamilyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId
	res.ApplyUser = req.ApplyUser

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnApplyFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// todo 检测人数是否达到上限

	// 检查玩家距离上次退出家族时间是否超过一天
	ok, err := familyservice.GlobalFamilyService.CheckUserLastLeaveFamilyTime(logger, uid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查玩家距离上次退出家族时间是否超过一天失败")
		return
	}
	if !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("玩家距离上次退出家族时间未超过一天")
		return
	}

	ok, err = familyservice.GlobalFamilyService.CheckApplyFamily(logger, req.GetFamilyId(), &familymodel.FamilyMember{
		UserID:         uid,
		NickName:       req.GetApplyUser().GetNickName(),
		Avatar:         req.GetApplyUser().GetAvatar(),
		Sex:            req.GetApplyUser().GetSex(),
		Level:          req.GetApplyUser().GetLevel(),
		PrivilegeLevel: int32(req.GetApplyUser().GetPrivilegeLevel()),
	})
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查加入家族失败")
		return
	}
	if !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("检查加入家族失败")
		return
	}

	// 申请加入家族
	familyInfo, err := familyservice.GlobalFamilyService.ApplyFamily(logger, req.GetFamilyId(), familymodel.FamilyMember{
		UserID:         uid,
		NickName:       req.GetApplyUser().GetNickName(),
		Avatar:         req.GetApplyUser().GetAvatar(),
		Sex:            req.GetApplyUser().GetSex(),
		Level:          req.GetApplyUser().GetLevel(),
		PrivilegeLevel: int32(req.GetApplyUser().GetPrivilegeLevel()),
	})
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("申请加入家族失败")
		return
	}

	res.ApplyList = familyInfo.DataToApplyUsersPb(logger)
	return
}

// OnComfirmApplyFamilyRQ 处理加入家族请求
func (f *Family) OnConfirmApplyFamilyRQ_10589_10590(s *session.Session, req *MazeFamily.ConfirmApplyFamilyRQ) (err error) {
	logger := log.Clone("OnComfirmApplyFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.ConfirmApplyFamilyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId
	res.ApplyUser = req.ApplyUser
	res.Result = req.Result

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnComfirmApplyFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	if req.GetResult() != MazeFamily.ApplyFamilyResult_APPLY_FAMILY_RESULT_AGREE &&
		req.GetResult() != MazeFamily.ApplyFamilyResult_APPLY_FAMILY_RESULT_REFUSE {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// 校验用户是否在家族的请求列表中
	ok, err := familyservice.GlobalFamilyService.CheckUserInApplyList(logger, req.GetFamilyId(), req.GetApplyUser().GetUserId())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("校验用户是否在家族的请求列表中失败")
		return
	}
	if !ok {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("用户不在家族的请求列表中")
		return
	}

	// todo 后续根据需求对 uid 做鉴权

	if req.GetResult() == MazeFamily.ApplyFamilyResult_APPLY_FAMILY_RESULT_AGREE {
		// 确认加入家族
		familyInfo, err := familyservice.GlobalFamilyService.AgreeApplyFamily(logger, req.GetFamilyId(), uid, familymodel.FamilyMember{
			UserID:         req.GetApplyUser().GetUserId(),
			NickName:       req.GetApplyUser().GetNickName(),
			Avatar:         req.GetApplyUser().GetAvatar(),
			Sex:            req.GetApplyUser().GetSex(),
			Level:          req.GetApplyUser().GetLevel(),
			PrivilegeLevel: int32(req.GetApplyUser().GetPrivilegeLevel()),
		})
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("确认加入家族失败")
			return err
		}

		err = familyservice.GlobalFamilyService.SetUserFamily(logger, req.GetApplyUser().GetUserId(), req.GetFamilyId())
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("设置玩家所在家族失败")
			return err
		}

		res.FamilyInfo = familyInfo.DataToFamilyInfoPb(logger)
		res.MemberList = familyInfo.DataToFamilyMembersPb(logger)
		res.ApplyList = familyInfo.DataToApplyUsersPb(logger)
	} else {
		// 拒绝加入家族
		familyInfo, err := familyservice.GlobalFamilyService.RefuseApplyFamily(logger, req.GetFamilyId(), uid, familymodel.FamilyMember{
			UserID:         req.GetApplyUser().GetUserId(),
			NickName:       req.GetApplyUser().GetNickName(),
			Avatar:         req.GetApplyUser().GetAvatar(),
			Sex:            req.GetApplyUser().GetSex(),
			Level:          req.GetApplyUser().GetLevel(),
			PrivilegeLevel: int32(req.GetApplyUser().GetPrivilegeLevel()),
		})
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("拒绝加入家族失败")
			return err
		}

		res.FamilyInfo = familyInfo.DataToFamilyInfoPb(logger)
		res.MemberList = familyInfo.DataToFamilyMembersPb(logger)
		res.ApplyList = familyInfo.DataToApplyUsersPb(logger)
	}
	return
}

// OnOperatePrivilegeRQ 修改某个用户的权限
func (f *Family) OnOperatePrivilegeRQ_10592_10593(s *session.Session, req *MazeFamily.OperatePrivilegeRQ) (err error) {
	logger := log.Clone("OnOperatePrivilegeRQ", uint64(s.UID()), 0)
	res := &MazeFamily.OperatePrivilegeRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId
	res.TargetPrivilegeLevel = req.TargetPrivilegeLevel

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnOperatePrivilegeRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 || len(req.GetOperateUsers()) <= 0 ||
		req.GetTargetPrivilegeLevel() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// todo 检查修改用户是否都在家族中
	for _, userID := range req.GetOperateUsers() {
		ok, err := familyservice.GlobalFamilyService.CheckUserInFamily(logger, req.GetFamilyId(), userID)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("检查用户是否在家族中失败")
			return err
		}
		if !ok {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("用户不在家族中")
			return err
		}
	}

	err = familyservice.GlobalFamilyService.OperatePrivilege(logger, uid, req.GetFamilyId(), int32(req.GetTargetPrivilegeLevel()), req.GetOperateUsers())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("操作权限失败")
		return
	}
	return
}

// OnKickFamilyRQ 踢出家族
func (f *Family) OnKickFamilyRQ_10594_10595(s *session.Session, req *MazeFamily.KickFamilyRQ) (err error) {
	logger := log.Clone("OnKickFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.KickFamilyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId
	res.KickUsers = req.KickUsers

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnKickFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 || len(req.GetKickUsers()) <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// 检查用户和踢出用户是否在同一个家族
	for _, userID := range req.GetKickUsers() {
		ok, err := familyservice.GlobalFamilyService.CheckUserInFamily(logger, req.GetFamilyId(), userID)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("检查用户是否在家族中失败")
			return err
		}

		if !ok {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("用户不在家族中")
			return err
		}
	}

	// 踢出家族
	err = familyservice.GlobalFamilyService.KickFamily(logger, uid, req.GetFamilyId(), req.GetKickUsers())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("踢出家族失败")
		return
	}

	// 删除玩家绑定家族
	for _, userID := range req.GetKickUsers() {
		err = familyservice.GlobalFamilyService.DeleteUserFamily(logger, userID)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.Wrap("删除玩家绑定家族失败")
			return
		}
	}

	// 发送踢出家族id包
	familyservice.GlobalFamilyService.SendKickFamilyID(logger, uid, req.GetFamilyId(), req.GetKickUsers())
	return
}

// OnExitFamilyRQ 退出家族
func (f *Family) OnExitFamilyRQ_10597_10598(s *session.Session, req *MazeFamily.ExitFamilyRQ) (err error) {
	logger := log.Clone("OnExitFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.ExitFamilyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.FamilyId = req.FamilyId

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnExitFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 退出家族
	err = familyservice.GlobalFamilyService.ExitFamily(logger, uid, req.GetFamilyId())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("退出家族失败")
		return
	}

	// 删除用户绑定家族
	err = familyservice.GlobalFamilyService.DeleteUserFamily(logger, uid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("删除用户绑定家族失败")
		return
	}

	// 记录玩家退出家族时间
	err = familyservice.GlobalFamilyService.SetUserLastLeaveFamilyTime(logger, uid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("记录玩家退出家族时间失败")
		return
	}

	// 发送退出家族id包
	familyservice.GlobalFamilyService.SendExitFamilyID(logger, uid, req.GetFamilyId())
	return
}

// 获取用户家族
func (f *Family) OnGetUserFamliyRQ_10600_10601(s *session.Session, req *MazeFamily.GetUserFamilyRQ) (err error) {
	logger := log.Clone("OnGetUserFamliyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetUserFamilyRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetUserFamliyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	// 检查参数
	if uid <= 0 || req.GetUserId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 获取用户家族
	familyID, err := familyservice.GlobalFamilyService.GetUserFamily(logger, uid)
	res.FamilyId = proto.Int32(familyID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取用户家族失败")
		return
	}
	return
}

func (f *Family) OnDissolutionFamilyRQ_10602_10603(s *session.Session, req *MazeFamily.DissolutionFamilyRQ) (err error) {
	logger := log.Clone("OnDissolutionFamilyRQ", uint64(s.UID()), 0)
	res := &MazeFamily.DissolutionFamilyRS{}

	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnDissolutionFamilyRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	// 检查参数
	if uid <= 0 || req.GetFamilyId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	// 解散家族
	err = familyservice.GlobalFamilyService.DissolutionFamily(logger, uid, req.GetFamilyId())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("解散家族失败")
		return
	}
	return
}
