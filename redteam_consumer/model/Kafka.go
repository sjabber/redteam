package model

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// todo 추후 환경변수로 변경
	key = "qlwkndlqiwndlian"
)

type User struct {
	UserNo        int    `json:"user_no"`
	Email         string `json:"email"`
	PasswordHash  string `json:"-"`
	Password      string `json:"password"`
	PasswordCheck string `json:"password_check"` //회원가입에서 사용된다.
	Name          string `json:"name"`
}

// 프로젝트 시작(Consumer)에서 사용하는 구조체
type ProjectStart struct {
	PNo            int    `json:"p_no"`
	TmpNo          int    `json:"tmp_no"`
	TargetNo       int    `json:"target_no"` // 훈련대상자들 번호
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

const (
	topic         = "redteam"
	brokerAddress = "localhost:9092"
)

// 프로젝트가 종료된 다음에 p_status 값 자동으로 변경시키는거 추가해야함.
// 자동으로 프로젝트 실행됐을 때 버튼눌렀을때와 동일하게 작동하도록도 수정해야함.
// 작동되는 시간 언제로 할지 논리적으로 따져서 결정하기.
//cron(스케쥴링)을 활용하여 매일 특정 시간에 프로젝트들의 진행상황을 체크한다.
func AutoStartProject() {
	wg := &sync.WaitGroup{}
	wg.Add(1) // WaitGroup 의 고루틴 개수 1개 증가

	kor, _ := time.LoadLocation("Asia/Seoul")
	c := cron.New(cron.WithLocation(kor))

	// 매 특정 시간마다 프로젝트들의 날짜를 점검하여 실행 종료시킨다.
	//c.AddFunc("0 0 * * *", Auto) // 매일 정각에 프로젝트를 검토후 진행&종료시킨다.
	c.AddFunc("09-10 14 * * *", Auto) // 매일 오후 2시 9-10분 마다 실행할 프로젝트를 검토한다.
	c.AddFunc("54-56 5 * * *", Auto) // 매일 오후 2시 9-10분 마다 실행할 프로젝트를 검토한다.
	c.AddFunc("58-59 23 * * *", Auto) //하루 끝에 최종적으로 한번더 검토한다.
	c.Start()
	wg.Wait()
}

func Auto() {
	// crontab(스케줄러)으로 설정한 시간이 되면 project_info 테이블에서
	// 프로젝트 번호, 템플릿 번호, 사용자번호, 상태, 시작일, 종료일을 조회하여
	// 조건을 충족하면 프로젝트를 자동실행시킨다.

	var Pno, TmlNo, num, status, startDate, endDate string

	// DB 연결
	conn, err := ConnectDB()
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
	defer conn.Close()

	query := `SELECT p_no, tml_no, user_no, p_status, 
       				 to_char(p_start_date, 'YYYY-MM-DD') as start_date, to_char(p_end_date, 'YYYY-MM-DD') as end_date
			  FROM project_info`
	rows, err := conn.Query(query)
	if err != nil {
		log.Println(err)
		//_ = fmt.Errorf("%v", err)
		//panic(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&Pno, &TmlNo, &num, &status, &startDate, &endDate)
		if err != nil {
			log.Println(err)
			//_ = fmt.Errorf("%v", err)
			//panic(err.Error())
		}

		// 날짜가 오늘 && 프로젝트가 예약상태인 경우 -> 프로젝트를 실행시킨다.
		if startDate == time.Now().Format("2006-01-02") && status == "0" {

			// 만약 템플릿이 삭제된 템플릿일 경우 프로젝트의 상태를 오류로 변경한다.
			if TmlNo == "0" {
				_, err = conn.Exec(`UPDATE project_info
 								SET p_start_date = now(), p_status = 3
 								WHERE user_no = $1 AND p_no = $2`, num, Pno)
				if err != nil {
					log.Println(err)
					//_ = fmt.Errorf("%v", err)
				}
				// 오류가 있으므로 다음 프로젝트로 넘어간다.
				continue
			} else {
				// 프로젝트번호, 템플릿번호, 훈련 대상자번호, 사용자번호만 조회해서 TOPIC 에 적재한다.
				query2 := `SELECT p_no, tml_no, target_no, user_no
				FROM project_target_info
				WHERE p_no = $1 and user_no = $2
				ORDER BY target_no;`

				rows2, err := conn.Query(query2, Pno, num)
				if err != nil {
					log.Println(err)
					//_ = fmt.Errorf("%v", err)
					//panic(err.Error())
				}

				// 카프카에 넣을 메일 내용
				// 보내는사람 이메일, 받는사람 이메일, 받는사람 이름, 메일제목, 메일 내용
				w := kafka.Writer{
					Addr:  kafka.TCP(brokerAddress),
					Topic: topic,
				}

				msg := ProjectStart2{}

				// Kafka producer
				go func() {
					for rows2.Next() {
						// DB 로부터 토픽에 작성할 내용들을 불러온다.
						err = rows2.Scan(&msg.PNo, &msg.TmpNo, &msg.TargetNo, &msg.UserNo)
						if err != nil {
							log.Println(err)
							//_ = fmt.Errorf("%v", err)
							//panic(err.Error())
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
					//panic(err.Error())
				}
			}

			// 프로젝트가 종료일이면 종료한다.
		} else if endDate == time.Now().Format("2006-01-02") && (status == "0" || status == "1") {
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

// Kafka producer function
func produce(messages []byte, w kafka.Writer) {
	err := w.WriteMessages(context.Background(), kafka.Message{
		//Key: []byte("Key"),
		Value: messages,
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
}

// kafka consumer function
func (p *ProjectStart) Consumer() {
	// Kafka consumer
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     "RedTeam",
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
       								smtp_id as sender_email,
       								smtp_host,
       								smtp_port,
       								smtp_id,
       								smtp_pw
								FROM (SELECT target_name, target_email, target_organize, target_position, target_phone
    								  FROM target_info
    								  WHERE target_no = $1 AND user_no = $2) as T
        									LEFT JOIN template_info as tmp on tmp_no = $3
        									LEFT JOIN smtp_info si on si.user_no = $2;`, p.TargetNo, p.UserNo, p.TmpNo)

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

	time.Sleep(8 * time.Second)
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
		s := "<html>\n<body>\n<img width=0" + " height=0" + " src=\"http://localhost:5000/api/CountTarget?" +
			"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\">\n" +
			//"<a href=\"http://localhost:5000/api/CountTarget?" +
			//"tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=false&download=false\"></a>" +
			"\n</body>\n</html>"

		str = strings.Replace(str, "{{count_ip}}", s, -1)
	}

	// 접속 사이트에 링크를 첨부한다.
	if strings.Contains(str, "{{link_ip}}") {
		s := "tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=true&download=false"
		str = strings.Replace(str, "{{link_ip}}", s, -1)
	}

	//if strings.Contains(str, "{{download_ip}}") {
	//	s := "tNo=" + str5 + "&pNo=" + str6 + "&email=true&link=true&download=false"
	//	str = strings.Replace(str, "{{link_ip}}", s, -1)
	//}


	// 감염 사이트 링크는 자바스크립트에서 처리한다.

	return str
}

// 복호화
func Decrypt(b cipher.Block, ciphertext []byte) []byte {

	if len(ciphertext)%aes.BlockSize != 0 {
		// 암호화된 데이터의 길이기 블록크기의 배수가 아니면 작동이 안됨.
		fmt.Println("The length of decrypted data must be a multiple of the block size. ")
		return nil
	}
	// todo 추후 환경변수로 변경
	iv := []byte("0987654321654321")

	plaintext := make([]byte, len(ciphertext)) // 평문 데이터를 저장할 공간을 생성한다.
	mode := cipher.NewCBCDecrypter(b, iv) // 암호화 블록과 초기화 벡터를 넣어 복호화된 블록모드의 인스턴스를 생성한다.

	mode.CryptBlocks(plaintext, ciphertext) //

	padding := plaintext[len(plaintext)-1] //가장 마지막 값(패딩값)을 가져온다.
	plaintext = plaintext[:len(plaintext)-int(padding)] // 패딩값을 빼준다.

	return plaintext
}