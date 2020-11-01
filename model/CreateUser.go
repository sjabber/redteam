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

func (u *User) CreateUsers() (int, error) {

	num := 7
	db, err := ConnectDB()
	if err != nil {
		return num, fmt.Errorf("db connection error")
	}
	defer db.Close()

	// 셋중 아무것도 입력하지 않았을 경우
	if len(u.Password) < 1 || len(u.PasswordCheck) < 1 || len(u.Email) < 1 || len(u.Name) < 1{
		num = 0
		return num, fmt.Errorf(" 정보를 입력해 주세요. ")
	}

	// 이메일 형식을 검사하는 정규식
	var validEmail, _ = regexp.MatchString(
		"^[_a-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", u.Email)
	//^[_a-z0-9+-.]+@[a-z0-9-]+.[a-z0-9-]*.[a-z]{2,4}

	// 이메일형식검사
	if validEmail != true {
		num = 2 // 402에러, 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf(" 이메일 형식이 올바르지 않습니다. ")
	}

	// 비밀번호 길이 16자 미만
	if len(u.Password) > 16 || len(u.PasswordCheck) > 16 {
		num = 2 // 402에러, 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf(" 비밀번호는 16 글자 이하여야 합니다. ")
	}

	// 비밀번호 길이 8자 이상
	if len(u.Password) < 7 || len(u.PasswordCheck) < 7 {
		num = 2 // 402에러, 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf(" 비밀번호는 적어도 8글자 이상이어야 합니다. ")
	}

	// 비밀번호 형식검사 검증
	if CheckPassword(u.Password) != nil {
		num = 2 // 402에러, 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, CheckPassword(u.Password)
	}

	// 비밀번호와 비밀번호 확인이 일치하지 않을경우
	if u.Password != u.PasswordCheck {
		num = 1 // 401에러, 비밀번호가 일치하지 않습니다.
		return num, fmt.Errorf(" 비밀번호가 일치하지 않습니다. ")
	}

		//사용자가 보낸 이메일을 모두 소문자로 변경한다.
	u.Email = strings.ToLower(u.Email)

	query := "SELECT user_email FROM user_info WHERE user_email = $1"
	row := db.QueryRow(query, u.Email)

	userLookup := User{}
	err = row.Scan(&userLookup)
	if err != sql.ErrNoRows {
		num = 3 // 403에러, 이미 존재하는 이메일입니다.
		fmt.Println("found user : " + userLookup.Email)
		return num, fmt.Errorf(" 이미 존재하는 이메일 입니다. / ID already exist ")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		num = 5 // 405에러, 계정생성 오류발생
		return num, fmt.Errorf(" 계정을 생성하는데 오류가 발생하였습니다. ")
	}
	u.PasswordHash = string(pwdHash)

	_, err = db.Exec("INSERT INTO user_info " +
		"(user_name, user_email, user_pw, user_pw_hash, created_time) " +
		"VALUES($1, $2, $3, $4, $5)", u.Name, u.Email, u.Password, u.PasswordHash, time.Now())

	return num, err
}

 // 비밀번호 형식을 검사하는 메서드
func CheckPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf(" 비밀번호는 적어도 8글자 이상이어야 합니다. ")
	}
	num := `[0-9]{1}`
	a_z := `[a-zA-Z]{1}`
	//A_Z := `[A-Z]{1,}`
	symbol := `[!@#~$%^&*()+|_]{1,}`
	if b, err := regexp.MatchString(num, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 숫자를 하나 이상 포함해야 합니다. ")
	}
	if b, err := regexp.MatchString(a_z, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 영문을 하나 이상 포함해야 합니다. ")
	}
	//if b, err := regexp.MatchString(A_Z, pw); !b || err != nil {
	//	return fmt.Errorf("비밀번호는 적어도 대문자를 하나 이상 포함해야 합니다. ")
	//}
	if b, err := regexp.MatchString(symbol, pw); !b || err != nil {
		return fmt.Errorf("비밀번호는 적어도 특수문자를 하나 이상 포함해야 합니다. ")
	}
	return nil
}