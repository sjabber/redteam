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
		c.JSON(http.StatusBadRequest, gin.H{"target_registration_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"target_registration_error, register_account" : c.Keys["email"]})
	}
}

func GetTarget(c *gin.Context) {
	target, err := model.ReadTarget() //DB에 저장된 대상들을 읽어오는 메서드
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"targets": target, "register_account": c.Keys["email"]})
}

func DeleteTarget(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)
	err := target.DeleteTarget(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"target_deleting_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"delete_success, deleting_account" : c.Keys["email"]})
	}
}

func DownloadExcel(c *gin.Context) {
	err := model.Download()
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"file_download_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"file_down_OK, download_account": c.Keys["email"]})
	}
}

