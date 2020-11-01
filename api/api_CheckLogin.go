package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckLogin(c *gin.Context) {
	// route에 미들웨어 붙여서 사용할거라 제외
	c.JSON(http.StatusOK, gin.H{"email": c.Keys["email"], "name": c.Keys["name"]})
}