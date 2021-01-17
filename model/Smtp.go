package model

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
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

func (sm *Smtpinfo) IdPwCheck(conn *sql.DB) error {
	sm.SmtpId = strings.ToLower(sm.SmtpId)

	if sm.SmtpHost == "" || sm.SmtpPort == "" || sm.SmtpTimeout == "" || sm.SmtpId == "" || sm.SmtpPw == ""  {
		return fmt.Errorf("Please enter your Smtp information. ")
	}

	row := conn.QueryRow(`SELECT smtp_pw FROM smtp_info WHERE smtp_id = $1`, sm.SmtpId)
	err := row.Scan(&sm.SmtpPwHashCheck)
	if err != nil {
		return fmt.Errorf("this account does not exist. ")
	}

	err = bcrypt.CompareHashAndPassword([]byte(sm.SmtpPwHashCheck), []byte(sm.SmtpPw))
	if err != nil {
		return fmt.Errorf("This Password is incorrect. ")
	}

	return nil
}


func (sm *Smtpinfo) SmtpConnectionCheck(conn *sql.DB, num int) error {

	row := conn.QueryRow(`SELECT smtp_host, smtp_port, smtp_id, smtp_pw
 								FROM smtp_info
 								WHERE user_no = $1`, num)
	err := row.Scan(&sm.SmtpHost, &sm.SmtpPort, &sm.SmtpId, &sm.SmtpPw)
	if err != nil {
		return fmt.Errorf("this account does not exist. ")
	}

	// 비밀번호 일치여부 검사
	//err = bcrypt.CompareHashAndPassword([]byte(sm.SmtpPwHashCheck), []byte(sm.SmtpPw))
	//if err != nil {
	//	return fmt.Errorf("This Password is incorrect. ")
	//}

	// string -> int, smtp 연결을 테스트한다.
	port, _ := strconv.Atoi(sm.SmtpPort)
	d := gomail.NewDialer(sm.SmtpHost, port, sm.SmtpId, sm.SmtpPw)
	_, err = d.Dial()
	if err != nil {
		return fmt.Errorf("Smtp connecting failed. : %v ", err)
	}
	//sm.SendMail2()
	return nil
}

//func (sm *Smtpinfo) SendMail2() error {
//
//	port, _ := strconv.Atoi(sm.SmtpPort) //string -> int
//
//	d := gomail.NewDialer(sm.SmtpHost, port, sm.SmtpId, sm.SmtpPw)
//
//	s, err := d.Dial()
//
//	if err != nil {
//		return err
//	}
//
//	m := gomail.NewMessage()
//	// for _, r := range list {
//	m.SetHeader("From", sm.SmtpId) //보내는 사람
//	m.SetAddressHeader("To", sm.SmtpId, "받는분 이름") //받는사람
//	m.SetHeader("Subject", "smtp test") //메일 제목
//	m.SetBody("text/html", fmt.Sprintf("Hello %s!", " 김태호")) //내용
//
//	if err := gomail.Send(s, m); err != nil {
//		return fmt.Errorf(
//			"Could not send email to %q: %v ", "보내는 계정주소", err)
//	}
//	m.Reset()
//	// }
//	return nil
//}