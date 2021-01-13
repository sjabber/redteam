package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"redteam/route"
)

func main() {
	route.RegisterRoute(gin.Default())

	//ctx := context.Background()
	//go kafka.Produce(ctx)
	//kafka.Consume(ctx)
}
