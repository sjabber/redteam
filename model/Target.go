package model

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

// Target(훈련대상)을 관리하기 위한 json 구조체
type Target struct {
	TargetNo         int    `json:"tg_no"`
	TargetName       string `json:"tg_name"`
	TargetEmail      string `json:"tg_email"`
	TargetPhone      string `json:"tg_phone"`
	TargetOrganize   string `json:"tg_organize"` //소속
	TargetPosition   string `json:"tg_position"` //직급
	TargetTag		 string `json:"tg_tag"`      //태그
	TargetCreateTime string `json:"created_t"`
	TagNo			 int 	`json:"tag_no"`
}

// 삭제할 Target(훈련대상)의 시퀀스 넘버를 프론트엔드로 부터 받아오기 위한 변수
type TargetNumber struct {
	TargetNumber []string `json:"target_list"` //front javascript 와 이름을 일치시켜야함.
}

type Tag struct {
	TagNo			int 	`json:"tag_no"`
	TagName			string	`json:"tag_name"`
	TagCreateTime	string	`json:"created_t"`
}

func (t *Target) CreateTarget(conn *sql.DB, num int) error {

	t.TargetName = strings.Trim(t.TargetName, " ")
	t.TargetEmail = strings.Trim(t.TargetEmail, " ")
	t.TargetPhone = strings.Trim(t.TargetPhone, " ")
	t.TargetOrganize = strings.Trim(t.TargetOrganize, " ")
	t.TargetPosition = strings.Trim(t.TargetPosition, " ")

	if len(t.TargetName) < 1 {
		return fmt.Errorf("Target's name is empty ")
	} else if len(t.TargetEmail) < 1 {
		return fmt.Errorf(" Target's E-mail is empty ")
	} else if len(t.TargetPhone) < 1 {
		return fmt.Errorf(" Target's Phone number is empty ")
	} else if len(t.TargetOrganize) < 1 {
		return fmt.Errorf(" Target's Organize is empty")
	} else if len(t.TargetPosition) < 1 {
		return fmt.Errorf(" Target's Position is empty ")
	}

	// 추후 조건 좀더 꼼꼼하게 만들기..
	// ex) 엑셀파일 중간에 값이 비워져있는 경우 채워넣을 Default 값에 대한 조건 등...
	// 엑셀파일의 중간에 값이 없는 경우, 잘못된 형식이 들어가 있을경우 이를 검사할 필요가 있음.

	TagName := conn.QueryRow("SELECT tag_name FROM tag_info WHERE tag_no = $1", t.TagNo)
	errs := TagName.Scan(&t.TargetTag)
	if t.TagNo == 0 {
		t.TargetTag = "Null"
	} else if errs != nil{
		fmt.Println(errs)
		return fmt.Errorf("Tag's name Inquirying error. ")
	}

	query1 := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position, target_tag, user_no) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)" +
		"RETURNING target_no"

	row := conn.QueryRow(query1, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition, t.TargetTag, num)
	err := row.Scan(&t.TargetNo)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Target create error. ")
	}

	return nil
	}


// todo 보완필요!!! -> 현재 이름, 이메일, 태그 중 하나라도 값이 없으면 리스트목록에 뜨지않는 오류가 존재한다. 태그값이 없어도 표시되도록 해야함.
func ReadTarget(num int) ([]Target, []Tag, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, nil, fmt.Errorf("데이터베이스 연결 오류")
	}

	query := "SELECT target_name, target_email, target_phone, target_organize, target_position, target_tag, modified_time, target_no FROM target_info WHERE user_no = $1 ORDER BY target_no asc "
	// ROW_NUMBER() over (ORDER BY target_no) RNUM FROM target_info WHERE user_no = $1 ORDER BY target_no asc

	rows, err := db.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return nil, nil, fmt.Errorf("Targets query Error. ")
	}

	var targets []Target
	tg := Target{}
	for rows.Next() {
		err = rows.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag, &tg.TargetCreateTime, &tg.TargetNo)

		if err != nil {
			fmt.Printf("Targets scanning Error. : %v", err)
			continue
		}

		targets = append(targets, tg)
	}

	query2 := "SELECT tag_no, tag_name, modify_t FROM tag_info ORDER BY tag_no asc"
	tags, err2 := db.Query(query2)
	if err2 != nil{
		fmt.Println(err2)
		return nil, nil, fmt.Errorf("Tag query error. ")
	}

	var tag []Tag
	tg2 := Tag{}
	for tags.Next() {
		err2 = tags.Scan(&tg2.TagNo, &tg2.TagName, &tg2.TagCreateTime)

		if err2 != nil {
			fmt.Printf("Tags scanning Error. : %v", err)
			continue
		}

		tag = append(tag, tg2)
	}

	return targets, tag, nil
}

