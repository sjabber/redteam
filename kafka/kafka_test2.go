package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func Configure() {
	topic := "redteam"
	partition := 1

	connect, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	connect.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err = connect.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
		)
	if err != nil {
		log.Fatal("failed to write message:", err)
	}

	if err := connect.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

////
	w := kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		Topic: "redteam",
	}

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte("key-A"),
			Value: []byte("Hello World!"),
		})

}