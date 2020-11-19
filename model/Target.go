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
	TargetTag        string `json:"tg_tag"`      //태그
	TargetCreateTime string `json:"created_time"`
	TagNo            string `json:"tag_no"`
}

// 삭제할 Target(훈련대상)의 시퀀스 넘버를 프론트엔드로 부터 받아오기 위한 변수
type TargetNumber struct {
	TargetNumber []string `json:"target_list"` //front javascript 와 이름을 일치시켜야함.
}

type Tag struct {
	TagNo         int    `json:"tag_no"`
	TagName       string `json:"tag_name"`
	TagCreateTime string `json:"created_t"`
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
	if t.TagNo == "" {
		t.TargetTag = "Null"
	} else if errs != nil {
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
func ReadTarget(num int, page int) ([]Target, int, int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("DB connection error")
	}

	var pageNum int // 몇번째 페이지부터 가져올지 결정하는 변수
	var pages int // 총 페이지 수
	var total int // 총 훈련대상자들의 수를 담을 변수


	// ex) 1페이지 -> 1~10, 2페이지 -> 11~20
	// 페이지번호에 따라 가져올 목록이 달라진다.
	pageNum = (page - 1) * 20

	// 대상목록들을 20개씩만 잘라서 반하여 페이징처리한다.
	query := `
    SELECT
       target_name,
       target_email,
       target_phone,
       target_organize,
       target_position,
       target_tag,
       modified_time,
       target_no
    FROM (SELECT ROW_NUMBER() over (ORDER BY target_no) AS row_num,
             target_no,
             target_name,
             target_email,
             target_phone,
             target_organize,
             target_position,
             target_tag,
             modified_time
          FROM target_info
          WHERE user_no = $1
         ) AS T
    WHERE row_num > $2
    LIMIT 20;
`
	rows, err := db.Query(query, num, pageNum)
	if err != nil {
		fmt.Println(err)
		return nil, 0, 0, fmt.Errorf("Targets query Error. ")
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
		// 10개로 제한된 대상이 구조체에 담긴다.
		targets = append(targets, tg)
	}

	// 전체 타겟(훈련대상)의 수를 반환한다.
	query = `
    select count(target_no) 
    from target_info 
    where user_no = $1`

	page_count := db.QueryRow(query, num)
	_ = page_count.Scan(&pages) // 훈련 대상자들의 전체 수를 pages 에 바인딩.

	total = (pages / 20) + 1 // 전체훈련 대상자들을 토대로 전체 페이지수를 계산한다.

	// 각각 표시할 대상 20개, 대상의 총 갯수, 총 페이지 수, 에러를 반환한다.
	return targets, pages, total, nil
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

	// str -> 일괄등록하기 위한 업로드 경로가 담기는 변수
	f, err := excelize.OpenFile(str)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	i := 2

	for {
		str := strconv.Itoa(i)

		t.TargetName = f.GetCellValue("Sheet1", "A"+str)
		t.TargetEmail = f.GetCellValue("Sheet1", "B"+str)
		t.TargetPhone = f.GetCellValue("Sheet1", "C"+str)
		t.TargetOrganize = f.GetCellValue("Sheet1", "D"+str)
		t.TargetPosition = f.GetCellValue("Sheet1", "E"+str)
		t.TargetTag = f.GetCellValue("Sheet1", "F"+str)

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

	i := 2
	// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 1은 전부 같은 경로로 수정, api_Target.go 파일의 todo 1 참고)
	// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
	// 서버에 있는 sample 파일에 내용을 작성한 다음 다른 이름의 파일로 클라이언트에게 전송한다.
	f, err := excelize.OpenFile("./Spreadsheet/sample.xlsx")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Open Spreadsheet Error. ")
	}

	index := f.NewSheet("Sheet1")

	for rows.Next() {
		tg := Target{}
		err = rows.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag, &tg.TargetCreateTime)
		if err != nil {
			fmt.Printf("Targets scanning error : %v", err)
			continue
		}

		if len(tg.TargetTag) < 1 {
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


// todo 4 : 대상 / 태그 조인 테이블로 부터 태그를 가져오도록 수정한다.
// 그전에 조인테이블을 만들어야겠지?

func GetTag(num int) []Tag {
	db, err := ConnectDB()
	if err != nil {
		return nil
	}

	query := "SELECT tag_no, tag_name, modified_time FROM tag_info WHERE user_no = $1 ORDER BY tag_no asc"
	tags, err := db.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var tag []Tag
	tg := Tag{}
	for tags.Next() {
		err = tags.Scan(&tg.TagNo, &tg.TagName, &tg.TagCreateTime)

		if err != nil {
			fmt.Printf("Tags scanning Error. : %v", err)
			continue
		}

		tag = append(tag, tg)
	}

	return tag
}
