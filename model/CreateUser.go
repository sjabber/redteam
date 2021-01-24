package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	_ "regexp"
	"strings"
	"time"
)

var (
	tokenSecret  = []byte(os.Getenv("TOKEN_SECRET"))
	tokenRefresh = []byte(os.Getenv("TOKEN_REFRESH"))
	key = "qlwkndlqiwndlian"
)

func (u *User) CreateUsers() (int, error) {

	num := 200
	db, err := ConnectDB()
	if err != nil {
		return num, fmt.Errorf("db connection error")
	}
	defer db.Close()

	// 입력하지 않은 정보가 존재하는지 검사
	// 400에러 : 필수요청 변수가 없는경우
	if len(u.Password) < 1 || len(u.PasswordCheck) < 1 || len(u.Email) < 1 || len(u.Name) < 1 {
		num = 400 // 정보를 입력해 주세요.
		return num, fmt.Errorf(" Please enter the information. ")
	}

	// 이메일 형식을 검사하는 정규식
	var validEmail, _ = regexp.MatchString(
		"^[_A-Za-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", u.Email)

	// 402에러 : 회원가입 계정 이메일이나 비밀번호 형식이 잘못된 경우
	if validEmail != true {
		num = 402 // 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf("Email format is incorrect. ")
	}

	// 비밀번호 길이 16자 미만
	if len(u.Password) > 17 || len(u.PasswordCheck) > 17 {
		num = 402 // 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf("Password must be 16 characters or less. ")
	}

	// 비밀번호 길이 8자 이상
	if len(u.Password) < 7 || len(u.PasswordCheck) < 7 {
		num = 402 // 비밀번호나 이메일 형식이 올바르지 않습니다.
		return num, fmt.Errorf("Password must be at least 8 characters long. ")
	}

	// 비밀번호 형식검사 검증
	if CheckPassword(u.Password) != nil {
		num = 402 // 비밀번호 검증시 상황별로 에러메시지 출력
		return num, CheckPassword(u.Password)
	}

	// 비밀번호와 비밀번호 확인이 일치하지 않는지 검사
	// 401에러 : 비밀번호와 비밀번호 확인이 일치하지 않을 경우
	if u.Password != u.PasswordCheck {
		num = 401 // 비밀번호가 일치하지 않습니다.
		return num, fmt.Errorf("Passwords do not match. ")
	}

	//사용자가 보낸 이메일을 모두 소문자로 변경한다.
	u.Email = strings.ToLower(u.Email)

	query := "SELECT user_email FROM user_info WHERE user_email = $1"
	row := db.QueryRow(query, u.Email)

	// 403에러 : 회원가입시 이미 존재하는 계정이 있는경우
	userLookup := User{}
	err = row.Scan(&userLookup)
	if err != sql.ErrNoRows {
		num = 403  // 이미 존재하는 이메일입니다.
		fmt.Println("found user : " + userLookup.Email)
		return num, fmt.Errorf("This account already exists. ")
	}

	// 패스워드 해싱
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		num = 500 // 405에러, 계정생성 오류발생
		return num, fmt.Errorf("There was an error creating an account. ")
	}
	u.PasswordHash = string(pwdHash)

	// 회원 정보를 등록한다.
	_, err = db.Exec(`INSERT INTO user_info (user_name, user_email, user_pw_hash, created_time)
        VALUES ($1, $2, $3, $4)`, u.Name, u.Email, u.PasswordHash, time.Now())
	if err != nil {
		return 500, fmt.Errorf("A database error has occurred. (1)")
	}


	// Note smtp_info 테이블에 삽입될 비밀번호를 암호화한다.
	block, err := aes.NewCipher([]byte(key)) // key 값을 블록에 집어넣는다.
	if err != nil {
		num = 500
		fmt.Println(err)
		return num, err
	}

	// AES 암호화 알고리즘 적용
	pwdAES := base64.StdEncoding.EncodeToString(Encrypt(block, []byte(u.Password)))// byte[]타입

	// Smtp 디폴트 정보 삽입
	query = "SELECT user_no FROM user_info WHERE user_email = $1"
	row = db.QueryRow(query, u.Email)

	err = row.Scan(&u.UserNo)
	_, err = db.Exec(`INSERT INTO smtp_info (user_no, smtp_id, smtp_pw)
        VALUES ($1, $2, $3)`, u.UserNo, u.Email, pwdAES)
	if err != nil {
		return 500, fmt.Errorf("A database error has occurred. (2)")
	}

	return num, err
}

