package model

import (
	"crypto/aes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
)

type Smtpinfo struct {
	SmtpNo          int    `json:"smtp_no"`
	SmtpHost        string `json:"smtp_host"`
	SmtpPort        string `json:"smtp_port"`
	SmtpProtocol    string `json:"smtp_protocol"`
	SmtpTls         string `json:"smtp_tls"`
	SmtpTimeout     string `json:"smtp_timeout"`
	SmtpId          string `json:"smtp_id"`
	SmtpPw          string `json:"smtp_pw"`
	SmtpPwHashCheck string `json:"smtp_identity"` // smtp 비밀번호 해시값 체크
}

//func (sm *Smtpinfo) IdPwCheck(conn *sql.DB) error {
//	sm.SmtpId = strings.ToLower(sm.SmtpId)
//
//	if sm.SmtpHost == "" || sm.SmtpPort == "" || sm.SmtpTimeout == "" || sm.SmtpId == "" || sm.SmtpPw == ""  {
//		return fmt.Errorf("Please enter your Smtp information. ")
//	}
//
//	row := conn.QueryRow(`SELECT smtp_pw FROM smtp_info WHERE smtp_id = $1`, sm.SmtpId)
//	err := row.Scan(&sm.SmtpPwHashCheck)
//	if err != nil {
//		return fmt.Errorf("this account does not exist. ")
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(sm.SmtpPwHashCheck), []byte(sm.SmtpPw))
//	if err != nil {
//		return fmt.Errorf("This Password is incorrect. ")
//	}
//
//	return nil
//}


func (sm *Smtpinfo) SmtpConnectionCheck(conn *sql.DB, num int) error {

	row := conn.QueryRow(`SELECT smtp_host, smtp_port, smtp_id, smtp_pw
 								FROM smtp_info
 								WHERE user_no = $1`, num)
	err := row.Scan(&sm.SmtpHost, &sm.SmtpPort, &sm.SmtpId, &sm.SmtpPw)
	if err != nil {
		return fmt.Errorf("this account does not exist. ")
	}

	// smtp 패스워드 복호화 작업 수행
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	password, _ := base64.StdEncoding.DecodeString(sm.SmtpPw)
	password = Decrypt(block, password)
	sm.SmtpPw = string(password)

	// string -> int, smtp 연결을 테스트한다.
	port, _ := strconv.Atoi(sm.SmtpPort)
	d := gomail.NewDialer(sm.SmtpHost, port, sm.SmtpId, sm.SmtpPw)
	_, err = d.Dial()
	if err != nil {
		return fmt.Errorf("Smtp connecting failed. : %v ", err)
	}
	return nil
}
