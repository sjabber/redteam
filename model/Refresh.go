package model

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Refresh 토큰을 검사하는데 사용할 메서드
// Login.go 의 IsTokenValid 와 반환값 빼고는 동일함.
func RefreshTokenValid(tokenString string) (bool, User) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v",	token.Header["alg"])
		}
		return tokenRefresh, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			SugarLogger.Errorf("token invalid : %v", err)
			return false, User{}
		}

		//SugarLogger.Infof("token expired : %v", err) // Refresh Token 만료
		return false, User{}
	}

	// 위에서 ok가 true, token 의 valid 값이 true 면 여기서 true 를 반환하며 검증을 완료
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// + 대시보드 API 로 반환될 구조체를 정의함.
		user := User{
			Email: claims["user_email"].(string),
			Name: claims["user_name"].(string),
			UserNo: int(claims["user_no"].(float64)),
		}
		return true, user
	} else {
		// 예상 밖 토큰 토큰인증에 문제가 발생한 경우 로그
		SugarLogger.Errorf("unexpected error, ok : %v, token Valid : %v", ok, token.Valid)
		return false, User{}
	}
}

// 토큰을 재발행 해주는 메서드
func (u *User) GetNewToken() (string, string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = u.Email
	claims["user_name"] = u.Name
	claims["user_no"] = u.UserNo
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)

	rfClaims := jwt.MapClaims{}
	rfClaims["user_email"] = u.Email
	rfClaims["user_name"] = u.Name
	rfClaims["user_no"] = u.UserNo
	rfClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rfToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &rfClaims)
	refreshToken, err := rfToken.SignedString(tokenRefresh)


	return authToken, refreshToken, err
}


