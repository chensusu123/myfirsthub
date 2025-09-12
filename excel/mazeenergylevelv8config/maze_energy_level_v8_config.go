package mazeenergylevelv8config

func GetKey(energyId int32, level int32) int32 {
	return energyId*10000 + level
}
