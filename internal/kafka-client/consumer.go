package kafkaclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"example_project/internal/database"
	discountaggregator "example_project/internal/discount_aggregator"
	"example_project/internal/leveledlog"
	"example_project/internal/structs"
)

func (k *Kafka) ConsumeEvent(doneChan chan bool) *kafka.Consumer {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	logsChan := make(chan kafka.LogEvent, 1000)
	configMap := k.creteConsumerConfig(logsChan)
	consumerInstance, err := kafka.NewConsumer(&configMap)
	if err != nil {
		logger.Fatal("failed to create consumer: %s", err.Error())
	}

	err = consumerInstance.SubscribeTopics([]string{k.Config.MftTopicName}, nil)
	if err != nil {
		logger.Fatal("failed to subscribe to topic, topic_name: %s, reason: %s", k.Config.MftTopicName, err.Error())
	}
	// consumed-event handler
	go func(done chan bool) {

		for e := range consumerInstance.Events() {
			k.handleConsumerEvent(e)
		}
		logger.Warning("consumer handle event routine stopped")
		done <- true
	}(doneChan)

	// consumer log handler
	// go func() {
	// 	for e := range c.Logs() {
	// 		logger.KafkaLog(e.Message)
	// 	}
	// }()
	return consumerInstance
}

func (k *Kafka) creteConsumerConfig(logsChan chan kafka.LogEvent) kafka.ConfigMap {
	configMap := make(kafka.ConfigMap)
	configMap["bootstrap.servers"] = k.Config.KafkaBrokerServers
	configMap["group.id"] = k.Config.KafkaGroupID
	configMap["client.id"] = fmt.Sprintf("discount_aggregator_host_%v", os.Getenv("HOSTNAME"))
	configMap["auto.offset.reset"] = "earliest"
	configMap["go.logs.channel.enable"] = true
	configMap["go.logs.channel"] = logsChan
	configMap["go.events.channel.enable"] = true
	configMap["go.events.channel.size"] = 1000
	leveledlog.Logger.Debug("config map for kafka is : %+v\n", configMap)
	return configMap
}

func (k *Kafka) handleConsumerEvent(e kafka.Event) {
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
		} else {
			// commit event of kafka
			fmt.Printf("consumed event to topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			var payload structs.TransactionCompletePayload
			json.Unmarshal(ev.Value, &payload)
			k.Wg.Add(1)
			discountLogTable := database.DiscountLogTable(*k.Db)
			unprocessedTransactionLogTable := database.UnprocessedTransactionLogTable(*k.Db)
			merchantMapTable := database.MerchantMapTable(*k.Db)

			trxEventDb := discountaggregator.DbRepo{
				DiscountLog:               &discountLogTable,
				UnprocessedTransactionLog: &unprocessedTransactionLogTable,
				MerchantMap:               &merchantMapTable,
			}
			discountaggregator.HandleDiscountEvent(k.Wg, nil, 0, trxEventDb, &payload)
		}
	default:
		fmt.Printf("Ignored event: %s\n", ev)
	}
}
