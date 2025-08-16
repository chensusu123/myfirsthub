package kafka

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanokafkaproducer"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
)

type KafkaProducer struct {
	serverdepend.DependInit
	database.KafkaProducer
}

var GflowKafka *KafkaProducer

func NewKafkaProduce(serviceName string, name string) *KafkaProducer {
	rt := &KafkaProducer{}
	xx := nanokafkaproducer.NewKafkaProducer(serviceName, name)
	rt.KafkaProducer = xx
	rt.DependInit = xx
	return rt
}

func init() {
	GflowKafka = NewKafkaProduce("maze_main_server_flow.kafka", "game_flow_produce")
	serverdepend.RegisterDepend(GflowKafka)
}
