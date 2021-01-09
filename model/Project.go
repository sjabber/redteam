package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	PDescription string   `json:"p_description"` // 프로젝트 설명
	PNo          int      `json:"p_no"`          // 진짜 프로젝트 번호
	FakeNo       int      `json:"fake_no"`       // 화면에 표시될 프로젝트의 번호
	PName        string   `json:"p_name"`        // 프로젝트 이름
	TagArray     []string `json:"tag_no"`        // 등록할 태그 대상자들
	PStatus      string   `json:"p_status"`      // 프로젝트 진행행태
	TemplateNo   string   `json:"tmp_no"`        // 적용할 템플릿 번호나 이름
	Infection    int      `json:"infection"`     // 감염비율
	Targets      int      `json:"targets"`       // 훈련 대상자수
	StartDate    string   `json:"start_date"`    // 프로젝트 시작일
	EndDate      string   `json:"end_date"`      // 프로젝트 종료일
}

type ProjectStart struct {
	PNo            int `json:"p_no"`
	TargetNo       int
	TargetName     string
	TargetEmail    string
	TargetOrganize string
	TargetPosition string
	TargetPhone    string
	MailTitle      string
	MailContent    string
	SenderEmail    string
	SmtpHost       string
	SmtpPort       string
	SmtpId         string
	SmtpPw         string
}

const (
	topic = "redteam"
	brokerAddress = "localhost:9092"
	partition = 1
)

func (p *Project) ProjectCreate(conn *sql.DB, num int) error {

	// 프로젝트 생성시 값이 제대로 들어오지 않은 경우 에러를 반환한다.
	if p.PName == "" || p.TemplateNo == "" || p.StartDate == "" || p.EndDate == "" || len(p.TagArray) < 1 {
		return fmt.Errorf("Please enter all necessary information. ")
	}

	// 프로젝트 이름 형식검사
	var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{1,30}$", p.PName)
	if validName != true {
		return fmt.Errorf("Project Name format is incorrect. ")
	}

	switch len(p.TagArray) {
		case 1:
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
 										tag1, tag2, tag3, p_status, user_no) 
 										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

			_, err := conn.Exec(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], 0, 0, 0, num)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Project Create error. ")
			}
		case 2:
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
										tag1, tag2, tag3, p_status, user_no) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

			_, err := conn.Exec(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], 0, 0, num)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Project Create error. ")
			}
		case 3:
			query := `INSERT INTO project_info (p_name, p_description, p_start_date, p_end_date, tml_no,
										tag1, tag2, tag3, p_status, user_no) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

			_, err := conn.Exec(query, p.PName, p.PDescription, p.StartDate, p.EndDate, p.TemplateNo,
				p.TagArray[0], p.TagArray[1], p.TagArray[2], 0, num)
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("Project Create error. ")
			}
	}

	return nil
}

func ReadProject(conn *sql.DB, num int) ([]Project ,error) {
	query := `SELECT row_num,
       p_no,
       tmp_name,
       p_name,
       p_status,
       to_char(p_start_date, 'YYYY-MM-DD'),
       to_char(p_end_date, 'YYYY-MM-DD'),
       T.tag1,
       T.tag2,
       T.tag3,
       T.infection,
       COUNT(target_no)
FROM (SELECT ROW_NUMBER() over (ORDER BY p_no) AS row_num,
             p_no,
             ti.tmp_name,
             p_name,
             p_status,
             p_start_date,
             p_end_date,
             p.tag1,
             p.tag2,
             p.tag3,
             p.infection
      FROM project_info as p
               LEFT JOIN template_info ti on p.tml_no = ti.tmp_no
      WHERE p.user_no = $1
     ) AS T
         LEFT JOIN target_info ta on user_no = ta.user_no
WHERE T.tag1 > 0 and (T.tag1 = ta.tag1 or T.tag1 = ta.tag2 or T.tag1 = ta.tag3)
   or T.tag2 > 0 and (T.tag2 = ta.tag1 or T.tag2 = ta.tag2 or T.tag2 = ta.tag3)
   or T.tag3 > 0 and (T.tag3 = ta.tag1 or T.tag3 = ta.tag2 or T.tag3 = ta.tag3)
GROUP BY row_num, p_no, tmp_name, p_name, p_status, to_char(p_start_date, 'YYYY-MM-DD'),
         to_char(p_end_date, 'YYYY-MM-DD'), T.tag1, T.tag2, T.tag3, T.infection	
