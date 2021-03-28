package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func CreateUser(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil { //model.User 바인딩 오류검사
		model.SugarLogger.Errorf("JSON binding error : %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk": false,
		})
		return
	}

	num, err := user.CreateUsers()
	if err != nil { //user.CreateUser 에 대한 오류검사
		c.JSON(num, gin.H{
			"isOk": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isOk": true,
		"You are now a member, account ": user.Email,
	})
}
