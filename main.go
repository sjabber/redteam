package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"redteam/model"
	"redteam/route"
)

func main() {
	//go crontab.AutoStartProject() -> redteam_prodcuer 로 이동함.
	//go appKafka.Consumer() -> Springboot 가 담당함.

	model.InitLogger()
	defer model.SugarLogger.Sync()
	route.RegisterRoute(gin.Default())
}