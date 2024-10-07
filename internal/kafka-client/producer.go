package kafkaclient

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"example_project/internal/config"
	"example_project/internal/leveledlog"
)

func ProducerEvent(cfg *config.Config, done chan bool) *kafka.Producer {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	logsChan := make(chan kafka.LogEvent, 1000)

	producerConf := createProducerConfig(cfg, logsChan)
	producerInstance, err := kafka.NewProducer(&producerConf)
	if err != nil {
		logger.Fatal("Failed to create producer: %s", err)
	}

	go func() {
		for e := range producerInstance.Events() {
			handleProducerEvent(e)
		}
		done <- true
	}()

	// go func() {
	// 	for e := range producerInstance.Logs() {
	// 		logger.KafkaLog(e.Message)
	// 	}
	// }()

	// users := [...]string{"eabara", "jsmith", "sgarcia", "jbernard", "htanaka", "awalther"}
	// items := [...]string{"book", "alarm clock", "t-shirts", "gift card", "batteries"}

	// for n := 0; n < 10; n++ {
	// 	key := users[rand.Intn(len(users))]
	// 	data := items[rand.Intn(len(items))]
	// 	p.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 		Key:            []byte(key),
	// 		Value:          []byte(data),
	// 	}, nil)

	// }
	return producerInstance

}

func createProducerConfig(cfg *config.Config, logsChan chan kafka.LogEvent) kafka.ConfigMap {
	producerConf := make(kafka.ConfigMap)
	producerConf["bootstrap.servers"] = cfg.KafkaBrokerServers
	producerConf["go.logs.channel.enable"] = true
	producerConf["go.logs.channel"] = logsChan
	producerConf["go.events.channel.size"] = 1000
	return producerConf
}

func handleProducerEvent(e kafka.Event) {
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
		} else {
			fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}
}
