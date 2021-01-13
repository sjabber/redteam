package model

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"regexp"
	"strconv"
	"strings"
)

// Target(훈련대상)을 관리하기 위한 json 구조체
type Target struct {
	TargetNo         int       `json:"tg_no"`
	TargetName       string    `json:"tg_name"`
	TargetEmail      string    `json:"tg_email"`
	TargetPhone      string    `json:"tg_phone"`
	TargetOrganize   string    `json:"tg_organize"` //소속
	TargetPosition   string    `json:"tg_position"` //직급
	TargetTag        [3]string `json:"tg_tag"`      //태그, 추후에 slice 로 변경한다..
	TargetCreateTime string    `json:"created_t"`
	TagArray         []string  `json:"tag_no"` // 태그 입력받을 때 사용
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

// 대상관리 페이지 맨 처음 GetTag가 실행되기 때문에
// 태그들을 전역변수인 해시맵에 담아서 트랜잭션 횟수를 줄인다.
var Hashmap = make(map[int]string) // 태그값들을 담아놓을 변수.

//
func isValueIn(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (t *Target) CreateTarget(conn *sql.DB, num int) (int, error) {

	// status 200 설정. 에러발생시 변경됨.
	errcode := 200

	t.TargetName = strings.Trim(t.TargetName, " ")
	t.TargetEmail = strings.Trim(t.TargetEmail, " ")
	t.TargetPhone = strings.Trim(t.TargetPhone, " ")
	t.TargetOrganize = strings.Trim(t.TargetOrganize, " ")
	t.TargetPosition = strings.Trim(t.TargetPosition, " ")

	if len(t.TargetName) < 1 {
		errcode = 400
		return errcode, fmt.Errorf("Target's name is empty ")
	} else if len(t.TargetEmail) < 1 {
		errcode = 400
		return errcode, fmt.Errorf(" Target's E-mail is empty ")
	} else if len(t.TargetPhone) < 1 {
		errcode = 400
		return errcode, fmt.Errorf(" Target's Phone number is empty ")
	} else if len(t.TargetOrganize) < 1 {
		errcode = 400
		return errcode, fmt.Errorf(" Target's Organize is empty ")
	} else if len(t.TargetPosition) < 1 {
		errcode = 400
		return errcode, fmt.Errorf(" Target's Position is empty ")
	}

	//else if len(t.TagArray) < 1 {
	//	return fmt.Errorf(" Target's Tag is empty ")
	//}

	// 추후 조건 좀더 꼼꼼하게 만들기..
	// ex) 엑셀파일 중간에 값이 비워져있는 경우 채워넣을 Default 값에 대한 조건 등...

	// 이메일 형식검사
	var validEmail, _ = regexp.MatchString(
		"^[_A-Za-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", t.TargetEmail)

	if validEmail != true {
		errcode = 402
		return errcode, fmt.Errorf("Email format is incorrect. ")
	}

	// 이름 형식검사 (한글, 영어 이름만 허용)
	var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{1,30}$", t.TargetName)

	if validName != true {
		errcode = 402
		return errcode, fmt.Errorf("Name format is incorrect. ")
	}

	// 핸드폰 형식검사
	var phoneNumber, _ = regexp.MatchString(
		"^[0-9]{9,11}$", t.TargetPhone)

	if phoneNumber != true {
		errcode = 402
		return errcode, fmt.Errorf("Phone number format is incorrect. ")
	}

	// t.TagArray 값이 비어있으면 에러나는 관계로 값을 채워준다.
	for i := 1; i <= 3; i++ {
		if len(t.TagArray) < i {
			t.TagArray = append(t.TagArray, "0")
		}
	}

	// 엑셀파일의 중간에 값이 없는 경우, 잘못된 형식이 들어가 있을경우 이를 검사할 필요가 있음.

	query1 := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position," +
		"tag1, tag2, tag3, user_no) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err := conn.Exec(query1, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition,
		t.TagArray[0], t.TagArray[1], t.TagArray[2], num)
	if err != nil {
		errcode = 500
		fmt.Println(err)
		return errcode, fmt.Errorf("Target create error. ")
	}

	return errcode, nil
}

// todo 보완필요!!! -> 현재 이름, 이메일, 태그 중 하나라도 값이 없으면 리스트목록에 뜨지않는 오류가 존재한다. 태그값이 없어도 표시되도록 해야함.
func ReadTarget(conn *sql.DB, num int, page int) ([]Target, int, int, error) {
	var pageNum int // 몇번째 페이지부터 가져올지 결정하는 변수
	var pages int   // 총 페이지 수
	var total int   // 총 훈련대상자들의 수를 담을 변수

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
       tag1,
       tag2,
       tag3,
       to_char(modified_time, 'YYYY-MM-DD'),
       target_no
    FROM (SELECT ROW_NUMBER() over (ORDER BY target_no) AS row_num,
             target_no,
             target_name,
             target_email,
             target_phone,
             target_organize,
             target_position,
             tag1,
             tag2,
             tag3,
             modified_time
          FROM target_info
          WHERE user_no = $1
         ) AS T
    WHERE row_num > $2
    ORDER BY target_no asc
    LIMIT 20;
`
	// 조건에 맞는 데이터를 조회한다.
	rows, err := conn.Query(query, num, pageNum)
	if err != nil {
		fmt.Println(err)
		return nil, 0, 0, fmt.Errorf("Target's query Error. ")
	}

	var targets []Target
	tg := Target{}

	//태그의 번호가 담길 변수
	var tags [3]int

	// 목록들을 하나하나 읽어들여온다.
	for rows.Next() {
		err = rows.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tags[0], &tags[1], &tags[2], &tg.TargetCreateTime, &tg.TargetNo)
		if err != nil {
			fmt.Printf("Targets scanning Error. : %v", err)
			continue
		}

		// Note 전화번호에 하이픈(-)을 추가하여 사용자에게 보여준다.
		var sub [3]string
		phone := []rune(tg.TargetPhone)

		if len(tg.TargetPhone) < 10 {
			sub[0] = string(phone[0:2])
			sub[1] = string(phone[2:5])
			sub[2] = string(phone[5:9])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if string(phone[1:2]) == "2" && len(tg.TargetPhone) == 10 {
			sub[0] = string(phone[0:2])
			sub[1] = string(phone[2:6])
			sub[2] = string(phone[6:10])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if len(tg.TargetPhone) == 10 {
			sub[0] = string(phone[0:3])
			sub[1] = string(phone[3:6])
			sub[2] = string(phone[6:10])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if len(tg.TargetPhone) == 11 {
			sub[0] = string(phone[0:3])
			sub[1] = string(phone[3:7])
			sub[2] = string(phone[7:11])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		}

		// Note 해당 대상(타겟)의 태그값을 여기서 읽어들어온다.
		// 데이터베이스에서 일일히 조회할 필요없이 GetTag()를 통해 가져오고 저장한 값과 비교해 조회속도를 높인다.
	Loop1:
		for i := 0; i < len(tags); i++ {
			if tags[i] == 0 {
				tg.TargetTag[i] = ""
				continue Loop1
			}

			for key, val := range Hashmap {
				if key == tags[i] {
					tg.TargetTag[i] = val
					break
				} else {
					tg.TargetTag[i] = ""
				}
			}
		}

		targets = append(targets, tg)

		tg.TargetTag[0] = ""
		tg.TargetTag[1] = ""
		tg.TargetTag[2] = "" // slice 로 변경되면 다른 방식으로 값을 비운다.
	}// for문 끝.

	// 전체 타겟(훈련대상)의 수를 반환한다.
	query = `
    select count(target_no) 
    from target_info 
    where user_no = $1`

	pageCount := conn.QueryRow(query, num)
	_ = pageCount.Scan(&total) // 훈련 대상자들의 전체 수를 pages 에 바인딩.

	pages = (total / 20) + 1 // 전체훈련 대상자들을 토대로 전체 페이지수를 계산한다.

	// 각각 표시할 대상 20개, 대상의 총 갯수, 총 페이지 수, 에러를 반환한다.
	return targets, total, pages, nil
}

func (t *TargetNumber) DeleteTarget(conn *sql.DB, num int) error {

	for i := 0; i < len(t.TargetNumber); i++ {
		number, _ := strconv.Atoi(t.TargetNumber[i])

		if t.TargetNumber == nil {
			return fmt.Errorf("Please enter the number of the object to be deleted. ")
		}

		// target_info 테이블에서 대상을 지운다.
		_, err := conn.Exec("DELETE FROM target_info WHERE user_no = $1 AND target_no = $2", num, number)
		if err != nil {
			return fmt.Errorf("Error deleting target ")
		}
	}

	return nil
}

// 반복해서 읽고 값을 넣는것을 메서드로 구현하고 API는 이걸 그냥 사용하기만 하면됨.
// Excel 파일로부터 대상의 정보를 일괄적으로 읽어 DB에 등록한다.
func (t *Target) ImportTargets(conn *sql.DB, uploadPath string, num int) error {

	// str -> 일괄등록하기 위한 업로드 경로 + 파일 이름이 담기는 변수
	f, err := excelize.OpenFile(uploadPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	user1 := strconv.Itoa(num) //int -> string

	// Bulk insert 하기 위해 값들을 쌓아놓을 변수
	var BigString string

	// 태그 값을 db 조회없이 비교, 대조하여 삽입하여 성능을 높이기 위한 변수
	var list []string
	pt := 0
	for _,  value := range Hashmap {
		list = append(list, value)
		pt++
	}


	i := 2 // 2행부터 값을 읽어온다.
	for i <= 501 {
		str := strconv.Itoa(i)

		t.TargetName = f.GetCellValue("Sheet1", "A"+str)
		t.TargetEmail = f.GetCellValue("Sheet1", "B"+str)

		// 이메일 형식검사
		var validEmail, _ = regexp.MatchString(
			"^[_A-Za-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", t.TargetEmail)

		// 이름 형식검사 (한글, 영어 이름만 허용)
		var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{2,30}$", t.TargetName)

		// 필수적인 정보가 누락됐거나 형식이 잘못된 경우 그 즉시 입력을 중단한다.
		if validName != true || t.TargetName == "" {
			break
		} else if validEmail != true || t.TargetEmail == "" {
			// todo 추후 이메일 중복도 체크해야한다.
			break
		}

		t.TargetPhone = f.GetCellValue("Sheet1", "C"+str)

		// 핸드폰 형식검사
		var validPhone, _ = regexp.MatchString(
			"^[0-9]{9,11}$", t.TargetPhone)

		// 핸드폰 번호 형식이 올바르지 않을 경우에는 공백처리한다.
		if validPhone != true {
			t.TargetPhone = ""
		}

		var sub [3]string

		// 여기부터는 선택정보.
		t.TargetOrganize = f.GetCellValue("Sheet1", "D"+str)
		t.TargetPosition = f.GetCellValue("Sheet1", "E"+str)
		sub[0] = f.GetCellValue("Sheet1", "F"+str)
		sub[1] = f.GetCellValue("sheet1", "G"+str)
		sub[2] = f.GetCellValue("sheet1", "H"+str)

	Loop1:
		for i := 0; i < len(sub); i++ {
			if isValueIn(sub[i], list) {
			Loop2:
				for key, val := range Hashmap {
					if val == sub[i] {
						sub[i] = strconv.Itoa(key)
						break Loop2
					}
				}
			} else {
				sub[i] = "0"
				continue Loop1
			}
		}

		//Bulk insert로 삽입할 내용들을 텍스트로 만든다.
		//Note xlsx 파일은 psql 에서 인코딩문제로 bulk insert 불가, csv, txt 등은 가능함.
		BigString += "('" + t.TargetName + "', '" + t.TargetEmail + "', '" +
			t.TargetPhone + "', '" + t.TargetOrganize + "', '" + t.TargetPosition + "', " +
			user1 + "," + sub[0] + "," + sub[1] + "," + sub[2] + ")," + "\n"

		i++
	}

	BigString = BigString[:len(BigString)-2]

	query := "INSERT INTO target_info (target_name, target_email, target_phone," +
		"target_organize, target_position, user_no, tag1, tag2, tag3) VALUES" +
		BigString

	_, err = conn.Exec(query)
	if err != nil {
		fmt.Println(err)
	}

	//bulkFile.Close()

	return nil
}

// DB에 저장된 값들을 읽어 엘셀파일에 일괄적으로 작성하여 저장한다.
func ExportTargets(conn *sql.DB, num int, tagNumber int) error {

	var tags [3]int
	// tagNumber 가 0인 경우 (전체 선택)
	if tagNumber == 0 {
		query := `
         SELECT target_no, target_name, target_email, target_phone, target_organize, target_position, 
         		to_char(modified_time, 'YYYY-MM-YY HH24:MI'), tag1, tag2, tag3
         from target_info
         WHERE user_no = $1
         ORDER BY target_no`

		rows, err := conn.Query(query, num)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Database error. ")
		}


		// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 1은 전부 같은 경로로 수정, api_Target.go 파일의 todo 1 참고)
		// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
		// 서버에 있는 sample 파일에 내용을 작성한 다음 다른 이름의 파일로 클라이언트에게 전송한다.
		f, err := excelize.OpenFile("./Spreadsheet/sample.xlsx")
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("Open Spreadsheet Error. ")
		}
		index := f.NewSheet("Sheet1")

		i := 2
		for rows.Next() {
			tg := Target{}
			err = rows.Scan(&tg.TargetNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
				&tg.TargetPosition, &tg.TargetCreateTime, &tags[0], &tags[1], &tags[2])
			if err != nil {
				fmt.Printf("Target scanning error : %v", err)
				continue
			}

			//query = `SELECT tag_no FROM tag_target_info WHERE user_no = $1 AND target_no = $2`
			//rows2, err2 := conn.Query(query, num, tg.TargetNo)
			//if err != nil {
			//	fmt.Println(err)
			//	return fmt.Errorf("Database error. ")
			//}

			//k := 0 // 태그의 인덱스를 담을 변수
			//var tagNum string

		Loop1:
			for i := 0; i < len(tags); i++ {
				if tags[i] == 0 {
					tg.TargetTag[i] = ""
					continue Loop1
				}

				for key, val := range Hashmap {
					if key == tags[i] {
						tg.TargetTag[i] = val
						break
					} else {
						tg.TargetTag[i] = ""
					}
				}
			}

			//for k := 0; k < len(tags); k++ { // Note 이중 포문으로 태그의 이름을 조회한다.
			//	if tags[k] == "0" {
			//		tg.TargetTag[k] = ""
			//		continue
			//	}
			//
			//	tagName := conn.QueryRow(`SELECT tag_name FROM tag_info WHERE tag_no = $1`, tags[k])
			//	err = tagName.Scan(&tg.TargetTag[k])
			//	if err != nil {
			//		_ = fmt.Errorf("Target's Tag number query Error. ")
			//		continue
			//	}
			//}

			str := strconv.Itoa(i)
			f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
			f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
			f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
			f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
			f.SetCellValue("Sheet1", "E"+str, tg.TargetPosition)
			f.SetCellValue("Sheet1", "F"+str, tg.TargetTag[0])
			f.SetCellValue("Sheet1", "G"+str, tg.TargetTag[1])
			f.SetCellValue("Sheet1", "H"+str, tg.TargetTag[2])
			f.SetCellValue("Sheet1", "I"+str, tg.TargetCreateTime)

			i++

			// 태그의 값을 마지막엔 비워준다.
			tg.TargetTag[0] = ""
			tg.TargetTag[1] = ""
			tg.TargetTag[2] = "" // slice 로 변경되면 다른 방식으로 값을 비운다.
		}

		f.SetActiveSheet(index)

		str := strconv.Itoa(num) //int -> string


		// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 3은 전부 같은 경로로 수정, api_Target.go 파일의 todo 3 참고)
		// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
		// 파일 이름에 str변수 (
		if err2 := f.SaveAs("./Spreadsheet/" + str + "/Registered_Targets.xlsx"); err2 != nil {
			fmt.Println(err2)
			return fmt.Errorf("Registered Target downloading Error. ")
		}

		return nil

		// todo -------------------아래부턴 특정 태그만 골라서 내보낼 경우에 해당함.-----------------------------------
	} else {

		query := `SELECT target_name, target_email, target_phone, target_organize, target_position,
 						 to_char(modified_time, 'YYYY-MM-YY HH24:MI'), tag1, tag2, tag3
				FROM (SELECT target_name,
             				 target_email,
            				 target_phone,
             				 target_organize,
             				 target_position,
             				 modified_time,
             				 tag1,
             				 tag2,
             				 tag3
      			FROM target_info
      			WHERE user_no = $1) as T
				WHERE tag1 = $2
   				OR tag2 = $2
   				OR tag3 = $2`

		result, err := conn.Query(query, num, tagNumber)
		if err != nil {
			fmt.Println(err)
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

		for result.Next() {
			tg := Target{}

			// 해당 태그에 속하는 대상들을 하나하나 가져온다.
			err = result.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone,
				&tg.TargetOrganize, &tg.TargetPosition, &tg.TargetCreateTime,
				&tags[0], &tags[1], &tags[2]) //조회한 값들을 하나하나 바인딩
			if err != nil {
				_ = fmt.Errorf("Target number scanning Error. ")
				continue
			}

		Loop2:
			for i := 0; i < len(tags); i++ {
				if tags[i] == 0 {
					tg.TargetTag[i] = ""
					continue Loop2
				}

				for key, val := range Hashmap {
					if key == tags[i] {
						tg.TargetTag[i] = val
						break
					} else {
						tg.TargetTag[i] = ""
					}
				}
			}

			str := strconv.Itoa(i)
			f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
			f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
			f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
			f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
			f.SetCellValue("Sheet1", "E"+str, tg.TargetPosition)
			f.SetCellValue("Sheet1", "F"+str, tg.TargetTag[0])
			f.SetCellValue("Sheet1", "G"+str, tg.TargetTag[1])
			f.SetCellValue("Sheet1", "H"+str, tg.TargetTag[2])
			f.SetCellValue("Sheet1", "I"+str, tg.TargetCreateTime)

			i++
		}
		f.SetActiveSheet(index)

		str := strconv.Itoa(num) //int -> string

		// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 3은 전부 같은 경로로 수정, api_Target.go 파일의 todo 3 참고)
		// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
		// 파일 이름에 str변수 (
		if err2 := f.SaveAs("./Spreadsheet/" + str + "/Registered_Targets.xlsx"); err != nil {
			fmt.Println(err2)
			return fmt.Errorf("Registered Target downloading Error. ")
		}
	}

	return nil
}

func (t *Tag) CreateTag(conn *sql.DB, num int) error {
	t.TagName = strings.Trim(t.TagName, " ")
	if len(t.TagName) < 1 {
		return fmt.Errorf(" Tag Name is empty. ")
	}

	_, err := conn.Exec("INSERT INTO tag_info(tag_name, user_no) VALUES ($1, $2)", t.TagName, num)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Tag create error. ")
	}

	return nil
}

func (t *Tag) DeleteTag(conn *sql.DB, num int) error {
	_, err := conn.Exec("DELETE FROM tag_info WHERE tag_no = $1 AND user_no = $2", t.TagNo, num)
	if err != nil {
		return fmt.Errorf("Error deleting tag on tag_info ")
	}

	return nil
}

// todo 4 : tag_info 에서 사용자 번호로 태그정보를 가져온다.
func GetTag(num int) []Tag {
	db, err := ConnectDB()
	if err != nil {
		return nil
	}

	var query string

	query = `SELECT tag_no, tag_name, to_char(modified_time, 'YYYY-MM-DD')
			  FROM tag_info
			  WHERE user_no = $1
			  ORDER BY tag_no asc
`
	tags, err := db.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var tag []Tag
	tg := Tag{}

	i := 0
	for tags.Next() {
		err = tags.Scan(&tg.TagNo, &tg.TagName, &tg.TagCreateTime)
		Hashmap[tg.TagNo] = tg.TagName

		if err != nil {
			fmt.Printf("Tags scanning Error. : %v", err)
			continue
		}

		tag = append(tag, tg)
		i++
	}

	//for key, val := range Hashmap {
	//	fmt.Println(key, val)
	//}

	return tag
}

func SearchTarget(conn *sql.DB, num int, page int, searchDivision string, searchText string) ([]Target, int, int, error) {
	var pageNum int // 몇번째 페이지부터 가져올지 결정하는 변수
	var pages int   // 총 페이지 수
	var total int   // 총 훈련대상자들의 수를 담을 변수

	// ex) 1페이지 -> 1~10, 2페이지 -> 11~20
	// 페이지번호에 따라 가져올 목록이 달라진다.
	pageNum = (page - 1) * 20

	// 대상목록들을 20개씩만 잘라서 반하여 페이징처리한다.
	query := "SELECT target_name, " +
		"target_email, " +
		"target_phone, " +
		"target_organize, " +
		"target_position, " +
		"tag1, " +
		"tag2, " +
		"tag3, " +
		"to_char(modified_time, 'YYYY-MM-DD')," +
		"target_no " +
		"FROM (SELECT ROW_NUMBER() over (ORDER BY target_no) AS row_num," +
		"target_no, " +
		"target_name, " +
		"target_email, " +
		"target_phone, " +
		"target_organize, " +
		"target_position, " +
		"tag1, " +
		"tag2, " +
		"tag3, " +
		"modified_time " +
		"FROM target_info " +
		"WHERE user_no = $1 AND " + searchDivision + " LIKE $2" +
		") AS T " +
		"WHERE row_num > $3 " +
		"ORDER BY target_no asc " +
		"LIMIT 20;"

	searchText = "%" + searchText + "%"
	rows, err := conn.Query(query, num, searchText, pageNum)
	if err != nil {
		fmt.Println(err)
		return nil, 0, 0, fmt.Errorf("Target's query Error. ")
	}

	var targets []Target
	tg := Target{}

	var tags [3]int

	for rows.Next() { // 목록들을 하나하나 읽어들여온다.
		err = rows.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tags[0], &tags[1], &tags[2], &tg.TargetCreateTime, &tg.TargetNo)
		if err != nil {
			fmt.Printf("Targets scanning Error. : %v", err)
			continue
		}

		// Note 전화번호에 하이픈(-)을 추가하여 사용자에게 보여준다.
		var sub [3]string
		phone := []rune(tg.TargetPhone)

		if len(tg.TargetPhone) < 10 {
			sub[0] = string(phone[0:2])
			sub[1] = string(phone[2:5])
			sub[2] = string(phone[5:9])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if string(phone[1:2]) == "2" && len(tg.TargetPhone) == 10 {
			sub[0] = string(phone[0:2])
			sub[1] = string(phone[2:6])
			sub[2] = string(phone[6:10])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if len(tg.TargetPhone) == 10 {
			sub[0] = string(phone[0:3])
			sub[1] = string(phone[3:6])
			sub[2] = string(phone[6:10])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		} else if len(tg.TargetPhone) == 11 {
			sub[0] = string(phone[0:3])
			sub[1] = string(phone[3:7])
			sub[2] = string(phone[7:11])
			tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		}

		// Note 해당 대상(타겟)의 태그값을 여기서 읽어들어온다.
		// 데이터베이스에서 일일히 조회할 필요없이 GetTag()를 통해 가져오고 저장한 값과 비교해 조회속도를 높인다.
	Loop1:
		for i := 0; i < len(tags); i++ {
			if tags[i] == 0 {
				tg.TargetTag[i] = ""
				continue Loop1
			}

			for key, val := range Hashmap {
				if key == tags[i] {
					tg.TargetTag[i] = val
					break
				} else {
					tg.TargetTag[i] = ""
				}
			}
		}

		targets = append(targets, tg)

		tg.TargetTag[0] = ""
		tg.TargetTag[1] = ""
		tg.TargetTag[2] = "" // slice 로 변경되면 다른 방식으로 값을 비운다.
	}

	// 전체 타겟(훈련대상)의 수를 반환한다.
	query = "SELECT count(target_no) " +
		"FROM target_info " +
		"WHERE user_no = $1 AND " + searchDivision + " LIKE $2"

	pageCount := conn.QueryRow(query, num, searchText)
	_ = pageCount.Scan(&total) // 훈련 대상자들의 전체 수를 pages 에 바인딩.

	pages = (total / 20) + 1 // 전체훈련 대상자들을 토대로 전체 페이지수를 계산한다.

	// 각각 표시할 대상 20개, 대상의 총 갯수, 총 페이지 수, 에러를 반환한다.
	return targets, total, pages, nil
}