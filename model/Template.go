package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Template struct {
	FakeNo			int `json:"fake_no"`
	TmpNo 			int `json:"tmp_no"`
	//UserNo       int	`json:"user_no"`		  // 사용자(사원) 번호
	Division       string `json:"tmp_division"`  // 구분
	Kind           string `json:"tmp_kind"`     // 훈련유형
	FileInfo       string `json:"file_info"`   // 첨부파일 정보
	TmpName        string `json:"tmp_name"`   // 템플릿 이름
	MailTitle      string `json:"title"`     // 메일 제목
	SenderName     string `json:"sender_name"` // 보낸 사람
	Content		   string `json:"content"`  // 메일내용
	DownloadType   string `json:"download_type"` // 다운로드 파일 타입
	CreatedTime    string `json:"created_time"`  // 생성시간
	CreateRealTime time.Time
}

//템플릿 생성 메서드, json 형식으로 데이터를 입력받아서 DB에 저장한다.
//func (t *Template) Create(conn *sql.DB, userID string) error {
//	t.TmpName = strings.Trim(t.TmpName, " ")
//	if len(t.TmpName) < 1 {
//		return fmt.Errorf("The template name is empty. ")
//		// 템플릿 이름이 비어있습니다.
//	}
//
//	query := "INSERT INTO template_info (user_no, tmp_division, tmp_kind, file_info," +
//		" tmp_name, mail_title, sender_name, download_type) " +
//		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)" +
//		" RETURNING tmp_no, user_no"
//
//	row := conn.QueryRow(query, //t.UserNo,
//		t.Division, t.Kind, t.FileInfo, t.TmpName,
//		t.MailTitle, userID, t.DownloadType)
//
//	err := row.Scan(&t.TmpNo) //&t.UserNo
//
//	if err != nil {
//		return fmt.Errorf("An error occurred while creating the template. ")
//		// 템플릿을 생성하던 중에 오류가 발생하였습니다.
//	}
//
//	return nil
//}

// 템플릿 조회 메서드, 템플릿 테이블(template_info)의 모든 템플릿을 조회한다.
// 여기서 유일하게 솔루션을 사용하는 사용자들이 사용하게 되는 매서드
// 사용자들을 위해 http status를 정의한다.
func ReadAll(num int) ([]Template, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("DB connecting error. ")
		// DB를 연결하던 중에 오류가 발생하였습니다.
	}
	//defer db.Close()

	query := `SELECT 
	   row_num,
	   tmp_no,
       tmp_division,
       tmp_kind,
       file_info,
       tmp_name,
       mail_title,
       sender_name,
       download_type,
       created_time
FROM (SELECT ROW_NUMBER() over (ORDER BY tmp_no) AS row_num,
			 tmp_no,
             tmp_division,
             tmp_kind,
             file_info,
             tmp_name,
             mail_title,
             sender_name,
             download_type,
             created_time
      FROM template_info
      WHERE user_no = 0 OR user_no = $1
     ) AS T
ORDER BY row_num;`

	rows, err := db.Query(query, num)

	if err != nil {
		// 템플릿을 DB 로부터 읽어오는데 오류가 발생.
		return nil, fmt.Errorf("There was an error reading the template. ")
	}

	var templates []Template
	for rows.Next() {
		tmp := Template{}
		err = rows.Scan(&tmp.FakeNo, &tmp.TmpNo, &tmp.Division, &tmp.Kind,
			&tmp.FileInfo, &tmp.TmpName, &tmp.MailTitle, &tmp.SenderName,
			&tmp.DownloadType, &tmp.CreateRealTime)

		if err != nil {
			// 읽어온 정보를 바인딩하는데 오류가 발생.
			return nil, fmt.Errorf("Template scanning error : %v ", err)
		}
		tmp.CreatedTime = tmp.CreateRealTime.Format("2006-01-02 15:04")

		switch tmp.Division {
		case "1":
			tmp.Division = "기본"
		case "2":
			tmp.Division = "사용자"
		}

		switch tmp.Kind {
		case "1":
			tmp.Kind = "경고 안내"
		case "2":
			tmp.Kind = "피싱 유도"
		case "3":
			tmp.Kind = "실태 조사"
		}

		switch tmp.FileInfo {
		case "1":
			tmp.FileInfo = "EXE"
		case "2":
			tmp.FileInfo = "HTML"
		case "3":
			tmp.FileInfo = "Excel"
		}

		switch tmp.DownloadType {
		case "1":
			tmp.DownloadType = "링크 첨부"
		case "2":
			tmp.DownloadType = "파일 첨부"
		}

		templates = append(templates, tmp)
		// 읽어들인 값들을 전부 template 배열에 넣은 후에 반환하여 보여준다.
	}

	return templates, nil
}

