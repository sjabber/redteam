package model

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	_ "regexp"
	"strings"
	"time"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

func (u *User) CreateUser() error {

	db, err := ConnectDb()
	if err != nil {
		return fmt.Errorf("db connection error")
	}
	defer db.Close()

	//query := "SELECT user_pw FROM user_info WHERE user_id=$1"
	//err = db.QueryRow(query, u.Email).Scan(&u.Password, &u.PasswordConfirm)

	//이메일 형식을 검사하는 정규식
	var validEmail, _ = regexp.MatchString(
		"^[_a-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", u.Email)

	if len(u.Password) > 16 || len(u.PasswordConfirm) > 16 {
		return fmt.Errorf(" 비밀번호는 16 글자 이하여야 합니다. ")
	}

	if len(u.Password) < 8 || len(u.PasswordConfirm) < 8 {
		return fmt.Errorf(" 비밀번호는 적어도 8글자 이상이어야 합니다. ")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf(" 비밀번호가 일치하지 않습니다. ")
	}

	if CheckPassword(u.Password) != nil {
		return CheckPassword(u.Password)
		return fmt.Errorf(" 비밀번호 형식이 올바르지 않습니다. ")
	}

	//이메일형식검사 todo : 10/09 이메일 형식 오류 검사
	if validEmail != true {
		return fmt.Errorf(" 이메일 형식이 올바르지 않습니다. ")
	}

	//사용자가 보낸 이메일을 모두 소문자로 변경한다.
	u.Email = strings.ToLower(u.Email)

	query := "SELECT user_id FROM user_info WHERE user_id = $1"
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
	u.PasswordHash = string(pwdHash)

	_, err = db.Exec("INSERT INTO user_info " +
		"(user_id, user_name, user_pw, user_pw_hash, created_time) " +
		"VALUES($1, $2, $3, $4, $5)", u.Email, u.Name, u.Password, u.PasswordHash, time.Now())

	return err
}

//// JWT 토큰을 반환해 주는 메서드
//func (u *User) GetAuthToken() (string, error) {
//	claims := jwt.MapClaims{}
//	claims["authorized"] = true
//	claims["user_id"] = u.Email
//	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
//	authToken, err := token.SignedString(tokenSecret)
//	return authToken, err
//}

 // 비밀번호를 검사하는 메서드
 // todo : 10월 9일 -> 대소문자 구부없이 1개 이상 들어가도록 변경하겠습니다.
func CheckPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf(" 비밀번호는 적어도 8글자 이상이어야 합니다. ")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1,}`
	symbol := `[!@#~$%^&*()+|_]{1,}`
	if b, err := regexp.MatchString(num, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 숫자를 하나 이상 포함해야 합니다. ")
	}
	if b, err := regexp.MatchString(a_z, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 소문자를 하나 이상 포함해야 합니다. ")
	}
	if b, err := regexp.MatchString(A_Z, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 대문자를 하나 이상 포함해야 합니다. ")
	}
	if b, err := regexp.MatchString(symbol, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 특수문자를 하나 이상 포함해야 합니다. ")
	}
	return nil
}