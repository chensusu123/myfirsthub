package cluster

import "github.com/google/uuid"

func generatorID() string {
	return uuid.New().String()
}
