package model

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Refresh 토큰을 검사하는데 사용할 메서드
//  Login.go 의 IsTokenValid 와 반환값 빼고는 동일함.
func RefreshTokenValid(tokenString string) (bool, User) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("unexpected signing method : %v",
				token.Header["alg"])
		}
		return tokenRefresh, nil
	})

	if err != nil {
		fmt.Printf("에러내용 %v \n", err)
		return false, User{}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 대시보드 API 로 반환될 구조체를 정의함.
		user := User{
			Email: claims["user_email"].(string),
			Name: claims["user_name"].(string),
			UserNo: int(claims["user_no"].(float64)),
		}
		return true, user
	} else {
		fmt.Println("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, User{}
	}
}

// Access 토큰만 반환해 주는 메서드
func (u *User) GetAccessToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = u.Email
	claims["user_name"] = u.Name
	claims["user_no"] = u.UserNo
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)

	return authToken, err
}


