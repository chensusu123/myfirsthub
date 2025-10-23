package npclibrary

import (
	"maze_game_server/pb/common/NPCLibrary"
)

// 注册图鉴处理器
func RegisterNPCLibraryHandlers() {
	// TODO: 根据项目的实际注册方式来实现
	// 这里需要根据项目的具体架构来注册处理器
}

// 图鉴协议常量
const (
	// 客户端到服务器协议
	UN_TCP_PACK_CLI_ADD_NPC_LIBRARY_RQ     = 50501 // 添加图鉴请求
	UN_TCP_PACK_CLI_ADD_NPC_LIBRARY_RS     = 50502 // 添加图鉴响应
	UN_TCP_PACK_CLI_QUERY_NPC_LIBRARY_RQ   = 50503 // 查询图鉴请求
	UN_TCP_PACK_CLI_QUERY_NPC_LIBRARY_RS   = 50504 // 查询图鉴响应
	UN_TCP_PACK_CLI_UPGRADE_NPC_LIBRARY_RQ = 50505 // 升级图鉴请求
	UN_TCP_PACK_CLI_UPGRADE_NPC_LIBRARY_RS = 50506 // 升级图鉴响应
	UN_TCP_PACK_CLI_CHECK_NPC_LIBRARY_RQ   = 50507 // 检查图鉴请求
	UN_TCP_PACK_CLI_CHECK_NPC_LIBRARY_RS   = 50508 // 检查图鉴响应
	UN_TCP_PACK_CLI_GET_NEXT_LEVEL_INFO_RQ = 50509 // 获取下一级图鉴信息请求
	UN_TCP_PACK_CLI_GET_NEXT_LEVEL_INFO_RS = 50510 // 获取下一级图鉴信息响应
	UN_TCP_PACK_CLI_FORCE_COMPOSE_NPC_RQ   = 50511 // 强制合成NPC请求
	UN_TCP_PACK_CLI_FORCE_COMPOSE_NPC_RS   = 50512 // 强制合成NPC响应

	// 服务器到客户端协议
	UN_TCP_PACK_SVR_NPC_LIBRARY_NOTIFY = 50513 // 图鉴变化通知
)

// 获取协议消息类型映射
func GetNPCLibraryMessageTypes() map[int32]interface{} {
	return map[int32]interface{}{
		UN_TCP_PACK_CLI_ADD_NPC_LIBRARY_RQ:     &NPCLibrary.AddNPCLibraryRQ{},
		UN_TCP_PACK_CLI_ADD_NPC_LIBRARY_RS:     &NPCLibrary.AddNPCLibraryRS{},
		UN_TCP_PACK_CLI_QUERY_NPC_LIBRARY_RQ:   &NPCLibrary.QueryNPCLibraryRQ{},
		UN_TCP_PACK_CLI_QUERY_NPC_LIBRARY_RS:   &NPCLibrary.QueryNPCLibraryRS{},
		UN_TCP_PACK_CLI_UPGRADE_NPC_LIBRARY_RQ: &NPCLibrary.UpgradeNPCLibraryRQ{},
		UN_TCP_PACK_CLI_UPGRADE_NPC_LIBRARY_RS: &NPCLibrary.UpgradeNPCLibraryRS{},
		UN_TCP_PACK_CLI_CHECK_NPC_LIBRARY_RQ:   &NPCLibrary.CheckNPCLibraryRQ{},
		UN_TCP_PACK_CLI_CHECK_NPC_LIBRARY_RS:   &NPCLibrary.CheckNPCLibraryRS{},
		UN_TCP_PACK_CLI_GET_NEXT_LEVEL_INFO_RQ: &NPCLibrary.GetNextLevelInfoRQ{},
		UN_TCP_PACK_CLI_GET_NEXT_LEVEL_INFO_RS: &NPCLibrary.GetNextLevelInfoRS{},
		UN_TCP_PACK_CLI_FORCE_COMPOSE_NPC_RQ:   &NPCLibrary.ForceComposeNPCRQ{},
		UN_TCP_PACK_CLI_FORCE_COMPOSE_NPC_RS:   &NPCLibrary.ForceComposeNPCRS{},
		UN_TCP_PACK_SVR_NPC_LIBRARY_NOTIFY:     &NPCLibrary.NPCLibraryNotify{},
	}
}
