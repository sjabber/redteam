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
	SmtpNo       int    `json:"smtp_no"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     string `json:"smtp_port"`
	SmtpProtocol string `json:"smtp_protocol"`
	SmtpTls      string `json:"smtp_tls"`
	SmtpTimeout  string `json:"smtp_timeout"`
	SmtpId       string `json:"smtp_id"`
	SmtpPw       string `json:"smtp_pw"`
	SmtpIdentity string `json:"smtp_identity"`
}

func (sm *Smtpinfo) IdPwCheck(conn *sql.DB) error {
	sm.SmtpId = strings.ToLower(sm.SmtpId)

	if sm.SmtpHost == "" || sm.SmtpPort == "" || sm.SmtpTimeout == "" || sm.SmtpId == "" || sm.SmtpPw == ""  {
		return fmt.Errorf("Please enter your Smtp information. ")
	}

	row := conn.QueryRow(`SELECT smtp_pw FROM smtp_info WHERE smtp_id = $1`, sm.SmtpId)
	err := row.Scan(&sm.SmtpIdentity)
	if err != nil {
		return fmt.Errorf("this account does not exist. ")
	}

	err = bcrypt.CompareHashAndPassword([]byte(sm.SmtpIdentity), []byte(sm.SmtpPw))
	if err != nil {
		return fmt.Errorf("This Password is incorrect. ")
	}

	return nil
}


func (sm *Smtpinfo) SmtpConnectionCheck(num int) error {
	port, _ := strconv.Atoi(sm.SmtpPort)
	d := gomail.NewDialer(sm.SmtpHost, port, sm.SmtpId, sm.SmtpPw)

	_, err := d.Dial()

	if err != nil {
		return err
	}
	return nil
}

func (sm *Smtpinfo) SendMail() error {

	port, _ := strconv.Atoi(sm.SmtpPort) //string -> int

	d := gomail.NewDialer(sm.SmtpHost, port, sm.SmtpId, sm.SmtpPw)

	s, err := d.Dial()

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	// for _, r := range list {
	m.SetHeader("From", "보내는 계정주소")
	m.SetAddressHeader("To", "받는계정주소", "이름")
	m.SetHeader("Subject", "메일제목")
	m.SetBody("text/html", fmt.Sprintf("Hello %s!", "이름"))

	if err := gomail.Send(s, m); err != nil {
		return fmt.Errorf(
			"Could not send email to %q: %v ", "보내는 계정주소", err)
	}
	m.Reset()
	// }
	return nil
}