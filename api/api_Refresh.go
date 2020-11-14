package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func RefreshToken(c *gin.Context) {
	// refresh-token 쿠키를 요청한다.
	bearer, err := c.Request.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated.",
			"err": err.Error()})
		c.Abort()
	}
	// 해당 쿠키(refresh token)의 값을 검사한다.
	// 검사에 통과하면 GetAccessToken 메서드로 access-token 을 재발급 받는다.
	isValid, user := model.RefreshTokenValid(bearer.Value)
	if isValid == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated."})
		// access-token을 검증할 때 false (유효시간 만료 등)면 403에러를 반환한다.
		// 이 경우 전에 access 토큰을 한번 이상 재발급한 적이 있음.
		c.Abort()
	} else {
		//검증이 통과한 경우 반환받은 구조체안의 정보를 해당 키값에 세팅
		c.Set("email", user.Email)
		c.Set("name", user.Name)
		c.Set("number", user.UserNo)
		c.Next()

		// todo 여기에 위의 계정이 맞는지 검사하는 로직도 넣어주자.
		accessToken, err := user.GetAccessToken()
<<<<<<< HEAD
		if err == nil {
			c.SetCookie("access-token", accessToken, 5, "", "", false, true)
=======
		if err == nil { //여기서 토큰을 쿠키에 붙인다.
			c.SetCookie("access-token", accessToken, 900, "", "", false, true)
>>>>>>> 9a11b5c7bd79e3b4db31f3bc9b140909b433d657
			c.JSON(http.StatusOK, gin.H{
				"isOk": true,
			})
			log.Print("token refresh")
			return

		} else {
			// access 토큰이 발급되지 않은 경우 405에러를 반환한다.
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"isOk": false,
				"error": "authentication error occurred. ",
			})
			return
		}

		return
	}
}
