// config/initializers/kafka.go

package initializers

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func ConnectKafka(port string, topic string) {
	KafkaWriter = &kafka.Writer{
		Addr:  kafka.TCP(fmt.Sprintf("localhost:%s", port)),
		Topic: topic,
	}
}
