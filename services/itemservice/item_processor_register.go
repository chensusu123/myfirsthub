package itemservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"reflect"
)

type RegisterInfo struct {
	// itemID--> handleClass
	ItemMap map[int32]ItemProcessor
	// itemType--> handleClass
	TypeMap map[int32]ItemProcessor
	// bagType--> handleClass
	BagMap map[int32]ItemProcessor
}

func NewRegister() *RegisterInfo {
	return &RegisterInfo{
		ItemMap: make(map[int32]ItemProcessor),
		TypeMap: make(map[int32]ItemProcessor),
		BagMap:  make(map[int32]ItemProcessor),
	}
}

// RegisterByItemID itemID 注册
func (reg *RegisterInfo) RegisterByItemID(itemID int32, handleClass ItemProcessor) {
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
func (reg *RegisterInfo) RegisterByItemType(itemType int32, handleClass ItemProcessor) {
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
func (reg *RegisterInfo) RegisterByBagType(bagType int32, handleClass ItemProcessor) {
	r, ok := reg.BagMap[bagType]
	if ok {
		exist := reflect.TypeOf(r).Elem()
		input := reflect.TypeOf(handleClass).Elem()
		fmt.Println("RegisterByBagType already register bag type", bagType, "exist=", exist.Name(), "input=", input.Name())
		return
	}
	reg.BagMap[bagType] = handleClass
}

func (reg *RegisterInfo) getProcessorByID(ctx context.Context, itemID int32) (v ItemProcessor, err error) {
	logger := fklog.ContextAppLogger(ctx)
	if itemID == 0 {
		logger.CtxError(ctx, "getProcessorByID item id 0")
		err = fmt.Errorf("item id 0")
		return
	}
	itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx, itemID)
	if itemCfg == nil {
		logger.CtxError(ctx, "getProcessorByID not find item", zap.Int32("itemID", itemID))
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

// GroupItemsByProcessor 根据processor对物品进行分类
func (reg *RegisterInfo) GroupItemsByProcessor(ctx context.Context, items ...*ItemInfo) (classItemMap map[ItemProcessor][]*ItemInfo, err error) {
	if reg == nil {
		err = fmt.Errorf("RegisterInfo nil")
		return
	}

	classItemMap = make(map[ItemProcessor][]*ItemInfo)

	for idx, item := range items {
		var processor ItemProcessor
		processor, err = reg.getProcessorByID(ctx, item.ItemId)
		if err != nil {
			return classItemMap, err
		}
		_, ok := classItemMap[processor]
		if !ok {
			classItemMap[processor] = make([]*ItemInfo, 0)
		}
		classItemMap[processor] = append(classItemMap[processor], items[idx])
	}
	return
}
