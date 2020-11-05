package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func Login(c *gin.Context) {
	num := 0
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

	num, err = user.Login()
	if err != nil && num == 0 {
		// num 0 이면서 err == nil 이어야만 로그인에 성공한다.
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest,
			gin.H{"오류": err.Error()})
		return
	} else if err != nil && num == 1 {
		// 로그인이 실패한 경우
		log.Println(err.Error())
		c.JSON(http.StatusUnauthorized,
			gin.H{"오류": err.Error()})
		return
	}

	accessToken, refreshToken, err := user.GetAuthToken()
	if err == nil { //여기서 토큰을 쿠키에 붙인다.
		c.SetCookie("access-token", accessToken, 900, "", "", false, true)
		c.SetCookie("refresh-token", refreshToken, 86400, "", "", false, true)
		//https 사용시 refresh-token 의 secure -> true 로 변경한다.

		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"isOk": false,
		"error":   "authentication error occurred. ",
	})
	//
	//c.JSON(http.StatusOK,
	//	gin.H{"로그인 상태": true, "ID": user.GetID(),
	//		"name": user.GetName() +"님 반갑습니다."})
	return

}
