package log

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// Clone 获取logger
func Clone(compName string, UserID uint64, logID int64) (logger fklog.FKLogI) {
	if logID <= 0 {
		logID = time.Now().UnixNano()
	}
	logger = fklog.AppLogger().Clone(compName)
	logger.SetUid(UserID)
	logger.SetLogId(logID)
	return
}
