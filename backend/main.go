package main

import (
	_ "net/http"

	"redteam/backend/route"

	"github.com/gin-gonic/gin"
)

func main() {

	route.RegisterRoute(gin.Default())

}






