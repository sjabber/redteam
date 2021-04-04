package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	UserNo        int    `json:"user_no"`
	Email         string `json:"email"`
	PasswordHash  string `json:"-"`
	Password      string `json:"password"`
	PasswordCheck string `json:"password_check"` //회원가입에서 사용된다.
	Name          string `json:"name"`
}

// JWT 토큰을 반환해 주는 메서드
func (u *User) GetAuthToken() (string, string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = u.Email
	claims["user_name"] = u.Name
	claims["user_no"] = u.UserNo
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
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

// 패스워드가 확실한지 체크하고 사용자가 로그인상태인지 확인하는 메서드
func (u *User) IsAuthenticated(conn *sql.DB) (error, int, int) {

	loginCount := 0
	u.Email = strings.ToLower(u.Email)

	if u.Email == "" || u.Password == "" {
		SugarLogger.Info("login fail, Please enter your account information.")
		return fmt.Errorf("login fail"), 400, loginCount
	}

	row := conn.QueryRowContext(context.Background(),
		"SELECT user_no, user_name, user_pw_hash FROM user_info WHERE user_email = $1", u.Email)
	err := row.Scan(&u.UserNo, &u.Name, &u.PasswordHash)
	if err != nil || u.UserNo == 0 || u.Name == "" {
		SugarLogger.Info(err.Error())
		return fmt.Errorf("login fail"), 403, loginCount
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		conn.QueryRow(`select login_count from user_info where user_no = $1`, u.UserNo).Scan(&loginCount)
		if loginCount >= 5 {
			SugarLogger.Info("login fail, Exceeded number of login attempts.")
			return fmt.Errorf("login fail"), 408, loginCount
		}
		conn.QueryRow(`update user_info
							 set login_count = login_count + 1
							 where user_no = $1
							 and login_count < 5
							 returning login_count`, u.UserNo)
		SugarLogger.Info("login fail, The password is incorrect.")
		return fmt.Errorf("login fail"), 409, loginCount
	} else {
		// 5회 안에 로그인 성공 시 login count 초기화
		conn.Exec(`update user_info
						 set login_count = 0
						 where user_no = $1
						 and login_count < 5`, u.UserNo)
	}

	defer conn.Close()
	return nil, 200, loginCount
}

// 토큰이 유효한지 검사하는 메서드
// 미들웨어의 TokenAuthMiddleWare() 에서 사용된다.
func IsTokenValid(tokenString string) (bool, User) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// HMAC 을 사용하는 이유
		// REST API(표현상태 전송 API)가 요청을 받았을 때, 이 요청이 신뢰할 수 있는 호출인지
		// 확인하는 기법으로 요청이 부적절한지 정상적인지 확인할 수 있다.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v",	token.Header["alg"])
		}
		return tokenSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			SugarLogger.Errorf("token invalid : %v", err)
			return false, User{}
		}

		//SugarLogger.Infof("token expired : %v", err) // Access Token 만료
		return false, User{}
	}

	// 위에서 ok가 true, token 의 valid 값이 true 면 여기서 true 를 반환하며 검증을 완료
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// MapClaims 는 JSON 디코딩을 위해 map[string]interface{}를 사용함.
		// 디폴트 claims 타입임. map 은 java 의 해시같은 개념.
		// + 대시보드 API 로 반환될 구조체를 정의함.
		user := User{
			Email:  claims["user_email"].(string),
			Name:   claims["user_name"].(string),
			UserNo: int(claims["user_no"].(float64)),
		}
		return true, user
	} else {
		// 예상 밖 토큰인증에 문제가 발생한 경우 로그
		SugarLogger.Errorf("unexpected error, ok : %v, token Valid : %v", ok, token.Valid)
		return false, User{}
	}
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetID() string {
	return u.Email
}

func (u *User) DelUser(conn *sql.DB) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = conn.Exec(`
delete
from count_info
where target_no in (select target_no from target_info where user_no = $1)`, u.UserNo)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
delete
from target_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
delete
from project_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
delete
from smtp_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
delete
from template_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`
delete
from tag_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}

	_, err = conn.Exec(`
delete
from user_info
where user_no = $1`, u.UserNo)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
