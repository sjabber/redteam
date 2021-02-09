package model

import (
	"context"
	"crypto/aes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Connect		 string   `json:"connect"`
	Infection    string   `json:"infection"`     // 감염비율
	Targets      int      `json:"targets"`       // 훈련 대상자수
	StartDate    string   `json:"start_date"`    // 프로젝트 시작일
	EndDate      string   `json:"end_date"`      // 프로젝트 종료일
}

// 프로젝트 시작(Consumer)에서 사용하는 구조체
type ProjectStart struct {
	PNo            int    `json:"p_no"`
	TmpNo          int    `json:"tmp_no"`
	TargetNo       int    `json:"target_no"`       // 훈련대상자들 번호
	UserNo         int    `json:"user_no"`
	TargetName     string `json:"target_name"`     // 훈련대상의 이름
	TargetEmail    string `json:"target_email"`    // 훈련대상의 이메일주소
	TargetOrganize string `json:"target_organize"` // 훈련대상의 소속
	TargetPosition string `json:"target_position"` // 훈련대상의 직급
	TargetPhone    string `json:"target_phone"`    // 훈련대상 전화번호
	MailTitle      string `json:"mail_title"`      // 메일 제목
	MailContent    string `json:"mail_content"`    // 메일 내용
	SenderEmail    string `json:"sender_email"`    // 보내는사람(관리자) 이메일
	SmtpHost       string `json:"smtp_host"`
	SmtpPort       string `json:"smtp_port"`
	SmtpId         string `json:"smtp_id"`
	SmtpPw         string `json:"smtp_pw"`
}

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
 										tag1, tag2, tag3, user_no) 
 										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING p_no;`
			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], 0, 0, num)
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
 										tag1, tag2, tag3, user_no) 
 										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING p_no;`

			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], 0, num)
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
										tag1, tag2, tag3, user_no) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING p_no;`

			row = conn.QueryRow(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], p.TagArray[2], num)
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

	query = `SELECT tag_no, tag_name
			  FROM tag_info
			  WHERE user_no = $1
			  ORDER BY tag_no asc`
	hash, err := conn.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tg := Tag{}
	for hash.Next() {
		err = hash.Scan(&tg.TagNo, &tg.TagName)
		Hashmap[tg.TagNo] = tg.TagName

		if err != nil {
			fmt.Printf("Tags scanning Error. : %v", err)
			continue
		}
	}

	query = `SELECT row_num,
       				T.p_no,
       				tmp_no,
       				tmp_name,
       				p_name,
       				to_char(p_start_date, 'YYYY-MM-DD'),
       				to_char(p_end_date, 'YYYY-MM-DD'),
				    T.tag1,
				    T.tag2,
				    T.tag3,
                    COUNT(distinct pti.target_no),
       				COUNT(distinct ci.target_no) as Reading,
       				COUNT(CASE WHEN ci.link_click_status THEN 1 END) as Connection,
       				COUNT(CASE WHEN ci.download_status THEN 1 END) as Infection,
                    T.send_no,
       				T.p_status
			FROM (SELECT ROW_NUMBER() over (ORDER BY p_no) AS row_num,
					 p_no,
			         tmp_no,
					 ti.tmp_name,
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
					   LEFT JOIN template_info ti on p.tml_no = ti.tmp_no
			  	WHERE p.user_no = $1
			) AS T
			     LEFT JOIN project_target_info pti on T.user_no = pti.user_no AND T.p_no = pti.p_no
				 LEFT JOIN target_info ta on T.user_no = ta.user_no
				 LEFT JOIN count_info ci on ta.target_no = ci.target_no AND T.p_no = ci.project_no
		GROUP BY row_num, T.p_no, tmp_no, tmp_name, p_name, to_char(p_start_date, 'YYYY-MM-DD'),
				 to_char(p_end_date, 'YYYY-MM-DD'), T.tag1, T.tag2, T.tag3, T.send_no, T.p_status
		ORDER BY row_num;`

	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading the projects. ")
	}

	var tags [3]int
	var projects []Project // Project 구조체를 값으로 가지는 배열
	for rows.Next() {
		p := Project{}
		// 가숫자, 진숫자, 템플릿, 템플릿 이름, 프로젝트 이름, 시작일, 종료일, 태그123, 대상자수, 읽은사람 수, 감염자 수, 보낸수, 플젝상태
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TmlNo, &p.TemplateNo, &p.PName, &p.StartDate, &p.EndDate,
			&tags[0], &tags[1], &tags[2], &p.Targets, &p.Reading, &p.Connect, &p.Infection, &p.SendNo, &p.PStatus)
		if err != nil {
			return nil, fmt.Errorf("Project scanning error : %v ", err)
		}

		//Note 프로젝트 생성시 무조건 하나 이상의 태그가 들어가야 하기 때문에 하나 이상은 존재한다.
		// 그렇지 않을 경우 버그가 발생한다!!!
		//태그의 값을 이곳에 넣어준다.
	Loop1:
		for i := 0; i < len(tags); i++ {
			if tags[i] == 0 {
				p.TagArray = append(p.TagArray, "")
				continue Loop1
			}

			for key, val := range Hashmap {
				if key == tags[i] {
					p.TagArray = append(p.TagArray, val)
					break
				}
			}
		}

		projects = append(projects, p)
	}

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

//Note 프로젝트가 종료된 다음에 p_status 값 자동으로 변경시키는거 추가해야함.
// 자동으로 프로젝트 실행됐을 때 버튼눌렀을때와 동일하게 작동하도록도 수정해야함.
// 작동되는 시간 언제로 할지 논리적으로 따져서 결정하기.
//cron(스케쥴링)을 활용하여 매일 특정 시간에 프로젝트들의 진행상황을 체크한다.
func AutoStartProject() {
	wg := &sync.WaitGroup{}
	wg.Add(1) // WaitGroup 의 고루틴 개수 1개 증가

	kor, _ := time.LoadLocation("Asia/Seoul")
	c := cron.New(cron.WithLocation(kor))

	// 매 특정 시간마다 프로젝트들의 날짜를 점검하여 실행 종료시킨다.
	c.AddFunc("43 4 * * *", Auto) // 매일 정각에 프로젝트를 검토후 진행&종료시킨다.
	//c.AddFunc("48-50 19 * * *", Auto) //매일 12시 15-18분 마다 실행할 프로젝트를 검토한다.
	c.AddFunc("57-59 11 * * *", Auto) //하루 끝에 최종적으로 한번더 검토한다.
	c.Start()
	wg.Wait()
}

func Auto() {
	var Pno, num, status, startDate, endDate string

	// DB 연결
	conn, err := ConnectDB()
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
	defer conn.Close()

	query := `SELECT p_no, user_no, p_status, 
       				 to_char(p_start_date, 'YYYY-MM-DD') as start_date, to_char(p_end_date, 'YYYY-MM-DD') as end_date
			  FROM project_info`
	rows, err := conn.Query(query)
	if err != nil {
		//_ = fmt.Errorf("%v", err)

		log.Println(err)
		panic(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&Pno, &num, &status, &startDate, &endDate)
		if err != nil {
			//_ = fmt.Errorf("%v", err)
			log.Println(err)
			panic(err.Error())
		}

		// 날짜가 오늘 && 프로젝트가 예약상태인 경우 -> 프로젝트를 실행시킨다.
		if startDate == time.Now().Format("2006-01-02") && status == "0" {

			// 프로젝트번호, 템플릿번호, 훈련 대상자번호, 사용자번호 만 조회해서 TOPIC 에 적재한다.
			query2 := `SELECT p_no, tml_no, target_no, user_no
				FROM project_target_info
				WHERE p_no = $1 and user_no = $2
				ORDER BY target_no;`

			rows2, err := conn.Query(query2, Pno, num)
			if err != nil {
				//_ = fmt.Errorf("%v", err)
				log.Println(err)
				panic(err.Error())
			}

			//todo 카프카에 넣을 메일 내용
			// 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용
			w := kafka.Writer{
				Addr:  kafka.TCP(brokerAddress),
				Topic: topic,
			}

			msg := ProjectStart2{}

			// todo Kafka producer
			go func() {
				for rows2.Next() {
					// DB 로부터 토픽에 작성할 내용들을 불러온다.
					err = rows2.Scan(&msg.PNo, &msg.TmpNo, &msg.TargetNo, &msg.UserNo)
					if err != nil {
						//_ = fmt.Errorf("%v", err)
						log.Println(err)
						panic(err.Error())
					}

					// 카프카에 작성할 내용들을 json 형식으로 변경하여 전송한다.
					message, _ := json.Marshal(msg)
					produce(message, w)
				}
			}()

			// 프로젝트 시작날짜를 오늘로 변경 & 프로젝트 상태를 진행으로 변경한다.
			_, err = conn.Exec(`UPDATE project_info
 								SET p_status = 1
 								WHERE user_no = $1 AND p_no = $2`,
				num, Pno)
			if err != nil {
				//_ = fmt.Errorf("%v", err)
				log.Println(err)
				panic(err.Error())
			}

			// 프로젝트가 종료일이면 종료한다.
		} else if endDate == time.Now().Format("2006-01-02") && status == "0" || status == "1" {
			_, err = conn.Exec(`UPDATE project_info
 								SET p_status = 2
 								WHERE user_no = $1 AND p_no = $2`,
				num, Pno)
			if err != nil {
				//_ = fmt.Errorf("%v", err)
				log.Println(err)
				panic(err.Error())
			}
		} else {
			continue
		}
	}


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

	//todo 카프카에 넣을 메일 내용
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
					_ = fmt.Errorf("%v", err)
				}
				break
			}

			// 카프카에 작성할 내용들을 json 형식으로 변경하여 전송한다.
			message, _ := json.Marshal(msg)
			produce(message, w)
		}
	}()

	// 프로젝트 시작날짜를 오늘로 변경 & 프로젝트 상태를 진행으로 변경한다.
	_, err = conn.Exec(`UPDATE project_info
 								SET p_start_date = now(), p_status = 1
 								WHERE user_no = $1 AND p_no = $2`,
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

// todo kafka consumer function
func (p *ProjectStart) Consumer() {
	// todo Kafka consumer
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     "redteam",
		MinBytes:    5,                //5 바이트
		MaxBytes:    4132,             //1kB
		MaxWait:     3 * time.Second,  //3초만 기다린다.
		StartOffset: kafka.LastOffset, // GroupID 이전에 동일한 설정으로 데이터 사용한 적이
		// 있는 경우 중단한 곳부터 계속된다.
	})

	// 부득이하게 다른 DB connecting 방법 사용..
	conn, err := ConnectDB()
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용 순이다.

	chance := 0
	for {
		// ReadMessage 메서드는 우리가 다음 이벤트를 받을 때까지 차단된다.
		// json 객체로 받은 값들을 sendMail 메서드에 잘 적재시킨다.
		Msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("Could not read message : " + err.Error())
		}

		// json -> ProjectStart 객체 각각에 값으로 들어가도록 한다.
		json.Unmarshal(Msg.Value, p)

		row := conn.QueryRow(`SELECT target_name,
       								target_email,
       								target_organize,
       								target_position,
       								target_phone,
       								mail_title,
       								mail_content,
       								sender_name,
       								smtp_host,
       								smtp_port,
       								smtp_id,
       								smtp_pw
								FROM (SELECT target_name, target_email, target_organize, target_position, target_phone
    								  FROM target_info
    								  WHERE target_no = $1 AND user_no = $2) as T
        									LEFT JOIN template_info as tmp on tmp_no = $3
        									LEFT JOIN smtp_info si on si.user_no = $2`, p.TargetNo, p.UserNo, p.TmpNo)

		// 메일 전송에 필요한 정보들 바인딩
		err = row.Scan(&p.TargetName, &p.TargetEmail, &p.TargetOrganize, &p.TargetPosition, &p.TargetPhone, &p.MailTitle,
			&p.MailContent, &p.SenderEmail, &p.SmtpHost, &p.SmtpPort, &p.SmtpId, &p.SmtpPw)
		if err == sql.ErrNoRows {
			// 해당 대상자가 존재하지 않는 경우 // Note 이거 지우자 그냥. (project_info 테이블에서 un_send_no도 없애자)
			_, err = conn.Exec(`UPDATE project_info 
									SET un_send_no = project_info.un_send_no + 1 
									WHERE user_no = $1 AND p_no = $2`, p.UserNo, p.PNo)
			if err != nil {
				//continue // 에러나면 그냥 스킵해버리기...
				panic(err.Error())
			}
			log.Println("error occurred 1")
			continue // 대상자가 없는 경우에는 보내지 못한 메일의 개수를 하나 올리고 다음으로 넘어간다.

		} else if err != nil {
			// 정말 알 수 없는 에러가 난 케이스
			_, err = conn.Exec(`UPDATE project_info 
									SET p_status = 3
									WHERE user_no = $1 AND p_no = $2`, p.UserNo, p.PNo)
			log.Println("error occurred 2")
			continue // 이 경우 그냥 스킵
			//panic(err.Error()) // 이 경우 프로세스 중지.
		}


		// 파싱작업 수행
		p.MailContent = parsing(p.MailContent, p.TargetName, p.TargetOrganize,
			p.TargetPosition, p.TargetPhone, strconv.Itoa(p.TargetNo), strconv.Itoa(p.PNo))

		// smtp 패스워드 복호화 작업 수행
		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			panic(err.Error())
		}
		password, _ := base64.StdEncoding.DecodeString(p.SmtpPw)
		password = Decrypt(block, password)
		p.SmtpPw = string(password)

		// 가공된 메일 전송
		err = p.sendMail()
		if err != nil {
			panic("Could not send email : " + err.Error())
		}

		// 메일 보낸 수 +1
		_, err = conn.Exec(`UPDATE project_info 
								SET send_no = project_info.send_no + 1
								WHERE user_no = $1 AND p_no = $2`, p.UserNo, p.PNo)
		if err != nil {
			panic("Could not send email " + err.Error())
		}

		// 메일 3개 보내면 3분 간격으로 쉬어준다.
		chance++
		if chance == 3 {
			time.Sleep(3 * time.Minute)
			chance = 0
		}
	}
}

func (p *ProjectStart) sendMail() error {

	// string -> int, smtp 연결을 테스트한다.
	port, _ := strconv.Atoi(p.SmtpPort) //string -> int
	d := gomail.NewDialer(p.SmtpHost, port, p.SmtpId, p.SmtpPw)
	s, err := d.Dial()
	if err != nil {
		return fmt.Errorf("1번 : " + err.Error())
	}

	m := gomail.NewMessage()
	m.SetHeader("From", p.SenderEmail)
	m.SetAddressHeader("To", p.TargetEmail, p.TargetName)
	m.SetHeader("Subject", p.MailTitle)
	m.SetBody("text/html", fmt.Sprintf(p.MailContent))

	//send close
	if err := gomail.Send(s, m); err != nil {
		return fmt.Errorf("2번 : " + err.Error())
	}

	time.Sleep(10 * time.Second)
	m.Reset()
	return nil
}

// 메일 내용 파싱함수
// 이 파싱 메서드가 성능이 어떤지는 아직 제대로 점검해보지 않았음..
//메일내용, 대상이름, 대상소속, 대상직급, 대상전화번호, 대상번호, 프로젝트번호
func parsing(str string, str1 string, str2 string, str3 string, str4 string, str5 string, str6 string) string {

	if strings.Contains(str, "{{target_name}}") {
		str = strings.Replace(str, "{{target_name}}", str1, -1)
	}

	if strings.Contains(str, "{{target_organize}}") {
		str = strings.Replace(str, "{{target_organize}}", str2, -1)
	}

	if strings.Contains(str, "{{target_position}}") {
		str = strings.Replace(str, "{{target_position}}", str3, -1)
	}

	if strings.Contains(str, "{{target_phone}}") {
		str = strings.Replace(str, "{{target_phone}}", str4, -1)
	}

	// todo 추후 도메인 추가로 필요, 그때는 파싱이 아니라 접속한 사이트에 넣어야함. (접속, 감염 두개 추가필요!)
	if strings.Contains(str, "{{count_ip}}") {
		s := "<html>\n<body>\n<img src=\"http://localhost:5000/api/CountTarget?" +
			"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\">\n" +
			"<a href=\"http://localhost:5000/api/CountTarget?" +
			"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\"></a>\n</body>\n</html>"

		str = strings.Replace(str, "{{count_ip}}", s, -1)
	}

	return str
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
