package familymodel

import (
	"maze_game_server/io/redis/familyredis"
	"maze_game_server/lib/serialize"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

type FamilyListModel struct {
	Familys []int32 `json:"familys,omitempty"`
}

func LoadFamilyListModel(logger fklog.FKLogI) (r *FamilyListModel, err error) {
	r = &FamilyListModel{}
	if err = r.load(logger); err != nil {
		logger.ErrorWF("LoadFamilyListModel err", zap.Error(err))
		return
	}
	return
}

func (r *FamilyListModel) load(logger fklog.FKLogI) (err error) {
	value, err := familyredis.GetFamilyList()
	if err != nil {
		logger.ErrorWF("LoadFamilyInfoModel err",
			zap.Error(err))
		return err
	}

	if value == nil {
		r.Familys = make([]int32, 0)
		return nil
	}

	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.ErrorWF("LoadFamilyInfoModel err",
			zap.Error(err))
		return err
	}
	return
}

func (r *FamilyListModel) Save(logger fklog.FKLogI) (err error) {
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.ErrorWF("FamilyListModel Save Marshal err",
			zap.Error(err))
		return err
	}
	err = familyredis.SetFamilyIDs(value)
	if err != nil {
		logger.ErrorWF("Save SetFamilyIDs err")
		return
	}
	return
}

func (r *FamilyListModel) Delete(logger fklog.FKLogI) (err error) {
	err = familyredis.DelFamilyList()
	if err != nil {
		logger.ErrorWF("Delete DelFamilyList err")
		return
	}
	return
}

func (r *FamilyListModel) AddFamily(logger fklog.FKLogI, familyId int32) {
	r.Familys = append(r.Familys, familyId)
}

func (r *FamilyListModel) RemoveFamily(logger fklog.FKLogI, familyId int32) {
	for i, v := range r.Familys {
		if v == familyId {
			r.Familys = append(r.Familys[:i], r.Familys[i+1:]...)
			break
		}
	}
}