// 암호화
func Encrypt(b cipher.Block, plaintext []byte) []byte {
	// plaintext(원본)의 크기가 블록크기의 배수가 아니면 채워준다(padding)
	mod := len(plaintext) % aes.BlockSize
	// 남는 공간에 패딩값을 삽입한다. (남는공간이 3이면 3을 세번 넣는다.)
	padding := bytes.Repeat([]byte{byte(aes.BlockSize - mod)}, aes.BlockSize-mod)
	plaintext = append(plaintext, padding...) // ...은 padding 슬라이스의 모든 컬렉션을 표현하는것

	ciphertext := make([]byte, len(plaintext)) // 16바이트의 배수로 맞춘 byte 배열공간을 만든다.

	// CBC(Cipher Block Chaining) 방식을 사용하여 암호화
	// 첫번째 블록은 직전 암호문이 없으므로 초기화 백터(iv)를 입력받는다.
	iv := []byte("0987654321654321") // 16byte 의 크기를 갖는다. Note + iv값 환경변수 & 암호화
	mode := cipher.NewCBCEncrypter(b, iv) // 암호화 블록과 초기화 벡터를 넣어서 암호화 블록 mode 인스턴스를 생성한다.
	mode.CryptBlocks(ciphertext, plaintext) // 암호화 블록 mode 인스턴스로 암호화

	return ciphertext
}

// 복호화
func Decrypt(b cipher.Block, ciphertext []byte) []byte {

	if len(ciphertext)%aes.BlockSize != 0 {
		// 암호화된 데이터의 길이기 블록크기의 배수가 아니면 작동이 안됨.
		fmt.Println("The length of decrypted data must be a multiple of the block size. ")
		return nil
	}
	iv := []byte("0987654321654321")

	plaintext := make([]byte, len(ciphertext)) // 평문 데이터를 저장할 공간을 생성한다.
	mode := cipher.NewCBCDecrypter(b, iv) // 암호화 블록과 초기화 벡터를 넣어 복호화된 블록모드의 인스턴스를 생성한다.

	mode.CryptBlocks(plaintext, ciphertext) //

	padding := plaintext[len(plaintext)-1] //가장 마지막 값(패딩값)을 가져온다.
	plaintext = plaintext[:len(plaintext)-int(padding)] // 패딩값을 빼준다.

	return plaintext
}

// 비밀번호 형식을 검사하는 메서드
func CheckPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf(" Password must be at least 8 characters long. ")
	}
	num := `[0-9]{1}`
	a_z := `[a-zA-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1,}`
	//A_Z := `[A-Z]{1,}`

	// 비밀번호에는 하나 이상의 숫자가 포함되어야합니다.
	if b, err := regexp.MatchString(num, pw); !b || err != nil {
		return fmt.Errorf("Passwords must contain at least one number. ")
	}
	// 비밀번호에는 하나 이상의 영문자가 포함되어야합니다.
	if b, err := regexp.MatchString(a_z, pw); !b || err != nil {
		return fmt.Errorf("Password must contain at least one English letter. ")
	}
	//  비밀번호는 적어도 특수문자를 하나 이상 포함해야 합니다.
	if b, err := regexp.MatchString(symbol, pw); !b || err != nil {
		return fmt.Errorf("Passwords must contain at least one special character. ")
	}

	//if b, err := regexp.MatchString(A_Z, pw); !b || err != nil {
	//	return fmt.Errorf("비밀번호는 적어도 대문자를 하나 이상 포함해야 합니다. ")
	//}
	return nil
}
