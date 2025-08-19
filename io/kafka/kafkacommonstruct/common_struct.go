package kafkacommonstruct

import (
	"encoding/json"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
)

type ProductID uint32

func (f ProductID) MarshalJSON() ([]byte, error) {
	return json.Marshal(100010)
}

type ServerID int32

func (s ServerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(appconfig.GlobalConfig().Global.SectionID)
}

type KafkaCommon struct {
	DataBase  string    `json:"database"`   // 库名
	Table     string    `json:"table"`      // 表名
	ProductID ProductID `json:"product_id"` // 产品ID
	ServerID  ServerID  `json:"server_id"`  // 区服id
}
