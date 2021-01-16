package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	PNo          int      `json:"p_no"`          // 진짜 프로젝트 번호
	TmlNo        int      `json:"tml_no"`        // 템플릿 번호
	FakeNo       int      `json:"fake_no"`       // 화면에 표시될 프로젝트의 번호
	PName        string   `json:"p_name"`        // 프로젝트 이름
	PDescription string   `json:"p_description"` // 프로젝트 설명
	TagArray     []string `json:"tag_no"`        // 등록할 태그 대상자들
	PStatus      string   `json:"p_status"`      // 프로젝트 진행행태
	TemplateNo   string   `json:"tmp_no"`        // 적용할 템플릿 번호나 이름
	Infection    string   `json:"infection"`     // 감염비율
	SendNo       int      `json:"send_no"`       // 메일 보낸 횟수
	Targets      int      `json:"targets"`       // 훈련 대상자수
	StartDate    string   `json:"start_date"`    // 프로젝트 시작일
	EndDate      string   `json:"end_date"`      // 프로젝트 종료일
}

type ProjectStart struct {
	PNo            int    `json:"p_no"`
	UserNo         int    `json:"user_no"`
	TargetNo       int    `json:"target_no"`       // 훈련대상자들 번호
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

type ProjectNumber struct {
	ProjectNumber []string `json:"project_list"` //front javascript 와 이름을 일치시켜야함.
}

const (
	topic = "redteam"
	brokerAddress = "localhost:9092"
	partition = 1
)
var Msg string

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
       				p_no,
       				tmp_no,
       				tmp_name,
       				p_name,
       				p_status,
       				to_char(p_start_date, 'YYYY-MM-DD'),
       				to_char(p_end_date, 'YYYY-MM-DD'),
				    T.tag1,
				    T.tag2,
				    T.tag3,
				    T.send_no,
				    COUNT(ta.target_no),
       				COUNT(ci.target_no)
			FROM (SELECT ROW_NUMBER() over (ORDER BY p_no) AS row_num,
					 p_no,
			         tmp_no,
					 ti.tmp_name,
					 p_name,
					 p_status,
					 p_start_date,
					 p_end_date,
					 p.tag1,
					 p.tag2,
					 p.tag3,
					 p.send_no
				FROM project_info as p
					   LEFT JOIN template_info ti on p.tml_no = ti.tmp_no
			  	WHERE p.user_no = $1
			) AS T
				 LEFT JOIN target_info ta on user_no = ta.user_no
				 LEFT JOIN count_info ci on ta.target_no = ci.target_no AND T.p_no = ci.project_no
		WHERE T.tag1 > 0 and (T.tag1 = ta.tag1 or T.tag1 = ta.tag2 or T.tag1 = ta.tag3)
			or T.tag2 > 0 and (T.tag2 = ta.tag1 or T.tag2 = ta.tag2 or T.tag2 = ta.tag3)
			or T.tag3 > 0 and (T.tag3 = ta.tag1 or T.tag3 = ta.tag2 or T.tag3 = ta.tag3)
		GROUP BY row_num, p_no, tmp_no, tmp_name, p_name, p_status, to_char(p_start_date, 'YYYY-MM-DD'),
				 to_char(p_end_date, 'YYYY-MM-DD'), T.tag1, T.tag2, T.tag3, T.send_no
		ORDER BY row_num;`

	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading the projects. ")
	}

	var tags [3]int
	var projects []Project // Project 구조체를 값으로 가지는 배열
	for rows.Next() {
		p := Project{}
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TmlNo, &p.TemplateNo, &p.PName, &p.PStatus,
			&p.StartDate, &p.EndDate, &tags[0], &tags[1], &tags[2], &p.SendNo, &p.Targets, &p.Infection)
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

func (p *ProjectNumber) DeleteProject(conn *sql.DB, num int) error {

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

	return nil
}

// 사용자가 시작버튼을 누른경우에만 동작한다.
func (p *ProjectStart) StartProject(conn *sql.DB, num int) error {
	// 프로젝트 상태를 1로 변경하며 프로젝트를 실행한다.
/*	_, err := conn.Exec(`UPDATE project_info
 								SET p_status = 1
 								WHERE user_no = $1 AND p_no = $2`,
		num, p.PNo)*/
	_, err := conn.Exec(`UPDATE project_info
 								SET p_start_date = now()
 								WHERE user_no = $1 AND p_no = $2`,
		num, p.PNo)
	if err != nil {
		return fmt.Errorf("Error : updating project status. ")
	}

	query := `SELECT distinct mail_title,
               mail_content,
               sender_name,
               target_name,
               target_email,
               target_organize,
               target_position,
               target_phone,
               target_no,
               p_no
			FROM (SELECT p_no, tml_no, tag1, tag2, tag3, user_no
     			FROM project_info
     			WHERE p_no = $1
       			AND user_no = $2) as T
        			LEFT JOIN template_info tmp on tml_no = tmp.tmp_no
        			LEFT JOIN target_info ta on T.user_no = ta.user_no
			WHERE T.tag1 > 0 and (T.tag1 = ta.tag1 or T.tag1 = ta.tag2 or T.tag1 = ta.tag3)
   				or T.tag2 > 0 and (T.tag2 = ta.tag1 or T.tag2 = ta.tag2 or T.tag2 = ta.tag3)
   				or T.tag3 > 0 and (T.tag3 = ta.tag1 or T.tag3 = ta.tag2 or T.tag3 = ta.tag3)
   				GROUP BY mail_title, mail_content, sender_name, target_name, target_email, target_organize,
   				         target_position, target_phone, target_no, p_no
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

	msg := ProjectStart{}
	msg.UserNo = num

	// todo Kafka producer
	go func() {
		for rows.Next() {
			// DB 로부터 토픽에 작성할 내용들을 불러온다.
			err = rows.Scan(&msg.MailTitle, &msg.MailContent, &msg.SenderEmail,
				&msg.TargetName, &msg.TargetEmail, &msg.TargetOrganize,
				&msg.TargetPosition, &msg.TargetPhone, &msg.TargetNo, &msg.PNo)
			if err != nil {
				fmt.Errorf("Error : sql error ")
			}

			// 파싱작업 수행
			msg.MailContent = parsing(msg.MailContent, msg.TargetName, msg.TargetOrganize,
				msg.TargetPosition, msg.TargetPhone, strconv.Itoa(msg.TargetNo), strconv.Itoa(msg.PNo))

			// 카프카에 작성할 내용들을 하나의 띄어쓰기로 구분짓고 하나의 string 으로 묶는다.
			// 각각 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용 순이다.
			message, _ := json.Marshal(msg)
			produce(message, w)
		}
	}()

	return nil
}