ORDER BY row_num;`

	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading the projects. ")
	}

	var tags [3]int
	var projects []Project // Project 구조체를 값으로 가지는 배열
	for rows.Next() {
		p := Project{}
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TemplateNo, &p.PName, &p.PStatus,
			&p.StartDate, &p.EndDate, &tags[0], &tags[1], &tags[2], &p.Infection, &p.Targets)
		if err != nil {
			return nil, fmt.Errorf("Project scanning error : %v ", err)
		}

		//Note 프로젝트 생성시 무조건 하나 이상의 태그가 들어가야 하기 때문에 하나 이상은 존재한다.
		// 그렇지 않을 경우 버그가 발생한다!!!
		//태그의 값을 이곳에 넣어준다.
		Loop1 :
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

	return projects, nil
}

func (p *Project) EndProject(conn *sql.DB, num int) error {
	_, err := conn.Exec(`UPDATE project_info
 								SET p_status = 2
 								WHERE user_no = $1 AND p_no = $2`,
 								num, p.PNo)
	if err != nil {
		return fmt.Errorf("Error updating project status ")
	}

	return nil
}

func (p *ProjectStart) StartProject(conn *sql.DB, num int) error {
	// 프로젝트 상태를 1로 변경하며 프로젝트를 실행한다.
	_, err := conn.Exec(`UPDATE project_info
 								SET p_status = 1
 								WHERE user_no = $1 AND p_no = $2`,
		num, p.PNo)
	if err != nil {
		return fmt.Errorf("Error : updating project status. ")
	}
	return nil
}

func (p *ProjectStart) Kafka(conn *sql.DB, num int) error {
	query := `SELECT distinct mail_title,
               mail_content,
               sender_name,
               target_name,
               target_email,
               target_organize,
               target_position,
               target_phone,
               target_no
			FROM (SELECT tml_no, tag1, tag2, tag3, user_no
     			FROM project_info
     			WHERE p_no = $1
       			AND user_no = $2) as T
        			LEFT JOIN template_info tmp on tml_no = tmp.tmp_no
        			LEFT JOIN target_info ta on T.user_no = ta.user_no
			WHERE T.tag1 > 0 and (T.tag1 = ta.tag1 or T.tag1 = ta.tag2 or T.tag1 = ta.tag3)
   				or T.tag2 > 0 and (T.tag2 = ta.tag1 or T.tag2 = ta.tag2 or T.tag2 = ta.tag3)
   				or T.tag3 > 0 and (T.tag3 = ta.tag1 or T.tag3 = ta.tag2 or T.tag3 = ta.tag3)
   				GROUP BY mail_title, mail_content, sender_name, target_name, target_email, target_organize, target_position,
         				 target_phone, target_no
				ORDER BY target_no;`

	rows, err := conn.Query(query, p.PNo, num)
	if err != nil {
		return fmt.Errorf("project starting error : %v", err)
	}

	//todo 카프카
	// 카프카에 넣을 메일 내용
	// 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용
	w := kafka.Writer{
		Addr:  kafka.TCP(brokerAddress),
		Topic: topic,
	}

	i := 0
	for rows.Next() {
		// DB 로부터 토픽에 작성할 내용들을 불러온다.
		err = rows.Scan(&p.MailTitle, &p.MailContent, &p.SenderEmail,
			&p.TargetName, &p.TargetEmail, &p.TargetOrganize,
			&p.TargetPosition, &p.TargetPhone, &p.TargetNo)
		if err != nil {
			fmt.Errorf("Error : sql error ")
		}

		// 파싱작업 수행
		p.MailContent = Parsing(p.MailContent, p.TargetName, p.TargetOrganize,
			p.TargetPosition, p.TargetPhone)

		// 카프카에 작성할 내용들을 하나의 띄어쓰기로 구분짓고 하나의 string 으로 묶는다.
		// 각각 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용 순이다.
		message := p.SenderEmail + "//" + p.TargetEmail + "//" + p.TargetName + "//" + p.MailTitle + "//" + p.MailContent

		// todo Kafka producer
		err = w.WriteMessages(context.Background(), kafka.Message{
			//Key: []byte("Key"),
			Value: []byte(message),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		i++
	}

	// 메일을 보내기 위한 정보를 DB로 부터 가져온다.
	row := conn.QueryRow(`SELECT smtp_host, smtp_port, smtp_id, smtp_pw
 								FROM smtp_info
 								WHERE user_no = $1`, num)
	err = row.Scan(&p.SmtpHost, &p.SmtpPort, &p.SmtpId, &p.SmtpPw)
	if err != nil {
		return fmt.Errorf("this account does not exist. ")
	}
	// string -> int, smtp 연결을 테스트한다.
	port, _ := strconv.Atoi(p.SmtpPort)
	d := gomail.NewDialer(p.SmtpHost, port, p.SmtpId, p.SmtpPw)
	_, err = d.Dial()
	if err != nil {
		return fmt.Errorf("Smtp connecting failed. : %v ", err)
	}

	// todo Kafka consumer
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic: topic,
		GroupID: "redteam",
		MinBytes: 5, //5 바이트
		MaxBytes: 4132,//1kB
		MaxWait: 3 * time.Second, //3초만 기다린다.
		StartOffset: kafka.FirstOffset, // GroupID 이전에 동일한 설정으로 데이터 사용한 적이
		// 있는 경우 중단한 곳부터 계속된다.
	})

	for j := i; j > 0; j-- {
		// ReadMessage 메서드는 우리가 다음 이벤트를 받을 때까지 차단된다.
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}

		s := strings.Split(string(msg.Value), "//")
		p.SendMail(s[0], s[1], s[2], s[3], s[4])
	}

	return nil
}

func (p *ProjectStart) SendMail(send string, receive string, name string, title string, content string) error {

	port, _ := strconv.Atoi(p.SmtpPort) //string -> int

	d := gomail.NewDialer(p.SmtpHost, port, p.SmtpId, p.SmtpPw)

	s, err := d.Dial()

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	// for _, r := range list {
	m.SetHeader("From", send)
	m.SetAddressHeader("To", receive, name)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", fmt.Sprintf(content))

	if err := gomail.Send(s, m); err != nil {
		return fmt.Errorf(
			"Could not send email to %q: %v ", send, err)
	}
	m.Reset()
	return nil
}

// 메일 내용 파싱함수
// 이 파싱 메서드가 성능이 어떤지는 아직 제대로 점검해보지 않았음..
func Parsing(str string, str1 string, str2 string, str3 string, str4 string) string {

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

	// todo 추후 도메인 주소가 나오면 파싱할 항목을 하나 더 늘린다.
	//if strings.Contains(str, "{{count_ip}}") {
	//	str = strings.Replace(str, "{{count_ip}}", str5, -1)
	//}

	return str
}