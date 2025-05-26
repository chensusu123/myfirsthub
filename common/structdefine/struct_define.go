/*
@Author: xiaobo
@Date: 2023/11/30 14:13
@Description: 结构定义
*/

package structdefine

import "maze_game_server/pb/common/MazeCommon"

type AddItemRes struct {
	IsCheckErr  bool                   // 是否是检查过程出错,废弃
	SucItem     []*MazeCommon.MazeItem // 操作成功的道具
	FailItem    []*MazeCommon.MazeItem // 操作失败的道具(业务逻辑错误)
	TimeoutItem []*MazeCommon.MazeItem // 超时添加失败的道具(处理超时或者redis/rpc失败、超时)
	LessItem    []*MazeCommon.MazeItem // 扣道具时，检查阶段，不够扣的道具数量
	LimitItem   []*MazeCommon.MazeItem // 达到上限的道具
}

func (a *AddItemRes) Merge(b *AddItemRes) {
	if a == nil {
		return
	}
	if b == nil {
		return
	}
	a.SucItem = append(a.SucItem, b.SucItem...)
	a.FailItem = append(a.FailItem, b.FailItem...)
	a.TimeoutItem = append(a.TimeoutItem, b.TimeoutItem...)
	a.LessItem = append(a.LessItem, b.LessItem...)
}
