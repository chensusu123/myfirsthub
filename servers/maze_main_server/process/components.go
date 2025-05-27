package process

import (
	"maze_game_server/lib/codec"
	"maze_game_server/lib/nano/component"
	"maze_game_server/servers/maze_main_server/process/attr_calc"
	"maze_game_server/servers/maze_main_server/process/auth"
	"maze_game_server/servers/maze_main_server/process/buff"
	"maze_game_server/servers/maze_main_server/process/card"
	"maze_game_server/servers/maze_main_server/process/collect"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/equip_gm"
	"maze_game_server/servers/maze_main_server/process/game"
	"maze_game_server/servers/maze_main_server/process/game/energy"
	"maze_game_server/servers/maze_main_server/process/game/sweep"
	"maze_game_server/servers/maze_main_server/process/interact"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/servers/maze_main_server/process/rob"
)

func Components() (comps *component.Components, routes *codec.Routes) {
	comps = &component.Components{}
	routes = &codec.Routes{}
	// 注册nano组件与自定义包解析路由
	reg := func(comp component.Component) {
		comps.Register(comp)
		routes.Register(comp)
	}

	{
		reg(auth.NewAuth())          // 验证组件
		reg(game.NewGame())          // 关卡组件
		reg(attr_calc.NewProperty()) // 人物属性
		reg(energy.NewEnergy())      // 体力组件
		reg(sweep.NewSweep())        // 扫荡组件
		reg(buff.NewBuff())          // Buff组件
		reg(collect.NewCollect())    // 挂机组件
		reg(equip.NewEquip())        // 装备组件
		reg(equip_gm.NewEquipGM())   // 装备批量操作组件
		reg(item.NewItem())          // 道具组件
		reg(card.NewCard())          // 月卡
		reg(rob.NewRob())            // 掠夺
		reg(interact.NewInteract())  // 交互
	}

	return
}
