package main

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"redteam/route"
)

func main() {


	route.RegisterRoute(gin.Default())


}






