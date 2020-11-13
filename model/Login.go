package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UserNo          int    `json:"user_no"`
	Email           string `json:"email"`
	PasswordHash    string `json:"-"`
	Password        string `json:"password"`
	PasswordCheck 	string `json:"password_check"` //회원가입에서 사용된다.
	Name            string `json:"name"`
}


// Go는 메서드, 변수 이름이 대문자로 시작 -> public, 소문자로 시작 -> private
// User를 클래스로 보면 login()이 User의 메서드로 작용하는 느낌이다.
// 디비와 커넥션을 해서 로그인을 확인하는 로직 샘플
//	db, err := ConnectDB()
//	if err != nil {
//		return num, fmt.Errorf("db connection error")
//	}
//	defer db.Close()
//
//	query := "select user_no, user_name from user_info where user_email=$1"
//	err = db.QueryRow(query, u.Email).Scan(&u.UserNo, &u.Name) //쿼리의 내용을 err 에 저장


// JWT 토큰을 반환해 주는 메서드
func (u *User) GetAuthToken() (string, string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = u.Email
	claims["user_name"] = u.Name
	claims["user_no"] = u.UserNo
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
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
func (u *User) IsAuthenticated(conn *sql.DB) (error, int) {
	num := 200 // 기본적으로 200 (StatusOk의 값을 넣어놓는다.)

	if u.Email == "" || u.Password == "" {
		num = 400 // 아이디나 패스워드를 입력하지 않은경우.
		return fmt.Errorf("Please enter your account information. "), num
	}

	row := conn.QueryRowContext(context.Background(), "SELECT user_no, user_name, user_pw_hash FROM user_info WHERE user_email = $1", u.Email)
	err := row.Scan(&u.UserNo, &u.Name, &u.PasswordHash)
	if err != nil || u.UserNo == 0 || u.Name == "" {
		num = 403 // 일치하는 계정이 존재하지 않을 경우.
		return fmt.Errorf("this account does not exist. "), num
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		num = 401 // 패스워드가 일치하지 않을경우.
		return fmt.Errorf("The password is incorrect. "), num
	}

	return nil, num
}

// 토큰이 유효한지 검사하는 메서드
// 미들웨어의 TokenAuthMiddleWare() 에서 사용된다.
func IsTokenValid(tokenString string) (bool, User) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("unexpected signing method : %v",
				token.Header["alg"])
			// HMAC 을 사용하는 이유
			// REST API(표현상태 전송 API)가 요청을 받았을 때, 이 요청이 신뢰할 수 있는 호출인지
			// 확인하는 기법으로 요청이 부적절한지 정상적인지 확인할 수 있다.
		}
		return tokenSecret, nil
	})

	if err != nil {
		fmt.Errorf("Error content : %v \n", err)
		return false, User{}
	}

	//위에서 ok가 true, token 의 valid 값이 true 면 여기서 true 를 반환하며 검증을 완료
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// MapClaims 는 JSON 디코딩을 위해 map[string]interface{}를 사용함.
		// 디폴트 claims 타입임. map 은 java 의 해시같은 개념.
		// fmt.Println(claims)
		user := User{
			Email: claims["user_email"].(string),
			Name: claims["user_name"].(string),
			UserNo: int(claims["user_no"].(float64)),
		}
		//user := claims["user_id"]
		//userNo := claims["user_no"]
		return true, user
	} else { // 토큰인증에 문제가 발생한 경우 오류메시지
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, User{}
	}
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetID() string {
	return u.Email
}
