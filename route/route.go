package route

import (
	"github.com/gin-gonic/gin"
	"redteam/api"
	"redteam/middleware"
	"redteam/model"
)


func RegisterRoute(r *gin.Engine) {

	conn, err := model.ConnectDb()
	if err != nil {
		return
	}

	r.Use(middleware.DBMiddleware(*conn))

	apiV1 := r.Group("/api")
	{
		apiV1.GET("/Time", api.Time)
		apiV1.POST("/Login", api.Login)
		//apiV1.GET("/logout", api.Logout)
		//apiV1.GET("/checklogin", api.CheckLogin)
		apiV1.POST("/CreateUser", api.CreateUser)
		//apiV1.GET("/RefreshToken", api.RefreshToken)
	}

	setting := r.Group("/setting")
	{
		setting.POST("/PostTemplateList", middleware.AuthMiddleWare(), api.PostTemplateList)
		setting.GET("/GetTemplateList", api.GetTemplateList)
	}

	target := apiV1.Group("/target")
	target.Use(middleware.AuthMiddleWare())

	r.Run(":5000")
}
