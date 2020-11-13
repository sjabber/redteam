package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Template struct {
	TmpNo int `json:"tmp_no"`
	//UserNo       int	`json:"user_no"`		   // 사용자(사원) 번호
	Division     string `json:"division"`      // 구분
	Kind         string `json:"kind"`          // 훈련유형
	FileInfo     string `json:"file_info"`     // 첨부파일 정보
	TmpName      string `json:"tmp_name"`      // 템플릿 이름
	MailTitle    string `json:"mail_title"`    // 메일 제목
	SenderName   string `json:"sender_name"`   // 보낸 사람
	DownloadType string `json:"download_type"` // 다운로드 파일 타입
}

//템플릿 생성 메서드, json 형식으로 데이터를 입력받아서 DB에 저장한다.
func (t *Template) Create(conn *sql.DB, userID string) error {
	t.TmpName = strings.Trim(t.TmpName, " ")
	if len(t.TmpName) < 1 {
		return fmt.Errorf("The template name is empty. ")
		// 템플릿 이름이 비어있습니다.
	}

	query := "INSERT INTO template_info (user_no, tmp_division, tmp_kind, file_info," +
		" tmp_name, mail_title, sender_name, download_type) " +
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)" +
		" RETURNING tmp_no, user_no"

	row := conn.QueryRow(query, //t.UserNo,
		t.Division, t.Kind, t.FileInfo, t.TmpName,
		t.MailTitle, userID, t.DownloadType)

	err := row.Scan(&t.TmpNo)//&t.UserNo

	if err != nil {
		return fmt.Errorf("An error occurred while creating the template. ")
		// 템플릿을 생성하던 중에 오류가 발생하였습니다.
	}

	return nil
}

// 템플릿 조회 메서드, 템플릿 테이블(template_info)의 모든 템플릿을 조회한다.
// 여기서 유일하게 솔루션을 사용하는 사용자들이 사용하게 되는 매서드
// 사용자들을 위해 http status를 정의한다.
func ReadAll() ([]Template, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("DB connecting error. ")
		// DB를 연결하던 중에 오류가 발생하였습니다.
	}
	//defer db.Close()

	query := "SELECT tmp_no, tmp_division, tmp_kind, file_info, tmp_name," +
		" mail_title, sender_name, download_type FROM template_info"

	rows, err := db.Query(query)

	if err != nil {
		// 템플릿을 DB 로부터 읽어오는데 오류가 발생.
		return nil, fmt.Errorf("There was an error reading the template. ")
	}

	var templates []Template
	for rows.Next() {
		tmp := Template{}
		err = rows.Scan(&tmp.TmpNo, &tmp.Division, &tmp.Kind,
			&tmp.FileInfo, &tmp.TmpName, &tmp.MailTitle, &tmp.SenderName,
			&tmp.DownloadType)

		if err != nil {
			// 읽어온 정보를 바인딩하는데 오류가 발생.
			return nil, fmt.Errorf("Template scanning error : %v ", err)
		}
		templates = append(templates, tmp)
		// 읽어들인 값들을 전부 template 배열에 넣은 후에 반환하여 보여준다.
	}

	return templates, nil
}

// 템플릿 수정 메서드, 템플릿 번호(tmp_no)에 해당하는 템플릿을 수정한다.
func (t *Template) Update(conn *sql.DB) error {
	str := string(t.TmpNo) // int -> string 형변환
	if str == "" {
		return fmt.Errorf("Please enter the template number to be modified. ")
	}
	now := time.Now()
	_, err := conn.Exec("UPDATE template_info SET tmp_division=$1, tmp_kind=$2,"+
		"file_info=$3, tmp_name=$4, mail_title=$5, sender_name=$6, download_type=$7, modify_t=$8 "+
		"WHERE tmp_no = $9", t.Division, t.Kind, t.FileInfo,
		t.TmpName, t.MailTitle, t.SenderName, t.DownloadType, now, t.TmpNo)

	if err != nil {
		return fmt.Errorf("Error updating template ")
		// 템플릿 업데이트 오류
	}

	return nil
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
