package gmservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func (s *service) RegHttp(logger fklog.FKLogI) {
	// user
	s.SafeHttpRegister(logger, "/LookAssembleInfo", s.LookAssembleInfo)
	s.SafeHttpRegister(logger, "/SetBarrier", s.SetBarrier)
	s.SafeHttpRegister(logger, "/DumpBattleData", s.DumpBattleData)
	s.SafeHttpRegister(logger, "/attrs", s.Attrs)

	// buff
	s.SafeHttpRegister(logger, "/setMazeTempBuff", s.SetMazeTempBuff)

	// equip
	s.SafeHttpRegister(logger, "/AddEquip", s.AddEquip)
	s.SafeHttpRegister(logger, "/GmEquipPosLvUp", s.GmEquipPosLvUp)
	s.SafeHttpRegister(logger, "/GetEquipInfoByCfgId", s.GetEquipInfoByCfgId)
	s.SafeHttpRegister(logger, "/GetEquipInfoByGuid", s.GetEquipInfoByGuid)
	s.SafeHttpRegister(logger, "/FixDollEquipAttr", s.FixDollEquipAttr)
	s.SafeHttpRegister(logger, "/SendOneSuitEquip", s.SendOneSuitEquip)
	s.SafeHttpRegister(logger, "/ReInitDollEquip", s.ReInitDollEquip)
	s.SafeHttpRegister(logger, "/FixDollAttr", s.FixDollAttr)
	s.SafeHttpRegister(logger, "/FixEquipPosUnlock", s.FixEquipPosUnlock)
	s.SafeHttpRegister(logger, "/FixAssembleEquipInfo", s.FixAssembleEquipInfo)
	s.SafeHttpRegister(logger, "/SetEquipRollScore", s.SetEquipRollScore)
	s.SafeHttpRegister(logger, "/BatchAddEquip", s.BatchAddEquip)

	// ohter
	s.SafeHttpRegister(logger, "/ClearBag", s.ClearBag)
	s.SafeHttpRegister(logger, "/ClearBagNotAssemble", s.ClearBagNotAssemble)
	s.SafeHttpRegister(logger, "/generateUser", s.GenerateUser)
	s.SafeHttpRegister(logger, "/online", s.Online)

	// item
	s.SafeHttpRegister(logger, "/AddExp", s.AddExp)
	s.SafeHttpRegister(logger, "/addRefreshCost", s.AddRefreshCost)
	s.SafeHttpRegister(logger, "/addEnergy", s.AddEnergy)
	s.SafeHttpRegister(logger, "/addItem", s.AddItem)

	// excel
	s.SafeHttpRegister(logger, "/showSheet", s.ShowSheet)
	s.SafeHttpRegister(logger, "/GetExcelList", s.GetExcelList)
	s.SafeHttpRegister(logger, "/GetExcelSheet", s.GetExcelSheet)
	s.SafeHttpRegister(logger, "/GetExcelData", s.GetExcelData)
}
