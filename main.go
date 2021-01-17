package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	appKafka "redteam/api"
	"redteam/route"
)

func main() {
	go appKafka.Consumer()
	route.RegisterRoute(gin.Default())
}
