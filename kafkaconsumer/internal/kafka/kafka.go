package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Rolan335/grpcMessenger/kafkaconsumer/internal/repository/inmemory"
	"github.com/Rolan335/grpcMessenger/kafkaconsumer/internal/webhook"
)

var (
	topicName         = "createchat_topic"
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
)

func Init(brokers ...string) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0

	var err error
	consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	partitionConsumer, err = consumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		panic("failed to init partitionConsumer: " + err.Error())
	}
}

func ConsumeCreateChat(webhook *webhook.Caller, storage *inmemory.StorageTTL) {
	for message := range partitionConsumer.Messages() {
		if storage.IsExist(string(message.Value)) {
			fmt.Printf("message %s already processed\n", message.Value)
			continue
		}
		if err := webhook.Call(string(message.Value)); err != nil {
			fmt.Printf("failed to call webhook for %s: %s \n", message.Value, err.Error())
			continue
		}
		storage.Add(string(message.Value))
		fmt.Printf("webhook for %s successfully called\n", message.Value)
	}
}

func Close() {
	if consumer != nil {
		consumer.Close()
	}
	if partitionConsumer != nil {
		partitionConsumer.Close()
	}
}
