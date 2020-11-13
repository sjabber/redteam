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
		c.JSON(http.StatusOK, gin.H{
			"LogoutSuccess" : true,
		})
		return

	} else {
		log.Println("logout fail")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}


}

