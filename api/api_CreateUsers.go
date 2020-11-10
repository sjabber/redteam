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

		//case 400: // 정보를 전부 입력하지 않은 경우
		//case 401: // 이미 존재하는 계정이 있을 경우
		//case 402: // 비밀번호나 이메일 형식이 올바르지 않은 경우
		//case 403: // 비밀번호와 확인값이 일치하지 않는 경우
		//case 405: // 계정을 생성하는 도중 오류가 발생한 경우
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

	c.JSON(http.StatusOK, gin.H{"User's account": user.Email})
}