// 템플릿 수정 메서드, 템플릿 번호(tmp_no)에 해당하는 템플릿을 수정한다.
func (t *Template) Update(conn *sql.DB, num int) error {

	//switch t.Division {
	//case "기본":
	//	t.Division = "1"
	//case "사용자":
	//	t.Division = "2"
	//}
	//
	//switch t.Kind {
	//case "경고 안내":
	//	t.Kind = "1"
	//case "피싱 유도":
	//	t.Kind = "2"
	//case "실태조사":
	//	t.Kind = "3"
	//}
	//
	//switch t.FileInfo {
	//case "EXE":
	//	t.FileInfo = "1"
	//case "HTML":
	//	t.FileInfo = "2"
	//case "Excel":
	//	t.FileInfo = "3"
	//}
	//
	//if t.DownloadType == "링크 첨부" {
	//	t.DownloadType = "1"
	//} else if t.DownloadType == "파일 첨부" {
	//	t.DownloadType = "2"
	//}

	query := `INSERT INTO template_info(tmp_division, tmp_kind, file_info, tmp_name,
 	mail_title, mail_content, sender_name, download_type, user_no)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := conn.Exec(query, t.Division, t.Kind, t.FileInfo, t.TmpName, t.MailTitle,
			t.Content, t.SenderName, t.DownloadType, num)

	if err != nil {
		log.Panic(err)
		return fmt.Errorf("Error updating template ")
		// 템플릿 업데이트 오류
	}

	return nil
}

func Detail(tmp_no int) (Template, error) {
	db, err := ConnectDB()
	if err != nil {
		return Template{}, fmt.Errorf("DB connecting error. ")
		// DB를 연결하던 중에 오류가 발생하였습니다.
	}

	query := `SELECT tmp_division, tmp_kind, file_info, sender_name, mail_title, mail_content, download_type
	FROM template_info
	WHERE tmp_no = $1`

	//var Detail []Template
	tmp := Template{}

	tmpDetail := db.QueryRow(query, tmp_no)
	err = tmpDetail.Scan(&tmp.Division, &tmp.Kind, &tmp.FileInfo, &tmp.SenderName,
		&tmp.MailTitle, &tmp.Content, &tmp.DownloadType)

	if err != nil {
		// 읽어온 정보를 바인딩하는데 오류가 발생.
		return Template{}, fmt.Errorf("Template detail scanning error : %v ", err)
	}

	//Detail = append(Detail, tmp)

	return tmp, nil
}

// 템플릿 삭제 메서드, 템플릿 번호(tmp_no)에 해당하는 템플릿을 삭제한다.
func (t *Template) Delete(conn *sql.DB) error {
	str := string(t.TmpNo) // int -> string 형변환
	if str == "" {
		return fmt.Errorf("Please enter the template number to be deleted. ")
		// 삭제할 템플릿 번호를 입력해주세요.
	}
	_, err := conn.Exec("DELETE ROWS FROM template_info WHERE tmp_no = $1", t.TmpNo)
	if err != nil {
		fmt.Printf("Error updating template: (%v)", err)
		return fmt.Errorf("Error deleting template ")
	}

	return nil
}

// 템플릿 테이블의 모든 정보를 삭제한다. -> 아직 template API에는 적용안한상태.
func (t *Template) DeleteAll(conn *sql.DB) error {
	_, err := conn.Exec("DELETE FROM template_info")
	if err != nil {
		fmt.Printf("Error updating template: (%v)", err)
		return fmt.Errorf("Error deleting template ")
	}

	return nil

}
