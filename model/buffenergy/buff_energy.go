package buffenergy

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

func getKey(userID uint64, barrierID, areaID, areaIndex int32) string {
	return fmt.Sprintf("tmpbuff:energy:u:%d:barrier:%d:area:%d:areaIndex:%d", userID, barrierID, areaID, areaIndex)
}

type BuffEnergy struct {
	NowEnergy int32
}

func NewBuffEnergy(ctx context.Context, userID uint64, barrierID, areaID, areaIndex int32) (*BuffEnergy, error) {
	nowBuffEnergy := &BuffEnergy{}
	if err := nowBuffEnergy.load(ctx, userID, barrierID, areaID, areaIndex); err != nil {
		return nil, err
	}

	return nowBuffEnergy, nil
}

func (s *BuffEnergy) load(ctx context.Context, userID uint64, barrierID, areaID, areaIndex int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierID, areaID, areaIndex), s)
}

func (s *BuffEnergy) Save(ctx context.Context, userID uint64, barrierID, areaID, areaIndex int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierID, areaID, areaIndex), s)
}

func (s *BuffEnergy) Del(ctx context.Context, userID uint64, barrierID, areaID, areaIndex int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierID, areaID, areaIndex))
}
