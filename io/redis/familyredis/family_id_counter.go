package familyredis

import (
	"context"
	globalredis "maze_game_server/io/redis"
)

func getFamilyIdCounterKey() string {
	return "family:id:counter"
}

func CreateFamilyId() (familyId int32, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return
	}
	ret, err := db.Incr(context.TODO(), db.MakeSectionKey(getFamilyIdCounterKey())).Result()
	return int32(ret), err
}
