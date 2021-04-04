package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

func RefreshToken(c *gin.Context) {
	// refresh-token 쿠키를 요청한다.
	// refresh-token 이 없거나 오류를 발생시킬 경우 401에러를 반환한다.
	bearer, err := c.Request.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk" : false,
		})
		model.SugarLogger.Errorf("cookie error : %v", err.Error())
		c.Abort() // 이후의 핸들러 호출하지 않도록.
		return
	}
	// 해당 쿠키(refresh token)의 값을 검사한다.
	// 검사에 통과하면 GetAccessToken 메서드로 access-token 을 재발급 받는다.
	// access-token 을 검증할 때 false (유효시간 만료 등)면 401에러를 반환한다.
	isValid, user := model.RefreshTokenValid(bearer.Value)
	if isValid == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk" : false,
		})
		c.Abort() // 이후의 핸들러 호출하지 않도록.
	} else {
		// 검증에 통과한 후 쿠키에 access-token 을 붙여준다.
		// 검증이 통과한 경우 반환받은 구조체안의 정보를 해당 키값에 세팅
		c.Set("email", user.Email)
		c.Set("name", user.Name)
		c.Set("number", user.UserNo)
		c.Next()

		// todo 여기에 위의 계정이 맞는지 검사하는 로직도 넣어주자.

		//여기서 토큰을 쿠키에 붙인다.
		accessToken, refreshToken, err := user.GetNewToken()
		if err == nil {
			// 각 1시간, 1주일
			c.SetCookie("access-token", accessToken, 3600, "", "", true, true)
			c.SetCookie("refresh-token", refreshToken, 604800, "", "", true, true)
			c.JSON(http.StatusOK, gin.H{
				"isOk": true,
			})
			//model.SugarLogger.Info("token refresh")
			return
		} else {
			// access 토큰이 발급되지 않은 경우 401에러를 반환한다.
			c.JSON(http.StatusUnauthorized, gin.H{
				"isOk": false,
			})
			model.SugarLogger.Infof("%v", err.Error())
			return
		}
	}
}

func Verify(c *gin.Context) {
	// refresh-token 쿠키를 요청한다.
	bearer, err := c.Request.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk" : false,
		})
		model.SugarLogger.Errorf("cookie error : %v", err.Error())
		c.Abort() // 이후의 핸들러 호출하지 않도록.
		return
	}

	isValid, _ := model.RefreshTokenValid(bearer.Value)
	if isValid == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"isOk" : false,
		})
		c.Abort() // 이후의 핸들러 호출하지 않도록.
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
		//model.SugarLogger.Info("token refresh")
		return
	}
}
