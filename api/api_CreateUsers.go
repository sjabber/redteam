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

	num := 7
	num, err = user.CreateUsers()
	if err != nil { //user.CreateUser 에 대한 오류검사
		switch num {
		case 0:
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})

		case 1:
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"오류": err.Error()})

		case 2:
			log.Println(err.Error())
			c.JSON(http.StatusPaymentRequired, gin.H{"오류": err.Error()})

		case 3:
			log.Println(err.Error())
			c.JSON(http.StatusForbidden, gin.H{"오류": err.Error()})

		case 5:
			log.Println(err.Error())
			c.JSON(http.StatusMethodNotAllowed, gin.H{"오류": err.Error()})
		}
		return

	}

	accessToken, refreshToken, err := user.GetAuthToken()
	if err == nil { // JWToken 에 대한 오류검사
		c.JSON(http.StatusOK, gin.H{
			"access-token": accessToken,
			"refresh-token": refreshToken,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"사용자 아이디": user.Email})
}
