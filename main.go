package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"redteam/route"
)

func main() {
	//go crontab.AutoStartProject()
	//go appKafka.Consumer()
	route.RegisterRoute(gin.Default())
}
