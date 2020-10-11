package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func Login(c *gin.Context) {

	// github.com/dgrijalva/jwt-go

	user := model.User{}
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
		return
	}

	//로그인 자격증명을 검사한다.
	db, _ := c.Get("db")
	conn := db.(sql.DB)
	err = user.IsAuthenticated(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"토큰" : token,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"오류": "인증에 오류가 발생하였습니다. ",
	})

	err = user.Login()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest,
			gin.H{"오류": err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"로그인 상태": true, "ID": user.GetID(),
			"name": user.GetName() +"님 반갑습니다."})
	return

}
