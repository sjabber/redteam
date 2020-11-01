package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func RegTarget(c *gin.Context) {
	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)
	err := target.CreateTarget(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"훈련대상 등록 오류": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"대상 등록 성공, 등록한 관리자": c.Keys["email"]})
	}
}

func GetTarget(c *gin.Context) {
	target, err := model.ReadTarget() //DB에 저장된 대상들을 읽어오는 메서드
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"targets": target, "id": c.Keys["email"]})
}

func DeleteTarget(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)
	err := target.DeleteTarget(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"훈련대상 삭제 오류": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"대상 삭제 성공, 삭제한 관리자": c.Keys["email"]})
	}
}

func DownloadExcel(c *gin.Context) {
	err := model.Download()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"파일 다운로드 오류": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"파일 다운로드 성공, 다운로드 계정": c.Keys["email"]})
	}
}
