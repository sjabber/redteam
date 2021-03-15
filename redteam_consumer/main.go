package main

import (
	"fmt"
	"log"
	appKafka "redteam_consumer/api"
	crontab "redteam_consumer/model"
)

func main() {
	// 스케쥴링 & 카프카 모두 비동기로 작동시킨다.
	go crontab.AutoStartProject()
	go appKafka.Consumer()

	// 비동기 프로세스를 지속시키기 위한 입력문
	// end 를 치면 프로세스 종료됨.
	var command string
	fmt.Scanln(&command)

	Command(command)
}

func Command(command string) {
	if command != "end" {
		fmt.Scanln(&command)
		Command(command)
	}  else {
		log.Println("Kafka server shutdown")
	}
}