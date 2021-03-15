package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"regexp"
	"strconv"
)

// 프로젝트 생성 & 프로젝트 읽어오기에 사용하는 구조체
type Project struct {
	PNo          int      `json:"p_no"`          // 진짜 프로젝트 번호
	TmlNo        int      `json:"tml_no"`        // 템플릿 번호
	FakeNo       int      `json:"fake_no"`       // 화면에 표시될 프로젝트의 번호
	PName        string   `json:"p_name"`        // 프로젝트 이름
	PDescription string   `json:"p_description"` // 프로젝트 설명
	TagArray     []string `json:"tag_no"`        // 등록할 태그 대상자들
	PStatus      string   `json:"p_status"`      // 프로젝트 진행행태
	TemplateNo   string   `json:"tmp_no"`        // 적용할 템플릿 번호나 이름
	SendNo       int      `json:"send_no"`       // 메일 보낸 횟수
	Reading      string   `json:"reading"`       //읽은 사람
	Connect      string   `json:"connect"`
	Infection    string   `json:"infection"`  // 감염비율
	Targets      int      `json:"targets"`    // 훈련 대상자수
	StartDate    string   `json:"start_date"` // 프로젝트 시작일
	EndDate      string   `json:"end_date"`   // 프로젝트 종료일
}

//// 프로젝트 시작(Consumer)에서 사용하는 구조체
//type ProjectStart struct {
//	PNo            int    `json:"p_no"`
//	TmpNo          int    `json:"tmp_no"`
//	TargetNo       int    `json:"target_no"` // 훈련대상자들 번호
//	UserNo         int    `json:"user_no"`
//	TargetName     string `json:"target_name"`     // 훈련대상의 이름
//	TargetEmail    string `json:"target_email"`    // 훈련대상의 이메일주소
//	TargetOrganize string `json:"target_organize"` // 훈련대상의 소속
//	TargetPosition string `json:"target_position"` // 훈련대상의 직급
//	TargetPhone    string `json:"target_phone"`    // 훈련대상 전화번호
//	MailTitle      string `json:"mail_title"`      // 메일 제목
//	MailContent    string `json:"mail_content"`    // 메일 내용
//	SenderEmail    string `json:"sender_email"`    // 보내는사람(관리자) 이메일
//	SmtpHost       string `json:"smtp_host"`
//	SmtpPort       string `json:"smtp_port"`
//	SmtpId         string `json:"smtp_id"`
//	SmtpPw         string `json:"smtp_pw"`
//}

// 프로젝트 시작(Producer)에서 사용하는 구조체
type ProjectStart2 struct {
	PNo      int `json:"p_no"`
	TmpNo    int `json:"tmp_no"`
	TargetNo int `json:"target_no"`
	UserNo   int `json:"user_no"`
}

// 프로젝트 삭제에 사용하는 구조체
type ProjectDelete struct {
	ProjectNumber []string `json:"project_list"` //front javascript 와 이름을 일치시켜야함.
}

const (
	topic         = "redteam"
	brokerAddress = "localhost:9092"
)

//var Msg string

