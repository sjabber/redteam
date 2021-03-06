package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func ResultDetail(c *gin.Context) {
	// 계정번호
	num := c.Keys["number"].(int)

	// db연결
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	resultDetail, err := model.GetResultDetail(c.Query("p_no"), num, &conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isOk":         1,
		"status":       http.StatusOK,
		"resultDetail": resultDetail,
	})
}
