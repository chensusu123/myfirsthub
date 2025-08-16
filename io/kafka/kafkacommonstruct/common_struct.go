package kafkacommonstruct

import "encoding/json"

type ProductID uint32

func (f ProductID) MarshalJSON() ([]byte, error) {
	return json.Marshal(100010)
}

type KafkaCommon struct {
	DataBase  string    `json:"database"`   // 库名
	Table     string    `json:"table"`      // 表名
	SectionID string    `json:"section_id"` // 区服ID
	ProductID ProductID `json:"product_id"` // 产品ID
}
