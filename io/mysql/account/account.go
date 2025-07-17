package account

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/mysql"
)

type AccountTableInfo struct {
	AccountID    uint64 `gorm:"column:account_id"`
	Token        uint64 `gorm:"column:token"`
	WeixinID     string `gorm:"column:weixin_id"`
	DouyinID     string `gorm:"column:douyin_id"`
	Email        string `gorm:"column:email"`
	Password     string `gorm:"column:password"`
	ExistUseInfo uint8  `gorm:"column:exist_user_info"`
	CreateDt     int64  `gorm:"column:create_dt"`
	ChangeDt     int64  `gorm:"column:change_dt"`
}

// 从maze_register_serve复制过来的方法
// 通过account_id找到存储 t_account_info_x表位置
func GetAppTableIndex(id uint64) uint32 {
	return uint32((id >> 4) % 16)
}

func GetAccountInfo(logger fklog.FKLogI, accountId uint64) (*AccountTableInfo, error) {
	db, err := mysql.GetMysqlDb()
	if err != nil {
		logger.ErrorWF("GetMysqlDb fail", zap.Error(err))
		return nil, err
	}
	tableIndex := GetAppTableIndex(accountId)
	table := fmt.Sprintf("t_account_info_%d", tableIndex)
	accountTableInfo := &AccountTableInfo{}
	res := db.Table(table).Where("account_id = ?", accountId).First(accountTableInfo)
	if res.Error != nil {
		logger.ErrorWF("GetAccountInfo fail", zap.Error(err), zap.Uint64("accountId", accountId))
		return nil, res.Error
	}

	return accountTableInfo, nil
}
