package app

import (
	"fmt"
)

type User interface {
	// ID 用户ID
	UserID() uint64

	// Endpoint 终端标识
	Endpoint() string

	// String 唯一标识
	String() string
}

type user struct {
	userID   uint64
	endpoint string
}

// WrapUser 整合用户基本信息，返回User实例接口
func WrapUser(userID uint64, endpoint string) (u User, err error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if len(endpoint) <= 0 {
		endpoint = "default"
	}
	return &user{userID: userID, endpoint: endpoint}, nil
}

// UserID implements User.
func (u *user) UserID() uint64 {
	return u.userID
}

// Endpoint implements User.
func (u *user) Endpoint() string {
	return u.endpoint
}

// String implements User.
func (u *user) String() string {
	return fmt.Sprintf("%d:%s", u.userID, u.endpoint)
}
