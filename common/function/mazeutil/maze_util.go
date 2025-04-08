package mazeutil

func CheckMazeEquip(equipId int32) bool {
	if equipId%1000 < 200 {
		return true
	}
	return false
}
