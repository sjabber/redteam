package main

import (
	"github.com/gin-gonic/gin"
	"redteam/route"
)

func main() {

	route.RegisterRoute(gin.Default())

}