func (p *Project) ProjectCreate(conn *sql.DB, num int) (error, int) {

	ErrorCode := 200

	// 프로젝트 생성시 값이 제대로 들어오지 않은 경우 에러를 반환한다.
	if p.PName == "" || p.TemplateNo == "" || p.StartDate == "" || p.EndDate == "" || len(p.TagArray) < 1 {
		ErrorCode = 400
		return fmt.Errorf("Please enter all necessary information. "), ErrorCode
	}

	// 프로젝트 이름 형식검사
	var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{1,30}$", p.PName)
	if validName != true {
		ErrorCode = 400
		return fmt.Errorf("Project Name format is incorrect. "), ErrorCode
	}

	// 등록된 대상자 수를 조회한다.
	row := conn.QueryRow(`SELECT COUNT(p_no)
								FROM project_info
								WHERE user_no = $1;`, num)
	err := row.Scan(&p.PNo)
	if err != nil {
		ErrorCode = 500
		return fmt.Errorf("%v", err), ErrorCode
	}

	// 등록된 대상자 수 검사 (405에러)
	if p.PNo >= 10 {
		ErrorCode = 405
		return fmt.Errorf("The project is already full. "), ErrorCode
	}

	// 태그 중복제거
	keys := make(map[string]bool)
	ue := []string{}

	for _, value := range p.TagArray {
		if _, saveValue := keys[value]; !saveValue { // 중복제거 핵심포인트

			keys[value] = true
			ue = append(ue, value)
		}
	}

	p.TagArray = nil
	p.TagArray = ue

	// 태그 개수에 따른 입력
	switch len(p.TagArray) {
	case 1:
		var count int

		// 해당 태그를 가진 훈련대상자가 존재하는지 검증
		row := conn.QueryRow(`SELECT COUNT(target_no) as targets
										FROM target_info
										WHERE tag1 = $1 or tag2 = $1 or tag3 = $1`, p.TagArray[0])
		err := row.Scan(&count)
		if err != nil {
			ErrorCode = 500
			return fmt.Errorf("%v", err), ErrorCode // 에러 출력
		}

		if count < 1 {
			// 해당 태그를 가진 대상자가 존재하지 않는 경우
			ErrorCode = 402
			return fmt.Errorf("No target with the corresponding tag exists. "), ErrorCode

		} else {
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
 										tag1, tag2, tag3, sender_email, user_no) 
 										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING p_no;`
			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], 0, 0, "-", num)
			err = row.Scan(&p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}

			query = `INSERT INTO project_target_info (p_no, tml_no, target_no, user_no)
						 SELECT p_no, tml_no, target_no, T.user_no
						 FROM (select target_no, user_no, ti.tag1, ti.tag2, ti.tag3
								from target_info as ti
								where user_no = $1 AND (ti.tag1 = $2 or ti.tag2 = $2 or ti.tag3 = $2)) as T
								LEFT JOIN project_info pi on pi.user_no = T.user_no and pi.p_no = $3
						 WHERE T.tag1 >= 0 and (T.tag1 = pi.tag1 or T.tag1 = pi.tag2 or T.tag1 = pi.tag3)
									or T.tag2 >= 0 and (T.tag2 = pi.tag1 or T.tag2 = pi.tag2 or T.tag2 = pi.tag3)
									or T.tag3 >= 0 and (T.tag3 = pi.tag1 or T.tag3 = pi.tag2 or T.tag3 = pi.tag3)
						 ORDER BY p_no;`

			_, err = conn.Exec(query, num, p.TagArray[0], p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}
		}

	case 2:
		var count int

		// 해당 태그를 가진 훈련대상자가 존재하는지 검증
		row := conn.QueryRow(`SELECT COUNT(target_no) as targets
										FROM target_info
										WHERE tag1 = $1 or tag2 = $1 or tag3 = $1
												or tag1 = $2 or tag2 = $2 or tag3 = $2`, p.TagArray[0], p.TagArray[1])
		err := row.Scan(&count)
		if err != nil {
			ErrorCode = 500
			return fmt.Errorf("%v", err), ErrorCode // 에러 출력
		}

		if count < 1 {
			// 해당 태그를 가진 대상자가 존재하지 않는 경우
			ErrorCode = 402
			return fmt.Errorf("No target with the corresponding tag exists. "), ErrorCode

		} else {
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
 										tag1, tag2, tag3, sender_email, user_no) 
 										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING p_no;`

			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], 0, "-", num)
			err = row.Scan(&p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}

			query = `INSERT INTO project_target_info (p_no, tml_no, target_no, user_no)
						 SELECT p_no, tml_no, target_no, T.user_no
						 FROM (select target_no, user_no, ti.tag1, ti.tag2, ti.tag3
								from target_info as ti
								where user_no = $1 AND (ti.tag1 = $2 or ti.tag2 = $2 or ti.tag3 = $2
								    				or ti.tag1 = $3 or ti.tag2 = $3 or ti.tag3 = $3)) as T
								LEFT JOIN project_info pi on pi.user_no = T.user_no and pi.p_no = $4
						 WHERE T.tag1 >= 0 and (T.tag1 = pi.tag1 or T.tag1 = pi.tag2 or T.tag1 = pi.tag3)
									or T.tag2 >= 0 and (T.tag2 = pi.tag1 or T.tag2 = pi.tag2 or T.tag2 = pi.tag3)
									or T.tag3 >= 0 and (T.tag3 = pi.tag1 or T.tag3 = pi.tag2 or T.tag3 = pi.tag3)
						 ORDER BY p_no;`

			_, err = conn.Exec(query, num, p.TagArray[0], p.TagArray[1], p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}

		}
	case 3:
		var count int

		// 해당 태그를 가진 훈련대상자가 존재하는지 검증
		row := conn.QueryRow(`SELECT COUNT(target_no) as targets
										FROM target_info
										WHERE tag1 = $1 or tag2 = $1 or tag3 = $1
												or tag1 = $2 or tag2 = $2 or tag3 = $2
												or tag1 = $3 or tag2 = $3 or tag3 = $3`,
			p.TagArray[0], p.TagArray[1], p.TagArray[2])
		err := row.Scan(&count)
		if err != nil {
			ErrorCode = 500
			return fmt.Errorf("%v", err), ErrorCode // 에러 출력
		}

		if count < 1 {
			// 해당 태그를 가진 대상자가 존재하지 않는 경우
			ErrorCode = 402
			return fmt.Errorf("No target with the corresponding tag exists. "), ErrorCode

		} else {
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
										tag1, tag2, tag3, sender_email, user_no) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING p_no;`

			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], p.TagArray[2], "-", num)
			err = row.Scan(&p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}

			query = `INSERT INTO project_target_info (p_no, tml_no, target_no, user_no)
						SELECT p_no, tml_no, target_no, T.user_no
						FROM (select target_no, user_no, ti.tag1, ti.tag2, ti.tag3
							from target_info as ti
							where user_no = $1 AND (ti.tag1 = $2 or ti.tag2 = $2 or ti.tag3 = $2
												or ti.tag1 = $3 or ti.tag2 = $3 or ti.tag3 = $3
												or ti.tag1 = $4 or ti.tag2 = $4 or ti.tag3 = $4)) as T
							LEFT JOIN project_info pi on pi.user_no = T.user_no and pi.p_no = $5
							WHERE T.tag1 >= 0 and (T.tag1 = pi.tag1 or T.tag1 = pi.tag2 or T.tag1 = pi.tag3)
							or T.tag2 >= 0 and (T.tag2 = pi.tag1 or T.tag2 = pi.tag2 or T.tag2 = pi.tag3)
							or T.tag3 >= 0 and (T.tag3 = pi.tag1 or T.tag3 = pi.tag2 or T.tag3 = pi.tag3)
							ORDER BY p_no;`

			_, err = conn.Exec(query, num, p.TagArray[0], p.TagArray[1], p.TagArray[2], p.PNo)
			if err != nil {
				ErrorCode = 500
				return fmt.Errorf("Project Create error. "), ErrorCode
			}
		}

	}

	defer conn.Close()

	return nil, ErrorCode
}

func ReadProject(conn *sql.DB, num int) ([]Project, error) {
	// 프로젝트 읽어오기전에 해시테이블에 태그정보 한번 넣고 시작한다.
	var query string

	query = `SELECT row_num,
       T.p_no,
       tmp_no,
       tmp_name,
       p_name,
       to_char(p_start_date, 'YYYY-MM-DD'),
       to_char(p_end_date, 'YYYY-MM-DD'),
       COALESCE(tag_name1, '')                          as tag_name1,
       COALESCE(tag_name2, '')                          as tag_name2,
       COALESCE(tag_name3, '')                          as tag_name3,
       target_count                                     as Targets,
       COUNT(distinct ci.target_no)                     as Reading,
       COUNT(CASE WHEN ci.link_click_status THEN 1 END) as Connection,
       COUNT(CASE WHEN ci.download_status THEN 1 END)   as Infection,
       T.send_no,
       T.p_status
