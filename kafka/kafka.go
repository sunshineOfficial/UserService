package kafka

import (
	"user-service/ctx"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Message = kafka.Message

type Kafka interface {
	Producer(topicName string, options ...ProducerOption) Producer
	Consumer(log *zap.Logger, getCtx ctx.ProvideWithCancel, options ...ConsumerOption) (Consumer, error)
}

type KafkaImpl struct {
	brokers []string
}

func NewKafka(brokers []string) *KafkaImpl {
	return &KafkaImpl{
		brokers: brokers,
	}
}

func (k *KafkaImpl) Producer(topicName string, options ...ProducerOption) Producer {
	return NewProducer(k.brokers, topicName, options...)
}

func (k *KafkaImpl) Consumer(log *zap.Logger, getCtx ctx.ProvideWithCancel, options ...ConsumerOption) (Consumer, error) {
	return NewConsumer(log, getCtx, k.brokers, options...)
}
