package app

import (
	"errors"

	"maze_game_server/io/redis/im/group"
	"maze_game_server/io/redis/im/msgstore"
	"maze_game_server/io/redis/im/session"
)

type Group = group.Group
type Session = session.Session
type Message = msgstore.Message

var (
	ErrInvalidAppID  = errors.New("invalid app identity")
	ErrInvalidUserID = errors.New("invalid user identity")
)

// IsValidApp 校验应用ID是否有效
func IsValidApp(appID int32) (err error) {
	if appID <= 0 {
		return ErrInvalidAppID
	}
	return nil
}

// IsValidUser 校验用户ID是否有效
func IsValidUser(userID uint64) (err error) {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	return nil
}

var (
	apps = map[int32]App{
		1: &app{appID: 1, name: "爱地牢"},
	}
)

type App interface {
	// ID 获取应用ID
	ID() int32

	// Name 获取应用名
	Name() string
}

type app struct {
	appID int32
	name  string
}

// LoadApp 加载指定ID应用
func LoadApp(appID int32) (a App, err error) {
	if appID <= 0 {
		return nil, ErrInvalidAppID
	}
	if a, ok := apps[appID]; ok {
		return a, nil
	} else {
		return nil, ErrInvalidAppID
	}
}

// ID implements App.
func (a *app) ID() int32 {
	return a.appID
}

// Name implements App.
func (a *app) Name() string {
	return a.name
}
