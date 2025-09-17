package gmservice

import "context"

func (s *service) RegHttp(ctx context.Context) {
	// user
	s.SafeGETRegister(ctx, "/LookAssembleInfo", s.LookAssembleInfo)
	s.SafeGETRegister(ctx, "/SetBarrier", s.SetBarrier)
	s.SafeGETRegister(ctx, "/DumpBattleData", s.DumpBattleData)
	s.SafeGETRegister(ctx, "/attrs", s.Attrs)
	s.SafePOSTRegister(ctx, "/SetLevel", s.SetUserLevel)
	// buff
	s.SafePOSTRegister(ctx, "/setMazeTempBuff", s.SetMazeTempBuff)
	// equip
	s.SafePOSTRegister(ctx, "/AddEquip", s.AddEquip)
	s.SafePOSTRegister(ctx, "/GmEquipPosLvUp", s.GmEquipPosLvUp)
	s.SafeGETRegister(ctx, "/GetEquipInfoByCfgId", s.GetEquipInfoByCfgId)
	s.SafeGETRegister(ctx, "/GetEquipInfoByGuid", s.GetEquipInfoByGuid)
	s.SafeGETRegister(ctx, "/SendOneSuitEquip", s.SendOneSuitEquip)
	s.SafeGETRegister(ctx, "/ReInitDollEquip", s.ReInitDollEquip)
	s.SafeGETRegister(ctx, "/FixDollAttr", s.FixDollAttr)
	s.SafeGETRegister(ctx, "/FixEquipPosUnlock", s.FixEquipPosUnlock)
	s.SafeGETRegister(ctx, "/FixAssembleEquipInfo", s.FixAssembleEquipInfo)
	s.SafeGETRegister(ctx, "/SetEquipRollScore", s.SetEquipRollScore)
	s.SafeGETRegister(ctx, "/BatchAddEquip", s.BatchAddEquip)
	// ohter
	s.SafeGETRegister(ctx, "/ClearBag", s.ClearBag)
	s.SafeGETRegister(ctx, "/ClearBagNotAssemble", s.ClearBagNotAssemble)
	s.SafeGETRegister(ctx, "/generateUser", s.GenerateUser)
	s.SafeGETRegister(ctx, "/online", s.Online)
	// item
	s.SafePOSTRegister(ctx, "/AddExp", s.AddExp)
	s.SafePOSTRegister(ctx, "/addRefreshCost", s.AddRefreshCost)
	s.SafePOSTRegister(ctx, "/addEnergy", s.AddEnergy)
	s.SafePOSTRegister(ctx, "/addItem", s.AddItem)
	// excel
	s.SafeGETRegister(ctx, "/showSheet", s.ShowSheet)
	s.SafeGETRegister(ctx, "/GetExcelList", s.GetExcelList)
	s.SafeGETRegister(ctx, "/GetExcelSheet", s.GetExcelSheet)
	s.SafeGETRegister(ctx, "/GetExcelData", s.GetExcelData)

	// alliance
	s.SafeGETRegister(ctx, "/CreateAlliance", s.CreateAlliance)
	s.SafeGETRegister(ctx, "/GetAllianceInfo", s.GetAllianceInfo)
	s.SafeGETRegister(ctx, "/ChangeAllianceInfo", s.ChangeAllianceInfo)
	s.SafeGETRegister(ctx, "/BatchChangeAlliance", s.BatchChangeAlliance)
	s.SafeGETRegister(ctx, "/GetAllianceList", s.GetAllianceList)
	s.SafeGETRegister(ctx, "/GetUserAlliance", s.GetUserAlliance)

	// tmp
	s.SafeGETRegister(ctx, "/SetRecommendSize", s.SetRecommendSize)
}
