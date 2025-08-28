package kafkacommonstruct

import (
	"encoding/json"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
)

type ProductID uint32

func (f ProductID) MarshalJSON() ([]byte, error) {
	// 爱地牢产品id 100010
	return json.Marshal(100010)
}

type ServerID int32

func (s ServerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(appconfig.GlobalConfig().Global.SectionID)
}

type ServiceName string

func (s ServiceName) MarshalJSON() ([]byte, error) {
	return json.Marshal(appconfig.GlobalConfig().Server.AppName)
}

type CreateTime uint64

func (t CreateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Now().UnixNano() / 1000000)
}

type KafkaCommon struct {
	ServiceName ServiceName `json:"service_name"`
	DataBase    string      `json:"database"`    // 库名
	Table       string      `json:"table"`       // 表名
	ProductID   ProductID   `json:"product_id"`  // 产品ID
	ServerID    ServerID    `json:"server_id"`   // 区服ID
	GroupID     ServerID    `json:"group_id"`    // 分组ID
	CreateTime  CreateTime  `json:"create_time"` // 毫秒级时间戳 ms
}
