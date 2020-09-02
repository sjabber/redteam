package route

import (
	"github.com/gin-gonic/gin"
	"redteam/backend/api"
)

func RegisterRoute(r *gin.Engine) {
	apiV1 := r.Group("/api")
	apiV1.GET("/Time", api.Time)
	apiV1.POST("/Login", api.Login)
	apiV1.POST("/UserRegister", api.UserRegister)

	r.Run(":5000")
}
