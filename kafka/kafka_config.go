package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strconv"
	"time"
)

const (
	topic = "redteam"
	brokerAddress = "localhost:9092"
)

func Produce(ctx context.Context) {
	// 카운터 초기화
	i := 0

	// writer 를 브로커 주소, 토픽과 초기화
	w := kafka.Writer{
		Addr: kafka.TCP(brokerAddress),
		Topic: topic,

	}

	// Note 여기서 프로듀서가 메시지를 생성하는 것이다.
	// 프로젝트가 실행되면 프로듀서의 여기에다가 값을 밀어넣는것이여.
	for {
		// 각각의 카프카 메시지는 키와 밸류를 가진다.
		// 키는 어떤 파티션에 메시지를 적을지 결정한다.

		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// 몇번 적힌것인지 로그를 남긴다.
		fmt.Println("writes:", i)
		i++

		// 1초 슬립.
		time.Sleep(time.Second)
	}
}

func Consume(ctx context.Context) {
	// 새로운 reader를 브로커와 토픽과 함께 초기화한다.
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic: topic,
		GroupID: "redteam",
		MinBytes: 5, //5 바이트
		MaxBytes: 4132,//1kB
		MaxWait: 500 * time.Millisecond, //0.5초만 기다린다.
		StartOffset: kafka.FirstOffset, // GroupID 이전에 동일한 설정으로 데이터 사용한 적이
		// 있는 경우 중단한 곳부터 계속된다.
	})
	for {
		// ReadMessage 메서드는 우리가 다음 이벤트를 받을 때까지 차단된다.
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// 메시지를 받은 다음에는 이것의 값을 기록한다.
		fmt.Println("received: ", string(msg.Value))
	}
}

