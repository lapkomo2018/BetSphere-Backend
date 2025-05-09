package log

import (
	"github.com/IBM/sarama"
)

type KafkaHook struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaLogger(brokers []string, topic string) (*Logger, error) {
	hook, err := NewKafkaHook(brokers, topic)
	if err != nil {
		return nil, err
	}
	logger := New()
	logger.AddHook(hook)
	return logger, nil
}

func NewKafkaHook(brokers []string, topic string) (*KafkaHook, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &KafkaHook{producer: p, topic: topic}, nil
}

func (k *KafkaHook) Fire(entry *Entry) error {
	b, err := entry.Bytes()
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(b),
	}
	_, _, err = k.producer.SendMessage(msg)
	return err
}

func (k *KafkaHook) Levels() []Level {
	return AllLevels
}
