package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func GetTag(c *gin.Context) {
	num := c.Keys["number"].(int)

	c.JSON(http.StatusOK, gin.H{
		"isOk":   1,
		"status": http.StatusOK,
		"tags":   model.GetTag(num), // 태그들
	})
}
