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
		c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
		return
	}

	err = user.CreateUser()
	if err != nil { //user.CreateUser 에 대한 오류검사
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil { // JWToken 에 대한 오류검사
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"사용자 아이디": user.Email})
}