// todo Kafka producer function
func produce (messages []byte, w kafka.Writer) {
	err := w.WriteMessages(context.Background(), kafka.Message{
		//Key: []byte("Key"),
		Value: messages,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
}

// todo kafka consumer function
func (p *ProjectStart)  Consumer() {
	// todo Kafka consumer
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     "redteam",
		MinBytes:    5,                 //5 바이트
		MaxBytes:    4132,              //1kB
		MaxWait:     3 * time.Second,   //3초만 기다린다.
		StartOffset: kafka.LastOffset, // GroupID 이전에 동일한 설정으로 데이터 사용한 적이
		// 있는 경우 중단한 곳부터 계속된다.
	})

	for {
		// ReadMessage 메서드는 우리가 다음 이벤트를 받을 때까지 차단된다.
		// json 객체로 받은 값들을 sendMail 메서드에 잘 적재시킨다.
		Msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("Could not read message " + err.Error())
		}

		// json -> ProjectStart 객체 각각에 값으로 들어가도록 한다.
		json.Unmarshal(Msg.Value, p)

		// 토픽의 값이 비어있지 않다면 값을 읽어 메일을 전송한다.
		if p.SenderEmail != "" || p.TargetEmail != "" || p.MailTitle != "" || p.TargetName != ""{
			err = p.sendMail(p.UserNo)
			if err != nil {
				panic("Could not send email " + err.Error())
			}
		} else {
			continue
		}
	}
}

func (p *ProjectStart) sendMail(num int) error {
	// 부득이하게 다른 DB connecting 방법 사용..
	conn, err := ConnectDB()
	if err != nil {
		return fmt.Errorf("db connection error")
	}
	defer conn.Close()

	// 메일을 보내기 위한 정보를 DB로 부터 가져온다.
	row := conn.QueryRow(`SELECT smtp_host, smtp_port, smtp_id, smtp_pw
 								FROM smtp_info
 								WHERE user_no = $1`, num)
	err = row.Scan(&p.SmtpHost, &p.SmtpPort, &p.SmtpId, &p.SmtpPw)
	if err != nil {
		panic("failed to write message: " + err.Error())
	}

	// string -> int, smtp 연결을 테스트한다.
	port, _ := strconv.Atoi(p.SmtpPort) //string -> int
	d := gomail.NewDialer(p.SmtpHost, port, p.SmtpId, p.SmtpPw)
	s, err := d.Dial()
	if err != nil {
		panic("failed to write message: " + err.Error())
	}

	m := gomail.NewMessage()
	// for _, r := range list {
	m.SetHeader("From", p.SenderEmail)
	m.SetAddressHeader("To", p.TargetEmail, p.TargetName)
	m.SetHeader("Subject", p.MailTitle)
	m.SetBody("text/html", fmt.Sprintf(p.MailContent))

	//send close
	if err := gomail.Send(s, m); err != nil {
		return fmt.Errorf(
			"Could not send email to %q: %v ", p.SenderEmail, err)
	}

	_, err = conn.Exec(`UPDATE project_info 
	SET send_no = project_info.send_no + 1
	WHERE	user_no = $1 AND p_no = $2`, num, p.PNo)
	if err != nil {
		return fmt.Errorf(
			"Could not update information : %v ", err)
	}

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

	// todo 추후 도메인 주소가 나오면 파싱할 항목을 하나 더 늘린다.
	if strings.Contains(str, "{{count_ip}}") {
		s := "<html>\n<body>\n<img src=\"http://localhost:5000/api/CountTarget?" +
			"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\">\n" +
			"<a href=\"http://localhost:5000/api/CountTarget?" +
			"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\"></a>\n</body>\n</html>"

		str = strings.Replace(str, "{{count_ip}}", s, -1)
	}

	return str
}