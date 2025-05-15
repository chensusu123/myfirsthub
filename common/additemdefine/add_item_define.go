/*
@Author: xiaobo
@Date: 2023/11/30 14:20
@Description: 加物品相关结构、接口
*/

package additemdefine

import (
	"fmt"
	"reflect"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeItemsV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structdefine"
	"go.uber.org/zap"
)

type AddItemOption struct {
	OpType      int32                                         // / 添加类型 ENUM_OP_TYPE,必填
	TradeNumber uint64                                        // / 交易流水号
	RegIns      *RegisterInfo                                 `json:"-"` // 物品管理注册, 有的业务需要根据情况 再次添加物品
	CheckResMap map[ItemClassProcess]*structdefine.AddItemRes `json:"-"` // 提前校验所有物品上限，如果有该值则不再校验
}

func (a *AddItemOption) GetCheckRes(class ItemClassProcess) (*structdefine.AddItemRes, bool) {
	if a.CheckResMap == nil {
		return nil, false
	}
	res, ok := a.CheckResMap[class]
	return res, ok
}

type ItemClassProcess interface {
	// GatherItem 加道具
	GatherItem(userCtx fkserver.UserContext, option *AddItemOption, items ...*MazeCommon.MazeItem) (*structdefine.AddItemRes, error)
	// GatherItemCheck 加道具检查
	GatherItemCheck(userCtx fkserver.UserContext, option *AddItemOption, items []*MazeCommon.MazeItem) (*structdefine.AddItemRes, error)
	// DeductItem 扣道具
	DeductItem(userCtx fkserver.UserContext, option *AddItemOption, items []*MazeCommon.MazeItem) (*structdefine.AddItemRes, *MessageType.ErrorInfo)
	// DeductItemCheck 扣道具检查
	DeductItemCheck(userCtx fkserver.UserContext, option *AddItemOption, items []*MazeCommon.MazeItem) (*structdefine.AddItemRes, *MessageType.ErrorInfo)
	// GetItem 查询数据
	GetItem(userCtx fkserver.UserContext, option *AddItemOption, items []*MazeCommon.MazeItem) *MessageType.ErrorInfo
}

type RegisterInfo struct {
	// itemID--> handleClass
	ItemMap map[int32]ItemClassProcess
	// itemType--> handleClass
	TypeMap map[int32]ItemClassProcess
	// bagType--> handleClass
	BagMap map[int32]ItemClassProcess
	// 通用转发类型
	CommonTransMap map[int32]ItemClassProcess
}

func NewRegister() *RegisterInfo {
	return &RegisterInfo{
		ItemMap:        make(map[int32]ItemClassProcess),
		TypeMap:        make(map[int32]ItemClassProcess),
		BagMap:         make(map[int32]ItemClassProcess),
		CommonTransMap: make(map[int32]ItemClassProcess),
	}
}

// RegisterByItemID itemID 注册
func (reg *RegisterInfo) RegisterByItemID(itemID int32, handleClass ItemClassProcess) {
	r, ok := reg.ItemMap[itemID]
	if ok {
		exist := reflect.TypeOf(r).Elem()
		input := reflect.TypeOf(handleClass).Elem()
		fmt.Println("RegisterByItemID already register item id", itemID, "exist=", exist.Name(), "input=", input.Name())
		return

	}
	reg.ItemMap[itemID] = handleClass
}

// RegisterByItemType itemType 注册
func (reg *RegisterInfo) RegisterByItemType(itemType int32, handleClass ItemClassProcess) {
	r, ok := reg.TypeMap[itemType]
	if ok {
		exist := reflect.TypeOf(r).Elem()
		input := reflect.TypeOf(handleClass).Elem()
		fmt.Println("RegisterByItemType already register item type", itemType, "exist=", exist.Name(), "input=", input.Name())
		return
	}
	reg.TypeMap[itemType] = handleClass
}

// RegisterByBagType isBag注册  1:进仓库，2:进船背包 3:进迷宫背包
func (reg *RegisterInfo) RegisterByBagType(bagType int32, handleClass ItemClassProcess) {
	r, ok := reg.BagMap[bagType]
	if ok {
		exist := reflect.TypeOf(r).Elem()
		input := reflect.TypeOf(handleClass).Elem()
		fmt.Println("RegisterByBagType already register bag type", bagType, "exist=", exist.Name(), "input=", input.Name())
		return
	}
	reg.BagMap[bagType] = handleClass
}

func (reg *RegisterInfo) RegisterCommonTrans(idType int32, handleClass ItemClassProcess) {
	r, ok := reg.CommonTransMap[idType]
	if ok {
		exist := reflect.TypeOf(r).Elem()
		input := reflect.TypeOf(handleClass).Elem()
		fmt.Println("RegisterByPayType already register common trans", "exist=", exist.Name(), "input=", input.Name())
		return
	}
	reg.CommonTransMap[idType] = handleClass
}

func (reg *RegisterInfo) GetClassByID(logger fklog.FKLogI, itemID int32) (v ItemClassProcess, err error) {
	if itemID == 0 {
		logger.WarnWF("GetClassByID item id 0")
		err = fmt.Errorf("item id 0")
		return
	}
	itemCfg := GMazeItemsV8Cfg.Get(itemID)
	if itemCfg == nil {
		logger.ErrorWF("GetClassByID not find item", zap.Int32("itemID", itemID))
		return nil, fmt.Errorf("GMazeItemsV8Cfg nil %d", itemID)
	}

	ok := false

	// 根据进入背包类型
	v, ok = reg.BagMap[itemCfg.Is_bag]
	if ok {
		return v, nil
	}

	// 根据itemID
	v, ok = reg.ItemMap[itemID]
	if ok {
		return v, nil
	}

	// 根据itemType
	v, ok = reg.TypeMap[itemCfg.Type]
	if ok {
		return v, nil
	}

	err = fmt.Errorf("not find item register handleProcess.itemID=%d", itemID)
	return
}

// GetClassifiedItems 根据不同的物品对物品进行分类
func (reg *RegisterInfo) GetClassifiedItems(logger fklog.FKLogI, items ...*MazeCommon.MazeItem) (classItemMap map[ItemClassProcess][]*MazeCommon.MazeItem, err error) {
	if reg == nil {
		err = fmt.Errorf("RegisterInfo nil")
		return
	}

	classItemMap = make(map[ItemClassProcess][]*MazeCommon.MazeItem)

	for idx, item := range items {
		var handleProcess ItemClassProcess
		handleProcess, err = reg.GetClassByID(logger, item.GetItemId())
		if err != nil {
			return classItemMap, err
		}
		_, ok := classItemMap[handleProcess]
		if !ok {
			classItemMap[handleProcess] = make([]*MazeCommon.MazeItem, 0)
		}
		classItemMap[handleProcess] = append(classItemMap[handleProcess], items[idx])
	}
	return
}
