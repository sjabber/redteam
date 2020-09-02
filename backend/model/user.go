package model

import (
	"fmt"
	"log"
)

type User struct {
	UserNo   int
	Email    string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Password string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Name     string `json:"name"`
}

func (u *User) Login() error {
	// 자바에서 클래스 선언할때 변수들 public, private 그 밑에 메서드
	// 한 클래스안에 들어가야 메서드
	// User를 클래스로 보면 login()이 User의 메서드로 작용하는 느낌이다.

	// 디비와 커넥션을 해서 로그인을 확인하는 로직을 구현

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

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetID() string {
	return u.Email
}



