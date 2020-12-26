package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Template struct {
	FakeNo int `json:"fake_no"`
	TmpNo  int `json:"tmp_no"`
	//UserNo       int	`json:"user_no"`		  // 사용자(사원) 번호
	Division       string `json:"tmp_division"`  // 구분
	Kind           string `json:"tmp_kind"`      // 훈련유형
	FileInfo       string `json:"file_info"`     // 첨부파일 정보
	TmpName        string `json:"tmp_name"`      // 템플릿 이름
	MailTitle      string `json:"title"`         // 메일 제목
	SenderName     string `json:"sender_name"`   // 보낸 사람
	Content        string `json:"content"`       // 메일내용
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
func ReadAll(conn *sql.DB, num int) ([]Template, error) {

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

	rows, err := conn.Query(query, num)

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

	// t.TmpNo가 3이하면 템플릿이 생성되도록 코드를 짜고
	// 4이상이면 수정이 되도록 코드를 수정한다.
	var query string

	if t.TmpNo <= 3 {
		//Note 템플릿 생성
		query = `INSERT INTO template_info(tmp_division, tmp_kind, file_info, tmp_name,
 	mail_title, mail_content, sender_name, download_type, user_no)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err := conn.Exec(query, 2, t.Kind, t.FileInfo, t.TmpName, t.MailTitle,
			t.Content, t.SenderName, t.DownloadType, num)

		if err != nil {
			log.Panic(err)
			return fmt.Errorf("Error create template. ")
			// 템플릿 업데이트 오류
		}
		return nil

	} else if t.TmpNo >= 4 {
		//Note 템플릿 수정
		query = `UPDATE template_info SET tmp_division = $1, tmp_kind = $2, file_info = $3, tmp_name = $4,
                        sender_name = $5, mail_title = $6, mail_content = $7, 
                         download_type = $8 WHERE user_no = $9 AND tmp_no = $10;`

		_, err := conn.Exec(query, 2, t.Kind, t.FileInfo, t.TmpName, t.SenderName, t.MailTitle, t.Content,
			t.DownloadType, num, t.TmpNo)

		if err != nil {
			log.Panic(err)
			return fmt.Errorf("Error updating template. ")
			// 템플릿 업데이트 오류
		}
		return nil
	}

	return nil
}

func Detail(conn *sql.DB, userNo int, tmpNo int) (Template, error) {
	var query = `SELECT tmp_no, tmp_division, tmp_kind, tmp_name, file_info, sender_name, mail_title,
 	mail_content, download_type
	FROM template_info
	WHERE tmp_no = $1 and user_no = $2`

	if tmpNo <= 3 && tmpNo >= 0 {
		userNo = 0
	}

	//var Detail []Template
	tmp := Template{}

	tmpDetail := conn.QueryRow(query, tmpNo, userNo)
	err := tmpDetail.Scan(&tmp.TmpNo, &tmp.Division, &tmp.Kind, &tmp.TmpName, &tmp.FileInfo, &tmp.SenderName,
		&tmp.MailTitle, &tmp.Content, &tmp.DownloadType)

	if err != nil {
		// 읽어온 정보를 바인딩하는데 오류가 발생.
		return Template{}, fmt.Errorf("Template detail scanning error : %v ", err)
	}

	//Detail = append(Detail, tmp)

	return tmp, nil
}

// 템플릿 삭제 메서드, 템플릿 번호(tmp_no)에 해당하는 템플릿을 삭제한다.
func (t *Template) Delete(conn *sql.DB, userNo int) error {
	str := string(t.TmpNo) // int -> string 형변환
	if str == "" {
		return fmt.Errorf("Please enter the template number to be deleted. ")
		// 삭제할 템플릿 번호를 입력해주세요.
	}
	//Note 사용자번호(user_no)에 막혀서 기본 템플릿은 삭제가 되지 않는다. // 기본템플릿은 user_no가 0이기 때문.
	_, err := conn.Exec("DELETE FROM template_info WHERE tmp_no = $1 and user_no = $2", t.TmpNo, userNo)
	if err != nil {
		fmt.Printf("Error deleting template: (%v)", err)
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
