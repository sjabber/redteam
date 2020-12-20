package route

import (
	"github.com/gin-gonic/gin"
	"redteam/api"
	"redteam/middleware"
	"redteam/model"
)


func RegisterRoute(r *gin.Engine) {

	conn, err := model.ConnectDB()
	if err != nil {
		return
	}

	r.Use(middleware.DBMiddleware(*conn))
	r.Use(middleware.SetHeader)
	apiV1 := r.Group("/api")
	{
		apiV1.GET("/Time", api.Time)
		apiV1.POST("/Login", api.Login)
		apiV1.GET("/Logout", api.Logout)
		// apiV1.GET("/CheckLogin", api.CheckLogin)
		// checkdLogin 왜만드냐 -> 토큰방식 서로 각자 관리한다(토큰이 헤더에 계속 붙어있는것이 아님)
		// 요청이 올때마다 -> 다른 페이지, api 요청할때마다 토큰을 검증한다.
		// 우리가 데시보드라는 화면을 들어갔다 -> 데이터를

		apiV1.POST("/createUser", api.CreateUser)
		apiV1.GET("/RefreshToken", api.RefreshToken)
	}

	setting := r.Group("/setting")
	setting.Use(middleware.TokenAuthMiddleWare()) //로그인 이후에 사용할 api 들은 토큰검증이 필요
	{
		//setting.POST("/setTemplates", api.PostTemplateList) - 사용안함
		setting.GET("/getTemplates", api.GetTemplateList)
		setting.POST("/EditTemplate", api.EditTemplate)
		setting.POST("/deleteTemplates", api.DeleteTemplateList)
		setting.GET("/TemplateDetail", api.TemplateDetail)
		setting.GET("/getTag", api.GetTag)

		//setting.GET("/userSetting", api.GetUserSetting) //Note - spring boot
		//setting.POST("/userSetting", api.SetUserSetting) //Note - spring boot
		//setting.GET("/smtpSetting", api.GetSmtpSetting) //Note - spring boot
		//setting.POST("/smtpSetting", api.SetSmtpSetting) //Note - spring boot
		//setting.POST("/smtpConnectCheck", api.SmtpConnectionCheck) //Note - spring boot
	}

	// 대시보드
	apiV2 := r.Group("/api")
	apiV2.Use(middleware.TokenAuthMiddleWare())
	{
		apiV2.GET("/dashboard", api.Dashboard)
		//apiV2.POST("/projectCreate", api.ProjectCreate)
		//apiV2.GET("/getProject", api.getProject)
		//apiV2.GET("/endProjectList", api.endProjectList)
		//apiV2.GET("/bookingProjectList", api.BookingProjectList)
	}

	//r.LoadHTMLGlob("./ui/html/target/*")
	target := apiV1.Group("/target")
	target.Use(middleware.TokenAuthMiddleWare())
	//target.Static("/files", "C:/Users/Taeho/go/src/redteam")
	{
		target.GET("/getTarget", api.GetTarget)
		target.POST("/delTarget", api.DeleteTarget)
		target.POST("/regTarget", api.RegTarget)
		target.GET("/exportTarget", api.ExportTarget)
		target.POST("/delTag", api.DeleteTag)
		target.POST("/regTag", api.RegTag)

		target.GET("/downloadExcel", api.DownloadExcel)
		target.POST("/importTargets", api.ImportTargets)

		target.GET("/search", api.Search)
	}

	r.Run(":5000")
}
