package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Logout(c *gin.Context) {

	_, err := c.Cookie("access-token")
	if err == nil {
		log.Println("logout true")
		c.SetCookie("access-token", "", -1, "", "", false, true)
		c.SetCookie("refresh-token", "", -1, "", "",false, true)
		c.JSON(http.StatusOK, gin.H{"LogoutSuccess" : true})
		return
	} else {
		//로그아웃에 오류가 발생할경우 500에러를 반환한다.
		log.Println("Logout error occurred, account :", Account)
		c.AbortWithStatus(http.StatusInternalServerError) //500에러
		return
	}

}

