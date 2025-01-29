package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
)

var (
	topicName = "createchat_topic"
	producer  sarama.AsyncProducer
)

func Init(brokers ...string) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = time.Second * 5
	config.Producer.Retry.Max = 3
	var err error
	producer, err = sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic("failed to create AsyncProducer" + err.Error())
	}
	go func() {
		for success := range producer.Successes() {
			logger.LogKafkaSuccess(int(success.Partition), int(success.Offset))
		}
	}()
	go func() {
		for err := range producer.Errors() {
			logger.LogKafkaError(err)
		}
	}()
}

func Close() {
	if producer != nil {
		producer.Close()
	}
}

func CreateChatEvent(chatUUID string) {
	message := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder("chat_uuid"),
		Topic: topicName,
		Value: sarama.StringEncoder(chatUUID),
	}
	producer.Input() <- message
}