FROM (SELECT distinct ROW_NUMBER() over (ORDER BY p.p_no) AS row_num,
                      target_count,
                      p.p_no,
                      p.tml_no,
                      p_name,
                      p_start_date,
                      p_end_date,
                      p.tag1,
                      p.tag2,
                      p.tag3,
                      p.user_no,
                      p.send_no,
                      p.p_status
      FROM project_info as p
               LEFT JOIN (select count(target_no) as target_count, user_no, p_no
                          from project_target_info
                          group by user_no, p_no) as pti on p.user_no = pti.user_no AND p.p_no = pti.p_no
      WHERE p.user_no = $1
     ) AS T
         LEFT JOIN (SELECT tag_name as tag_name1, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti1 on ti1.tag_no = T.tag1
         LEFT JOIN (SELECT tag_name as tag_name2, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti2 on ti2.tag_no = T.tag2
         LEFT JOIN (SELECT tag_name as tag_name3, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti3 on ti3.tag_no = T.tag3
         LEFT JOIN template_info ti on ti.tmp_no = T.tml_no
         LEFT JOIN target_info ta on T.user_no = ta.user_no
         LEFT JOIN count_info ci on ta.target_no = ci.target_no AND T.p_no = ci.project_no
GROUP BY row_num, T.p_no, tmp_no, tmp_name, p_name, to_char(p_start_date, 'YYYY-MM-DD'),
         to_char(p_end_date, 'YYYY-MM-DD'), COALESCE(tag_name1, ''), COALESCE(tag_name2, ''), COALESCE(tag_name3, ''),
         target_count, T.send_no, T.p_status
ORDER BY row_num;`

	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading the projects. ")
	}

	// 태그들을 담을 배열
	var tags [3]string

	// Project 구조체를 값으로 가지는 배열
	var projects []Project

	for rows.Next() {
		p := Project{}
		// 가숫자, 진숫자, 템플릿, 템플릿 이름, 프로젝트 이름, 시작일, 종료일, 태그123, 대상자수, 읽은사람 수, 감염자 수, 보낸수, 플젝상태
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TmlNo, &p.TemplateNo, &p.PName, &p.StartDate, &p.EndDate,
			&tags[0], &tags[1], &tags[2], &p.Targets, &p.Reading, &p.Connect, &p.Infection, &p.SendNo, &p.PStatus)
		if err != nil {
			return nil, fmt.Errorf("Project scanning error : %v ", err)
		}

		// 태그이름 슬라이스에 바인딩
		for i := 0; i < len(tags); i++ {
			p.TagArray = append(p.TagArray, tags[i])
		}

		projects = append(projects, p)
	}

	// DB 커넥션 종료
	defer conn.Close()

	return projects, nil
}

// 프로젝트 종료버튼 누를경우 작동하는 기능, 현재는 사용안함.
//func (p *Project) EndProject(conn *sql.DB, num int) error {
//	_, err := conn.Exec(`UPDATE project_info
// 								SET p_status = 2
// 								WHERE user_no = $1 AND p_no = $2`,
// 								num, p.PNo)
//	if err != nil {
//		return fmt.Errorf("Error updating project status ")
//	}
//
//	return nil
//}

func (p *ProjectDelete) DeleteProject(conn *sql.DB, num int) error {

	for i := 0; i < len(p.ProjectNumber); i++ {
		number, _ := strconv.Atoi(p.ProjectNumber[i])

		if p.ProjectNumber == nil {
			return fmt.Errorf("Please enter the number of the object to be deleted. ")
		}

		// project_info 테이블에서 해당하는 프로젝트를 지운다.
		_, err := conn.Exec("DELETE FROM project_info WHERE user_no = $1 AND p_no = $2", num, number)
		if err != nil {
			return fmt.Errorf("Error deleting target ")
		}
	}

	defer conn.Close()

	return nil
}

// 사용자가 시작버튼을 누른경우에만 동작한다.
func (p *ProjectStart2) StartProject(conn *sql.DB, num int) error {

	// 프로젝트번호, 템플릿번호, 훈련 대상자번호, 사용자번호 만 조회해서 TOPIC 에 적재한다.
	query := `SELECT p_no, tml_no, target_no, user_no
				FROM project_target_info
				WHERE p_no = $1 and user_no = $2
				ORDER BY target_no;`

	rows, err := conn.Query(query, p.PNo, num)
	if err != nil {
		return fmt.Errorf("project starting error : %v", err)
	}

	// 카프카에 넣을 메일 내용
	// 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용
	w := kafka.Writer{
		Addr:  kafka.TCP(brokerAddress),
		Topic: topic,
	}

	msg := ProjectStart2{}

	// todo Kafka producer
	go func() {
		for rows.Next() {
			// DB 로부터 토픽에 작성할 내용들을 불러온다.
			err = rows.Scan(&msg.PNo, &msg.TmpNo, &msg.TargetNo, &msg.UserNo)
			if err != nil {
				_ = fmt.Errorf("%v", err)
			}

			// 만약 템플릿이 삭제된 템플릿일 경우 프로젝트의 상태를 오류로 변경한다.
			if msg.TmpNo == 0 {
				_, err = conn.Exec(`UPDATE project_info
 								SET p_start_date = now(), p_status = 3
 								WHERE user_no = $1 AND p_no = $2`, num, p.PNo)
				if err != nil {
					//_ = fmt.Errorf("%v", err)
					log.Println(err)
				}
				break
			}

			// 카프카에 작성할 내용들을 json 형식으로 변경하여 전송한다.
			message, _ := json.Marshal(msg)
			produce(message, w)
		}
	}()

	// 프로젝트 시작날짜를 오늘로 변경 & 프로젝트 상태를 진행으로 변경한다.
	_, err = conn.Exec(`UPDATE project_info SET p_start_date = now(), p_status = 1,
                        sender_email = (SELECT smtp_id FROM smtp_info WHERE user_no = $1)
                        WHERE user_no = $1 AND p_no = $2;`,
		num, p.PNo)
	if err != nil {
		return fmt.Errorf("Error : updating project status. ")
	}

	defer conn.Close()

	return nil
}

// todo Kafka producer function
func produce(messages []byte, w kafka.Writer) {
	err := w.WriteMessages(context.Background(), kafka.Message{
		//Key: []byte("Key"),
		Value: messages,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
}

func (p *Project) EndDateModify(conn *sql.DB, num int) (bool, error) {
	// 진행상태가 종료인 경우 종료일 변경 불가
	result := 0
	err := conn.QueryRow(`
		UPDATE project_info
		SET p_end_date = $1
		WHERE user_no = $2
		  AND p_no = $3
		  and p_end_date > now()
		  and p_end_date < $1
		  returning 1`, p.EndDate, num, p.PNo).Scan(&result)
	if result == 0 {
		return true, err
	}
	if err != nil {
		return false, err
	}

	defer conn.Close()
	return false, nil
}

func ProjectDetail(conn *sql.DB, userNo int, tmpNo int, pNo int) (Template, error) {
	var query = `SELECT tmp_no, tmp_division, tmp_kind, tmp_name, file_info,
       sender_email, mail_title, mail_content, download_type
	FROM template_info as ti
	    LEFT JOIN project_info pi on p_no = $1 and pi.user_no = $2
	WHERE tmp_no = $3 and ti.user_no = $4;`

	var userNo2 int

	// 기본템플릿 (0 - 4)은 user_no 가 0으로 설정되어있다.
	if tmpNo <= 4 && tmpNo >= 0 {
		userNo2 = 0
	} else {
		userNo2 = userNo
	}

	//var Detail []Template
	tmp := Template{}

	tmpDetail := conn.QueryRow(query, pNo, userNo, tmpNo, userNo2)
	// tmp.SenderName 은 smtp_info 테이블의 smtp_id 정보를 담는다.
	err := tmpDetail.Scan(&tmp.TmpNo, &tmp.Division, &tmp.Kind, &tmp.TmpName, &tmp.FileInfo, &tmp.SenderName,
		&tmp.MailTitle, &tmp.Content, &tmp.DownloadType)

	if err != nil {
		// 읽어온 정보를 바인딩하는데 오류가 발생.
		return Template{}, fmt.Errorf("Template detail scanning error : %v ", err)
	}

	//Detail = append(Detail, tmp)
	defer conn.Close()

	return tmp, nil
}
