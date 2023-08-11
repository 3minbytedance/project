package kafka

import (
	"context"
	"encoding/json"
	"project/models"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaManager *Manager

type Manager struct {
	Brokers []string
}

func NewKafkaManager(brokers []string) *Manager {
	return &Manager{
		Brokers: brokers,
	}
}

func (m *Manager) NewProducer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(m.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
}

func (m *Manager) NewConsumer(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   m.Brokers,
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
		GroupID:   "messages",
	})
}

func (m *Manager) ProduceMessage(producer *kafka.Writer, message *models.Message) error {
	messageBytes, err := json.Marshal(*message)
	if err != nil {
		return err
	}
	return producer.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})
}
