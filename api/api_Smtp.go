package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func SmtpConnectionCheck(c *gin.Context)  {
	var sm model.Smtpinfo

	num := c.Keys["number"].(int)

	// db연결
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	//err := c.BindJSON(&sm)
	//err = sm.IdPwCheck(&conn)
	err := sm.SmtpConnectionCheck(&conn, num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"isOK": 0,
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"isOK": 1,
	})
}

//func SendMail(c *gin.Context) {
//	var sm model.Smtpinfo
//	err := c.BindJSON(&sm)
//	err = sm.SendMail()
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"status": http.StatusBadRequest,
//			"isOk": 0,
//			"error": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"status": http.StatusOK,
//		"isOk": 1,
//	})
//}
