package route

import (
	"github.com/gin-gonic/gin"
	"redteam/api"
)

func RegisterRoute(r *gin.Engine)  {
	apiV1 := r.Group("/api")

	apiV1.GET("/Time", api.Time)

	r.Run(":5000")
}