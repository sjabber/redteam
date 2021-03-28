package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"net/http"
	"redteam/model"
	"time"
)

// 미들웨어
// DB의 key 값을 db로 설정한다.
func DBMiddleware(conn sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

// 헤더 HTTP 리퀘스트시 헤더를 꾸며주는 친구
// 크로스도메인 보안
func SetHeader(c *gin.Context) {
	// c.Header(cors.AllowOriginKey, "http://localhost:8888") 여기에 지정된것만 접근이 가능하도록 되어있다. -> 신뢰된 사이트
	// 지정되어 있지 않은 주소가 접근하면 CORS (Cross Origin Resource Sharing 문제를 야기한다.)
	// CSRF, XSS 등의 공격을 막기 위함.
	c.Header(cors.AllowOriginKey, "http://localhost:8888")
	c.Header(cors.AllowCredentialsKey, "true")
	c.Header(cors.AllowMethodsKey, "GET, POST, PUT, OPTIONS, DELETE")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	c.Header("Last-Modified", time.Now().String())
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "-1")
	if c.Request.Method == "OPTIONS" {
		println("OPTIONS")
		c.JSON(http.StatusOK,
			gin.H{"status": http.StatusOK})
		c.Abort()
		return
	}
}

func TokenAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer, err := c.Request.Cookie("access-token")
		if err != nil {
			model.SugarLogger.Errorf("cookie error : %v", err.Error())

			// access-token 을 읽어오는데 오류가 발생할 경우
			// recover.go 로 이동하기전에 403 에러를 반환시켜버림.
			c.AbortWithStatus(http.StatusForbidden)
			c.Abort()
			return
		}

		isValid, user := model.IsTokenValid(bearer.Value)
		if !isValid || bearer == nil {
			// access-token 을 검증할 때 false (유효시간 만료 등)면
			// recover.go 로 이동하기전에 403 에러를 반환시켜버림.
			c.AbortWithStatus(http.StatusForbidden)
			c.Abort()
			return
		} else {
			c.Set("number", user.UserNo)
			c.Set("email", user.Email)
			c.Set("name", user.Name)
			c.Next()
		}
	}
}

// 기존 함수 검증 미들웨어
//func TokenAuthMiddleWare() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		bearer, err := c.Request.Cookie("access-token")
//		if err != nil {
//			log.Println(err) // named cookie not present
//			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
//			// access-token 을 읽어오는데 오류가 발생할 경우 403에러를 반환한다.
//			c.Abort()
//		}
//
//		isValid, user := model.IsTokenValid(bearer.Value)
//		if isValid == false || bearer == nil {
//			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
//			// access-token 을 검증할 때 false (유효시간 만료 등)면 403에러를 반환한다.
//			c.Abort()
//		} else {
//			c.Set("number", user.UserNo)
//			c.Set("email", user.Email)
//			c.Set("name", user.Name)
//			c.Next()
//		}
//	}
//}
