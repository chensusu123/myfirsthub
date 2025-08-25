package idgenerator

import (
	"github.com/sony/sonyflake/v2"
)

var (
	flake *sonyflake.Sonyflake
)

func init() {
	var err error
	flake, err = sonyflake.New(sonyflake.Settings{})
	if err != nil {
		panic(err)
	}
}

// NextID generates a next unique ID as uint64.
func NextID() (uint64, error) {
	i, err := flake.NextID()
	if err != nil {
		return 0, err
	}
	return uint64(i), nil
}