func (t *TargetNumber) DeleteTarget(conn *sql.DB, num int) error {

	for i := 0; i < len(t.TargetNumber); i++ {
		number, _ := strconv.Atoi(t.TargetNumber[i])

		if t.TargetNumber == nil {
			return fmt.Errorf("Please enter the number of the object to be deleted. ")
		}

		_, err := conn.Exec("DELETE FROM target_info WHERE user_no = $1 AND target_no = $2", num, number)
		if err != nil {
			return fmt.Errorf("Error deleting target ")
		}
	}

	return nil
}

// 반복해서 읽고 값을 넣는것을 메서드로 구현하고 API는 이걸 그냥 사용하기만 하면됨.
func (t *Target) ImportTargets(conn *sql.DB, str string, num int) error {

	// todo 2 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 2는 전부 같은 경로로 수정, api_Target.go 파일의 todo 2 참고)
	f, err := excelize.OpenFile("C:/Users/Taeho/go/src/redteam/Spreadsheet/"+ str)
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
		t.TargetTag	= f.GetCellValue("Sheet1", "F"+str)

		if t.TargetName == "" {
			break
		} else if len(t.TargetTag) < 1 {
			t.TargetTag = "Null"
		}

		//	todo 4 : 추후 해당 목록에 적힌 글들의 값이 올바른 형식이 아닐경우 제외하도록 하는 코드도 삽입한다. -> 정규식 사용.
		query := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position, target_tag, user_no) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7)" +
			"RETURNING target_no"

		row := conn.QueryRow(query, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition, t.TargetTag, num)
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
		}

		if len(tg.TargetTag) < 1{
			tg.TargetTag = "Null"
		}

		str := strconv.Itoa(i)
		f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
		f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
		f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
		f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
		f.SetCellValue("Sheet1", "E"+str, tg.TargetPosition)
		f.SetCellValue("Sheet1", "F"+str, tg.TargetTag)
		f.SetCellValue("Sheet1", "G"+str, tg.TargetCreateTime)

		i++
	}

	f.SetActiveSheet(index)

	str := strconv.Itoa(num) //int -> string

	// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 3은 전부 같은 경로로 수정, api_Target.go 파일의 todo 3 참고)
	// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
	// 파일 이름에 str변수 (
	if err2 := f.SaveAs("./Spreadsheet/Registered_Targets" + str + ".xlsx"); err != nil {
		fmt.Println(err2)
		return fmt.Errorf("Registered Target downloading Error. ")
	}

	return nil
}

func (t *Tag) CreateTag(conn *sql.DB) error {
	t.TagName = strings.Trim(t.TagName, " ")
	if len(t.TagName) < 1 {
		return fmt.Errorf(" Tag Name is empty. ")
	}

	_, err := conn.Exec("INSERT INTO tag_info(tag_name) VALUES ($1)", t.TagName)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Tag create error. ")
	}

	return nil
}

func (t *Tag) DeleteTag(conn *sql.DB) error {

	// num (int) -> str (string) 변환
	str := strconv.Itoa(t.TagNo)
	if str == "" {
		return fmt.Errorf("Please enter the number of the object to be deleted. ")
	}

	_, err := conn.Exec("DELETE FROM tag_info WHERE tag_no = $1", t.TagNo)
	if err != nil {
		return fmt.Errorf("Error deleting target ")
	}

	return nil
}