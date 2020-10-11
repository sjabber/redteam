package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
	"strings"
)

// 미들웨어
// DB의 key 값을 db로 설정한다.
func DBMiddleware(conn sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Taeho ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"오류": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]
		//fmt.Printf("Bearer (%v) \n", token)
		isValid, userID := model.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("email", userID)
			c.Next()
		}
	}
}