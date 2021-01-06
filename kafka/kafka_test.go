package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func StartKafka() {
	conf := kafka.ReaderConfig{
		Brokers: []string {"localhost:2181"},
		Topic: "message",
		GroupID: "G1",
		MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)

	// 루프문은 이 브로커가 작동하고 있는지 점검한다.
	for {
		//포트 내 컨텍스트 개체는 해당 메시지를 읽는 동작
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Some error occured", err)
			continue
		}
		fmt.Println("Message is : ", string(m.Value))
	}
}