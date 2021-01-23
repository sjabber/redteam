package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)
// 토큰안에 이름도 넣는다.
var Account interface{}

func Dashboard(c *gin.Context) {

	Account = c.Keys["email"]
	//if Account == nil {
	//	c.JSON(http.StatusForbidden, gin.H{})
	//	return
	//}

	// email, name 을 출력할 수 있도록 만든다.
	c.JSON(http.StatusOK, gin.H{"email": c.Keys["email"], "name": c.Keys["name"]})

	return
}

// 맨위 전체 진행상황
func GetDashboardInfo1(c *gin.Context) {
	// 계정번호
	num := c.Keys["number"].(int)

	// db연결
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	Info1, err := model.GetDashboardInfo1(&conn, num)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"isOk": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"info1" : Info1,
		})
	}
}

// 진행중인 프로젝트 현황
//func GetDashboardInfo2(c *gin.Context) {
//	// 계정번호
//	num := c.Keys["number"].(int)
//
//	// db연결
//	db, _ := c.Get("db")
//	conn := db.(sql.DB)
//
//	Info1, err := model.GetDashboardInfo2(&conn, num)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"status": http.StatusInternalServerError,
//			"isOk": 0,
//		})
//		return
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"info2" : Info1,
//		})
//	}
//
//
//	return
//}

// 맨아래 전체 프로젝트 리스트
func OverallList() {

}