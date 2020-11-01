package model

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strings"
)

type Target struct {
	TargetNo         int    `json:"tg_no"`
	TargetName       string `json:"tg_name"`
	TargetEmail      string `json:"tg_email"`
	TargetPhone      string `json:"tg_phone"`
	TargetOrganize   string `json:"tg_organize"` //소속
	TargetPosition   string `json:"tg_position"` //직급
	TargetClassify   string `json:"tg_tag"`      //태그
	TargetCreateTime string `json:"created_t"`
}

func (t *Target) CreateTarget(conn *sql.DB) error {

	t.TargetName = strings.Trim(t.TargetName, " ")
	t.TargetEmail = strings.Trim(t.TargetEmail, " ")
	t.TargetPhone = strings.Trim(t.TargetPhone, " ")
	t.TargetOrganize = strings.Trim(t.TargetOrganize, " ")
	t.TargetPosition = strings.Trim(t.TargetPosition, " ")

	if len(t.TargetName) < 1 {
		return fmt.Errorf("이름이 비어있습니다. ")
	} else if len(t.TargetEmail) < 1 {
		return fmt.Errorf(" 이메일이 비어있습니다. ")
	} else if len(t.TargetPhone) < 1 {
		return fmt.Errorf(" 이메일이 비어있습니다. ")
	} else if len(t.TargetOrganize) < 1 {
		return fmt.Errorf(" 소속이 비어있습니다. ")
	} else if len(t.TargetPosition) < 1 {
		return fmt.Errorf(" 직급이 비어있습니다. ")
	}

	query := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING target_no"

	row := conn.QueryRow(query, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition)

	err := row.Scan(&t.TargetNo)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("대상 입력 중 오류가 발생하였습니다. ")
	}

	return nil
}

func ReadTarget() ([]Target, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 오류")
	}

	query := "SELECT target_no, target_name, target_email, target_phone, target_organize, target_position, modified_time from target_info"

	rows, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("대상자들을 불러오는데 오류가 발생하였습니다. ")
	}

	var targets []Target
	for rows.Next() {
		tg := Target{}
		err = rows.Scan(&tg.TargetNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetCreateTime)

		if err != nil {
			fmt.Printf("훈련대상 스캐닝 오류 : %v", err)
			continue
		}
		targets = append(targets, tg)

	}

	return targets, nil
}

func (t *Target) DeleteTarget(conn *sql.DB) error {
	str := string(t.TargetNo) // int -> string 형변환
	if str == "" {
		return fmt.Errorf("삭제할 대상의 번호를 입력해 주세요. ")
	}
	_, err := conn.Exec("DELETE FROM target_info WHERE target_no = $1", t.TargetNo)
	if err != nil {
		//fmt.Printf("Error deleting target: (%v)", err)
		return fmt.Errorf("Error deleting target ")
	}

	return nil
}

func Download() error {
	f := excelize.NewFile()
	style, err := f.NewStyle(`{"font": "bold"}`)
	if err != nil {
		fmt.Println(err)
	}
	categories := map[string]string{"A1": "이름", "B1": "이메일", "C1": "연락처", "D1": "소속",
		"E1": "직급", "F1": "태그"}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
		f.SetCellStyle("Sheet1", k, v, style)
	}
	f.SaveAs("./훈련대상.xlsx")
	return nil
}
