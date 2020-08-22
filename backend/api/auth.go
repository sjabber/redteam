package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"redteam/backend/model"
)

func Login(c *gin.Context) {

	// github.com/dgrijalva/jwt-go

	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest,
			gin.H{"err": err.Error()})
		return
	}

	err = user.Login()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest,
			gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"로그인 상태": true, "ID": user.GetID() ,"name": user.GetName() +"님 반갑습니다."})
	return

}
