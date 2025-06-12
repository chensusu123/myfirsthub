// @Author pangchenyang 2025/6/9 21:49:00
// @Desc: 
package userprofilemysql

import (
	"gorm.io/gorm"
	"time"
	"maze_game_server/io/mysql"
	"sync"
)

var db *gorm.DB
var once sync.Once

// func init() {
// 	var err error
// 	db, err = mysql.GetMysqlDb()
// 	if err != nil {
// 		db.Logger.Error(nil, "userprofilemysql GetMysqlDb fail, err:%v", err)
// 	}
// }

func InitMysql() {
	once.Do(func() {
		var err error
		db, err = mysql.GetMysqlDb()
		if err != nil {
			db.Logger.Error(nil, "userprofilemysql GetMysqlDb fail, err:%v", err)
		}
	})
}

// UserProfileMysql 用户资料表模型
type UserProfile struct {
	UserID    uint64    `gorm:"column:user_id;primaryKey"` // 用户ID
	Nickname  string    `gorm:"column:nickname;size:32"`   // 昵称
	IconUrl   string    `gorm:"column:iconUrl"`            // 头像ID
	Sex       uint8     `gorm:"column:sex"`                // 性别(0:未知,1:男,2:女)
	CreatedAt time.Time `gorm:"column:created_at"`         // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at"`         // 更新时间
}

func GetTableName() string {
	return "user_profiles"
}

// UserProfileMysql 用户资料数据库
type OpMysql struct {
	db *gorm.DB
}

// NewUserProfileMysql 构造函数
func NewUserProfileMysql(db *gorm.DB) *OpMysql {
	return &OpMysql{db: db}
}

// Create 创建用户资料
func Create(profile *UserProfile) error {
	return db.Create(profile).Error
}

// GetByUserID 根据用户ID获取资料
func GetByUserID(userID uint64) (*UserProfile, error) {
	var profile UserProfile
	err := db.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

// Update 更新用户资料
func Update(profile *UserProfile) error {
	return db.Save(profile).Error
}

// Delete 删除用户资料
func Delete(userID uint64) error {
	return db.Where("user_id = ?", userID).Delete(&UserProfile{}).Error
}

// GetByNickname 根据昵称获取用户资料
func GetByNickname(nickname string) (*UserProfile, error) {
	var profile UserProfile
	err := db.Where("nickname = ?", nickname).First(&profile).Error
	return &profile, err
}

// GetByCreatedAt 根据创建时间范围查询用户资料
func GetByCreatedAt(start, end time.Time) ([]UserProfile, error) {
	var profiles []UserProfile
	err := db.Where("created_at BETWEEN ? AND ?", start, end).Find(&profiles).Error
	return profiles, err
}

// BatchGetByUserIDs 根据用户ID列表批量查询用户资料
func BatchGetByUserIDs(userIDs []uint64) ([]*UserProfile, error) {
	var profiles []*UserProfile
	err := db.Where("user_id IN (?)", userIDs).Find(&profiles).Error
	return profiles, err
}

// List 分页查询用户资料
func List(page, pageSize int) ([]*UserProfile, error) {
	var profiles []*UserProfile
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&profiles).Error
	return profiles, err
}

// Count 获取用户资料总数
func Count() (int64, error) {
	var count int64
	err := db.Model(&UserProfile{}).Count(&count).Error
	return count, err
}
