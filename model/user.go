package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"time"

	"log"
)

type User struct {
	UserNo          int    `json:"user_no"`
	Email           string `json:"email"`
	PasswordHash    string `json:"-"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Name            string `json:"name"`
}

func (u *User) Login() error {
	// 자바에서 클래스 선언할때 변수들 public, private 그 밑에 메서드
	// 한 클래스안에 들어가야 메서드
	// User를 클래스로 보면 login()이 User의 메서드로 작용하는 느낌이다.
	// 디비와 커넥션을 해서 로그인을 확인하는 로직을 구현

	//user := models.User{}
	//err := l.ShouldBindJSON(&user)
	//if err != nil {
	//	l.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return err
	//}

	db, err := ConnectDb()
	if err != nil {
		return fmt.Errorf("db connection error")
	}
	defer db.Close()

	query := "select user_no, user_name from user_info where user_id=$1 and user_pw=$2"
	err = db.QueryRow(query, u.Email, u.Password).Scan(&u.UserNo, &u.Name) //쿼리의 내용을 err 에 저장

	if err == nil {
		log.Println("login true")

		return nil
	} else {
		log.Println("login false")
		return fmt.Errorf("login fail")
	}
}

// JWT 토큰을 반환해 주는 메서드
func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.Email
	claims["user_name"] = u.Name
	claims["user_no"] = u.UserNo
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

//패스워드가 확실한지 체크하고 사용자가 로그인상태인지 확인하는 메서드
func (u *User) IsAuthenticated(conn *sql.DB) error {
	row := conn.QueryRowContext(context.Background(), "SELECT user_pw_hash FROM user_info WHERE user_id = $1", u.Email)
	err := row.Scan(&u.PasswordHash)

	if err == pgx.ErrNoRows {
		fmt.Println("해당 계정이 존재하지 않습니다.")
		return fmt.Errorf("로그인 자격증명이 올바르지 않습니다. ")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("로그인 자격증명이 올바르지 않습니다. ")
	}

	return nil
}

func IsTokenValid(tokenString string) (bool, User) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("토큰 서명이 유요하지 않습니다. : %v",
				token.Header["alg"])
			// HMAC 을 사용하는 이유
			// REST API(표현상태 전송 API)가 요청을 받았을 때,
			// 이 요청이 신뢰할 수 있는 호출인지 확인하는 기법으로
			// 요청이 부적절한지 정상적인지 확인할 수 있다.
		}
		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("에러내용 %v \n", err)
		return false, User{}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// MapClaims 는 JSON 디코딩을 위해 map[string]interface{}를 사용함.
		// 디폴트 claims 타입임. map 은 java 의 해시같은 개념.
		// fmt.Println(claims)
		user := User{
			Email: claims["user_id"].(string),
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
