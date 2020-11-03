package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func RefreshToken(c *gin.Context) {
	// refresh-token 쿠키를 요청한다.
	bearer, err := c.Request.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"오류": "Not authenticated.",
			"err": err.Error()})
		c.Abort()
	}
	// 해당 쿠키(refresh token)의 값을 검사한다.
	// 검사에 통과하면
	isValid, user := model.RefreshTokenValid(bearer.Value)
	if isValid == false {
		c.JSON(http.StatusUnauthorized, gin.H{"오류": "Not authenticated."})
		c.Abort()
	} else {
		//검증이 통과한 경우 반환받은 구조체안의 정보를 해당 키값에 세팅
		c.Set("email", user.Email)
		c.Set("name", user.Name)
		c.Next()

		// 여기에 위의 계정이 맞는지 검사하는 로직도 넣어주자.

		accessToken, err := user.GetAccessToken()
		if err == nil { //여기서 토큰을 쿠키에 붙인다.
			c.SetCookie("access-token", accessToken, 900, "", "", false, false)
			c.JSON(http.StatusOK, gin.H{
				"isOk": true,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"isOk": false,
			})
		}
	}













}
