package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

func StartConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal(err)
	}

	partitionConsumer, err := consumer.ConsumePartition("user.login", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range partitionConsumer.Messages() {

			log.Println("🔍 RAW MESSAGE:", string(msg.Value)) // 👈 ADD THIS

			var event UserLoginEvent

			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("❌ JSON parse error:", err)
				continue
			}

			log.Printf(
				"📩 topic=%s partition=%d offset=%d event=%+v\n",
				msg.Topic,
				msg.Partition,
				msg.Offset,
				event,
			)
		}
	}()
}
