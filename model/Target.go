package model

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

type Target struct {
	TargetNo         int    `json:"tg_no"`
	TargetName       string `json:"tg_name"`
	TargetEmail      string `json:"tg_email"`
	TargetPhone      string `json:"tg_phone"`
	TargetOrganize   string `json:"tg_organize"` //소속
	TargetPosition   string `json:"tg_position"` //직급
	TargetTag		 string `json:"tg_tag"`      //태그
	TargetCreateTime string `json:"created_t"`
}

func (t *Target) CreateTarget(conn *sql.DB, num int) error {

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
		"VALUES ($1, $2, $3, $4, $5) WHERE user_no = $6" +
		"RETURNING target_no"

	row := conn.QueryRow(query, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition, num)

	err := row.Scan(&t.TargetNo)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("대상 입력 중 오류가 발생하였습니다. ")
	}

	return nil
}

func ReadTarget(num int) ([]Target, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 오류")
	}

	query := "SELECT target_no, target_name, target_email, target_phone, target_organize, target_position, target_tag, modified_time from target_info where user_no = $1"
	rows, err := db.Query(query, num)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("대상자들을 불러오는데 오류가 발생하였습니다. ")
	}

	var targets []Target
	for rows.Next() {
		tg := Target{}
		err = rows.Scan(&tg.TargetNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag, &tg.TargetCreateTime)

		if err != nil {
			fmt.Printf("훈련대상 스캐닝 오류 : %v", err)
			continue
		} else if len(tg.TargetTag) < 1{
			tg.TargetTag = " "
		}
		// 여기 조건 좀더 꼼꼼하게 만들어야함.

		targets = append(targets, tg)

	}

	return targets, nil
}

func (t *Target) DeleteTarget(conn *sql.DB, num int) error {
	str := string(t.TargetNo) // int -> string 형변환
	if str == "" {
		return fmt.Errorf("삭제할 대상의 번호를 입력해 주세요. ")
	}
	_, err := conn.Exec("DELETE FROM target_info WHERE target_no = $1 && user_no = $2", t.TargetNo, num)
	if err != nil {
		//fmt.Printf("Error deleting target: (%v)", err)
		return fmt.Errorf("Error deleting target ")
	}

	return nil
}

// 반복해서 읽고 값을 넣는것을 메서드로 구현하고 API는 이걸 그냥 사용하기만 하면됨.
func (t *Target) ImportTargets(conn *sql.DB, str string, num int) error {

	// todo 2 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 2는 전부 같은 경로로 수정, api_Target.go 파일의 todo 2 참고)
	f, err := excelize.OpenFile("C:/Users/Taeho/Downloads/"+ str)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	i  := 2

	for {
		str := strconv.Itoa(i)

		t.TargetName = f.GetCellValue("Sheet1", "A"+str)
		t.TargetEmail = f.GetCellValue("Sheet1", "B"+str)
		t.TargetPhone = f.GetCellValue("Sheet1", "C"+str)
		t.TargetOrganize = f.GetCellValue("Sheet1", "D"+str)
		t.TargetPosition = f.GetCellValue("Sheet1", "E"+str)

		if t.TargetName == "" {
			break
		}
		//	todo 4 : 추후 해당 목록에 적힌 글들의 값이 올바른 형식이 아닐경우 제외하도록 하는 코드도 삽입한다.
		query := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position) " +
			"VALUES ($1, $2, $3, $4, $5) WHERE user_no = &6" +
			"RETURNING target_no"

		row := conn.QueryRow(query, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition, num)
		err := row.Scan(&t.TargetNo)
		if err != nil {
			fmt.Println(err)
			break
		}

		i++
	}

	return nil
}

// DB에 저장된 값들을 읽어 엘셀파일에 일괄적으로 작성하여 저장한다.
func ExportTargets(num int) error {

	db, err := ConnectDB()
	query := "SELECT target_name, target_email, target_phone, target_organize, target_position, target_tag, modified_time from target_info WHERE user_no = $1"

	rows, err := db.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Database error. ")
	}


	i  := 2
	// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 1은 전부 같은 경로로 수정, api_Target.go 파일의 todo 1 참고)
	// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
	f, err := excelize.OpenFile("C:/Users/Taeho/go/src/redteam/Spreadsheet/sample.xlsx")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Open Spreadsheet Error. ")
	}

	index := f.NewSheet("Sheet1")

	for rows.Next() {
		tg := Target{}
		err = rows.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag , &tg.TargetCreateTime)
		if err != nil {
			fmt.Printf("Targets scanning error : %v", err)
			continue
		} else if len(tg.TargetTag) < 1{
			tg.TargetTag = " "
		}



		str := strconv.Itoa(i)
		f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
		f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
		f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
		f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
		f.SetCellValue("Sheet1", "E"+str, tg.TargetPhone)
		f.SetCellValue("Sheet1", "F"+str, tg.TargetTag)
		f.SetCellValue("Sheet1", "G"+str, tg.TargetCreateTime)



		i++
	}

	f.SetActiveSheet(index)

	str := strconv.Itoa(num)

	// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 3은 전부 같은 경로로 수정, api_Target.go 파일의 todo 3 참고)
	// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
	if err2 := f.SaveAs("C:/Users/Taeho/go/src/redteam/Spreadsheet/Registered_Targets" + str + ".xlsx"); err != nil {
		fmt.Println(err2)
		return fmt.Errorf("Registered Target downloading Error. ")
	}

	return nil
}
