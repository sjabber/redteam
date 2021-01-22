package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	appKafka "redteam/api"
	crontab "redteam/model"
	"redteam/route"
)

func main() {
	go crontab.AutoStartProject()
	go appKafka.Consumer()
	route.RegisterRoute(gin.Default())
}
