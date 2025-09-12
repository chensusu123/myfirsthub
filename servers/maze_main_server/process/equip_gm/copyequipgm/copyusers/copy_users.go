package copyusers

import (
	"context"
	"sync"
	"time"

	"maze_game_server/common/function/fileio"
	"maze_game_server/servers/maze_main_server/process/equip_gm/asynctask"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type CopyUsers struct {
	UserIdMap map[uint64][]uint64
}

func NewCopyUsers() *CopyUsers {
	info := new(CopyUsers)
	info.UserIdMap = make(map[uint64][]uint64)
	return info
}

func (m *CopyUsers) GetCopyUsers(userId uint64) []uint64 {
	return m.UserIdMap[userId]
}

// 加载用户
func (m *CopyUsers) LoadUser(ctx context.Context, path string, sp string) error {
	logger := fklog.ContextAppLogger(ctx)
	fr := fileio.NewDefFReaderEx(logger, sp)
	err := fr.Open(path)
	if err != nil {
		fkfmt.Println("err", err)
		return err
	}
	fr.SetDumpRow(1000)
	fkfmt.Println("open file", path, "succ")

	f := func(logger fklog.FKLogI, line []uint64) bool {
		if len(line) < 2 {
			return true
		}
		m.UserIdMap[line[0]] = append(m.UserIdMap[line[0]], line[1])
		return true
	}
	fr.Range(f)
	fr.Close()
	return nil
}

type CallBackF func(ctx context.Context, srcUserId uint64, dstUserId []uint64) bool

func (m *CopyUsers) RangeUser(ctx context.Context, f CallBackF) {
	logger := fklog.ContextAppLogger(ctx)
	var wg sync.WaitGroup
	for src, dstList := range m.UserIdMap {
		userLogger := logger.Clone("copy_task")
		userLogger.SetUid(src)
		userLogger.SetLogId(time.Now().UnixNano())

		tmpSrc := src
		var tmpDstList []uint64
		tmpDstList = append(tmpDstList, dstList...)
		wg.Add(1)
		asynctask.GWorkGroupBusiness.SendTask(tmpSrc, func() {
			wg.Done()
			_ = f(ctx, tmpSrc, tmpDstList)
		})
	}
	wg.Wait()
}
