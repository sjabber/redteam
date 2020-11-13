package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)


func CreateUser(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil { //model.User 바인딩 오류검사
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	num := 200
	num, err = user.CreateUsers()
	if err != nil { //user.CreateUser 에 대한 오류검사
		log.Println(err.Error())
		c.JSON(num, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := user.GetAuthToken()
	if err == nil { // JWToken 에 대한 오류검사
		c.JSON(http.StatusOK, gin.H{
			"access-token": accessToken,
			"refresh-token": refreshToken,
		})
		return
	} else {
		c.JSON(500, gin.H{"Create error occurred, account : ": user.Email})
		return
	}

}
