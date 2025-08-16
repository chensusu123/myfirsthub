package kafkacommonstruct

type KafkaCommon struct {
	DataBase  string `json:"database"`
	Table     string `json:"table"`
	SectionID string `json:"section_id"`
}
