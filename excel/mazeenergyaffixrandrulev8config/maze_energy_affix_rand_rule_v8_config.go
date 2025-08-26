package mazeenergyaffixrandrulev8config

func GetKey(randId int32, level int32) int32 {
	return randId*10000 + level
}
