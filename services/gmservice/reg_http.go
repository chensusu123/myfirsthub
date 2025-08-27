package gmservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func (s *service) RegHttp(logger fklog.FKLogI) {
	// user
	s.SafeGETRegister(logger, "/LookAssembleInfo", s.LookAssembleInfo)
	s.SafeGETRegister(logger, "/SetBarrier", s.SetBarrier)
	s.SafeGETRegister(logger, "/DumpBattleData", s.DumpBattleData)
	s.SafeGETRegister(logger, "/attrs", s.Attrs)
	// buff
	s.SafeGETRegister(logger, "/setMazeTempBuff", s.SetMazeTempBuff)
	// equip
	s.SafeGETRegister(logger, "/AddEquip", s.AddEquip)
	s.SafeGETRegister(logger, "/GmEquipPosLvUp", s.GmEquipPosLvUp)
	s.SafeGETRegister(logger, "/GetEquipInfoByCfgId", s.GetEquipInfoByCfgId)
	s.SafeGETRegister(logger, "/GetEquipInfoByGuid", s.GetEquipInfoByGuid)
	s.SafeGETRegister(logger, "/SendOneSuitEquip", s.SendOneSuitEquip)
	s.SafeGETRegister(logger, "/ReInitDollEquip", s.ReInitDollEquip)
	s.SafeGETRegister(logger, "/FixDollAttr", s.FixDollAttr)
	s.SafeGETRegister(logger, "/FixEquipPosUnlock", s.FixEquipPosUnlock)
	s.SafeGETRegister(logger, "/FixAssembleEquipInfo", s.FixAssembleEquipInfo)
	s.SafeGETRegister(logger, "/SetEquipRollScore", s.SetEquipRollScore)
	s.SafeGETRegister(logger, "/BatchAddEquip", s.BatchAddEquip)
	// ohter
	s.SafeGETRegister(logger, "/ClearBag", s.ClearBag)
	s.SafeGETRegister(logger, "/ClearBagNotAssemble", s.ClearBagNotAssemble)
	s.SafeGETRegister(logger, "/generateUser", s.GenerateUser)
	s.SafeGETRegister(logger, "/online", s.Online)
	// item
	s.SafeGETRegister(logger, "/AddExp", s.AddExp)
	s.SafeGETRegister(logger, "/addRefreshCost", s.AddRefreshCost)
	s.SafeGETRegister(logger, "/addEnergy", s.AddEnergy)
	s.SafeGETRegister(logger, "/addItem", s.AddItem)
	// excel
	s.SafeGETRegister(logger, "/showSheet", s.ShowSheet)
	s.SafeGETRegister(logger, "/GetExcelList", s.GetExcelList)
	s.SafeGETRegister(logger, "/GetExcelSheet", s.GetExcelSheet)
	s.SafeGETRegister(logger, "/GetExcelData", s.GetExcelData)
}
