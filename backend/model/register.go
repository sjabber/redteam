package model

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"time"
)

//type user User

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

func (u *User) Register() error {

	db, err := ConnectDb()
	if err != nil {
		return fmt.Errorf("db connection error")
	}
	defer db.Close()

	query := "SELECT user_pw FROM user_info WHERE user_id=$1"
	err = db.QueryRow(query, u.Password).Scan(&u.Password, &u.PasswordConfirm)


	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf(" 비밀번호는 적어도 4 글자 이상이어야 합니다. ")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf(" 비밀번호가 일치하지 않습니다. ")
	}

	//추후 Email 형식검사에 대한 부분을 보강할 예정입니다.
	if len(u.Email) < 4 {
		return fmt.Errorf(" 이메일은 적어도 4 글자 이상이어야 합니다. ")
	}

	//사용자가 보낸 이메일을 모두 소문자로 변경한다.
	u.Email = strings.ToLower(u.Email)

	query = "SELECT user_id FROM user_info WHERE user_id = $1"
	row := db.QueryRow(query, u.Email)

	userLookup := User{}
	err = row.Scan(&userLookup)
	if err != sql.ErrNoRows {
		fmt.Println("found user : " + userLookup.Email)
		return fmt.Errorf(" 이미 존재하는 ID입니다. / ID already exist ")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf(" 계정을 생성하는데 오류가 발생하였습니다. ")
	}
	u.PasswordConfirm = string(pwdHash)

	_, err = db.Exec("INSERT INTO user_info " +
		"(user_id, user_name, user_pw, created_time) " +
		"VALUES($1, $2, $3, $4)", u.Email, u.Name, u.Password, time.Now())

	return err
}

// JW 토큰은 일정시간이지나면 만료된다. 사용될 토큰을 반환하는 메서드
func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
 }