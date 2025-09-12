package familymodel

import (
	"context"

	"maze_game_server/io/redis/familyredis"
	"maze_game_server/lib/serialize"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type FamilyListModel struct {
	Familys []int32 `json:"familys,omitempty"`
}

func LoadFamilyListModel(ctx context.Context) (r *FamilyListModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	r = &FamilyListModel{}
	if err = r.load(ctx); err != nil {
		logger.CtxError(ctx, "LoadFamilyListModel err", zap.Error(err))
		return
	}
	return
}

func (r *FamilyListModel) load(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := familyredis.GetFamilyList(ctx)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyInfoModel err",
			zap.Error(err))
		return err
	}

	if value == nil {
		r.Familys = make([]int32, 0)
		return nil
	}

	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyInfoModel err",
			zap.Error(err))
		return err
	}
	return
}

func (r *FamilyListModel) Save(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.CtxError(ctx, "FamilyListModel Save Marshal err",
			zap.Error(err))
		return err
	}
	err = familyredis.SetFamilyIDs(ctx, value)
	if err != nil {
		logger.CtxError(ctx, "Save SetFamilyIDs err")
		return
	}
	return
}

func (r *FamilyListModel) Delete(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	err = familyredis.DelFamilyList(ctx)
	if err != nil {
		logger.CtxError(ctx, "Delete DelFamilyList err")
		return
	}
	return
}

func (r *FamilyListModel) AddFamily(ctx context.Context, familyId int32) {
	r.Familys = append(r.Familys, familyId)
}

func (r *FamilyListModel) RemoveFamily(ctx context.Context, familyId int32) {
	for i, v := range r.Familys {
		if v == familyId {
			r.Familys = append(r.Familys[:i], r.Familys[i+1:]...)
			break
		}
	}
}
