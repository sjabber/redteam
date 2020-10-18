package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func Login(c *gin.Context) {


	user := model.User{}
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
		return
	}

	// github.com/dgrijalva/jwt-go
	//로그인 자격증명을 검사한다.
	//db, _ := c.Get("db")
	//conn := db.(sql.DB)
	//err = user.IsAuthenticated(&conn) // 비밀번호 확인
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"오류": err.Error()})
	//	return
	//}
	// 1800 -> 30분

	err = user.Login()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest,
			gin.H{"오류": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {
		c.SetCookie("access-token", token, 1800, "", "", false, false)
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"isOk": false,
		"오류":   "인증에 오류가 발생하였습니다. ",
		//
	})
	//
	//c.JSON(http.StatusOK,
	//	gin.H{"로그인 상태": true, "ID": user.GetID(),
	//		"name": user.GetName() +"님 반갑습니다."})
	return

}
