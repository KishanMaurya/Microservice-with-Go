package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func InitProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	Producer, err = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal("Kafka producer error:", err)
	}

	log.Println("✅ Kafka Producer connected")
}

func Publish(topic string, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := Producer.SendMessage(msg)
	if err != nil {
		log.Println("❌ Kafka publish error:", err)
	}
	log.Printf("Sent to partition %d at offset %d", partition, offset)
}