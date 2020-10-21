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
	r.Use(middleware.SetHeader)

	apiV1 := r.Group("/api")
	{
		apiV1.GET("/Time", api.Time)
		apiV1.POST("/Login", api.Login)
		//apiV1.GET("/logout", api.Logout)
		//apiV1.GET("/checklogin", api.CheckLogin)
		// checkdLogin 왜만드냐 -> 토큰방식 서로 각자 관리한다(토큰이 헤더에 계속 붙어있는것이 아님)
		// 요청이 올때마다 -> 다른 페이지, api 요청할때마다 토큰을 검증한다.
		// 우리가 데시보드라는 화면을 들어갔다 -> 데이터를

		apiV1.POST("/CreateUser", api.CreateUser)
		//apiV1.GET("/RefreshToken", api.RefreshToken)
	}

	setting := r.Group("/setting")
	setting.Use(middleware.AuthMiddleWare()) //로그인 이후에 사용할 api 들은 토큰검증이 필요
	{
		setting.POST("/PostTemplateList", api.PostTemplateList)
		setting.GET("/GetTemplateList", api.GetTemplateList)
	}

	// 대시보드
	apiV2 := r.Group("/api")
	apiV2.Use(middleware.AuthMiddleWare())
	{
		apiV2.GET("/dashboard", api.Dashboard)
	}

	target := apiV1.Group("/target")
	target.Use(middleware.AuthMiddleWare())

	r.Run(":5000")
}
