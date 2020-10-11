package model

import (
	"database/sql"
	"fmt"
	"strings"
)

type Template struct {
	TmpNo        int
	UserNo       int	`json:"userno"`		   // 사용자(사원) 번호
	Division     string `json:"division"`      // 구분
	TrainingKind string `json:"training_kind"` // 훈련유형
	FileInfo     string `json:"file_info"`     // 첨부파일 정보
	TmpName      string `json:"tmp_name"`	   // 템플릿 이름
	MailTitle    string `json:"mail_title"`	   // 메일 제목
	SenderName   string `json:"sender_name"`   // 보낸 사람
	DownloadType string `json:"download_type"` // 다운로드 파일 타입
}

func (t *Template) Create(conn *sql.DB, userID string) error {
	//conn *sql.DB, userID string
	//db, err := ConnectDb()
	//if err != nil {
	//	return fmt.Errorf("데이터베이스 연결 오류")
	//}
	//defer db.Close()

	t.TmpName = strings.Trim(t.TmpName, " ")
	if len(t.TmpName) < 1 {
		return fmt.Errorf("템플릿 이름이 비어있습니다. ")
	}

	query := "INSERT INTO template_info (user_no, tmp_division, tmp_kind, file_info," +
		" tmp_name, mail_title, sender_name, download_type) " +
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)" +
		" RETURNING tmp_no, user_no"


	//t.UserNo
	row := conn.QueryRow(query, t.UserNo, t.Division, t.TrainingKind, t.FileInfo, t.TmpName,
		t.MailTitle, userID, t.DownloadType)

	err := row.Scan(&t.TmpNo, &t.UserNo)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("템플릿 생성도중 오류가 발생하였습니다. ")
	}

	return nil
}

func ReadAll() ([]Template, error) {

	db, err := ConnectDb()
	if err != nil {
		return nil, fmt.Errorf("데이터베이스 연결 오류")
	}
	//defer db.Close()

	query := "SELECT tmp_no, user_no, tmp_division, tmp_kind, file_info, tmp_name," +
		" mail_title, sender_name, download_type FROM template_info"

	rows, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("템플릿을 읽어오는데 오류가 발생하였습니다 ")
	}

	var templates []Template
	for rows.Next() {
		tmp := Template{}
		err = rows.Scan(&tmp.TmpNo, &tmp.UserNo, &tmp.Division, &tmp.TrainingKind,
			&tmp.FileInfo, &tmp.TmpName, &tmp.MailTitle, &tmp.SenderName,
			&tmp.DownloadType)

		if err != nil {
			fmt.Printf("템플릿 스캐닝 오류 : %v", err)
			continue
		}
		templates = append(templates, tmp)
		// 읽어들인 값들을 전부 template 에 넣은 후에 반환하여 보여준다.
	}

	return templates, nil
}

