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
	c.Header(cors.AllowOriginKey, "http://localhost:63343") // 여기에 지정된것만 접근이 가능하도록 되어있다. -> 신뢰된 사이트
	// 지정되어 있지 않은 주소가 접근하면 CORS (Cross Origin Resource Sharing 문제를 야기한다.)
	// CSRF, XSS 등의 공격을 막기 위함.
	//("http://ip주소 + 5000")
	//c.Header(cors.AllowOriginKey, "http://121.173.129.251")
	//msa
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
			c.JSON(http.StatusUnauthorized, gin.H{"오류": "Not authenticated.",
				"err": err.Error()})
			c.Abort()
		}

		isValid, user := model.IsTokenValid(bearer.Value)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"오류": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("email", user.Email)
			c.Set("name", user.Name)
			c.Next()
		}
	}
}
