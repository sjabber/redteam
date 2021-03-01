package model

import (
	"database/sql"
	"fmt"
)

type PInfo1 struct {
	Targets   int `json:"targets"`
	Scheduled int `json:"scheduled"`
	Ongoing   int `json:"ongoing"`
	Closed    int `json:"closed"`
}

type PInfo2 struct {
	PName		string	 `json:"p_name"`
	TmpName     string   `json:"tmp_name"`
	MailTitle   string   `json:"mail_title"`
	SenderEmail string   `json:"sender_email"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	TagArray    []string `json:"tag_no"`
	Targets     string   `json:"targets"`
	SendNo      string   `json:"send_no"`
	Reading     string   `json:"reading"`
	Connect     string   `json:"connect_no"`
	Infection   string   `json:"infection"`
}

// p_no ,p_name , p_status, end_date
type PInfo3 struct {
	PNo			int		`json:"p_no"`
	FakeNo		int		`json:"fake_no"`
	PName		string	`json:"p_name"`
	PStatus		string	`json:"p_status"`
	EndDate		string	`json:"end_date"`
}


func GetDashboardInfo1(conn *sql.DB, num int) (PInfo1, error) {

	var query string

	pi1 := PInfo1{}

	query = `SELECT count(distinct pti.target_no), ready, ing, closed
			 FROM project_target_info as pti
					  LEFT JOIN (SELECT ready, ing, closed, user_no
								 FROM (SELECT COUNT(case when p_status = 0 then 1 end) as ready,
											 COUNT(case when p_status = 1 then 1 end) as ing,
											 COUNT(case when p_status = 2 then 1 end) as closed,
											 user_no
							 		   FROM project_info
									   WHERE user_no = $1
									   GROUP BY user_no) as pi1) pi2 on pi2.user_no = pti.user_no
			 WHERE pti.user_no = $1
			 GROUP BY ready, ing, closed;`
	rows := conn.QueryRow(query, num)
	err := rows.Scan(&pi1.Targets, &pi1.Scheduled, &pi1.Ongoing, &pi1.Closed)
	if err == sql.ErrNoRows { //프로젝트가 존재하지 않는 경우
		pi1.Targets, pi1.Scheduled, pi1.Ongoing, pi1.Closed = 0, 0, 0 ,0
		return pi1, nil
	} else if err != nil {
		return PInfo1{}, err
	}

	defer conn.Close()

	return pi1, nil
}

func GetDashboardInfo2(conn *sql.DB, num int, pnum int) (PInfo2, error) {

	var query string

	query = `SELECT tag_no, tag_name
			  FROM tag_info
			  WHERE user_no = $1
			  ORDER BY tag_no asc`
	hash, err := conn.Query(query, num)
	if err != nil {
		fmt.Println(err)
		return PInfo2{}, err
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

	pi2 := PInfo2{}

	var tags [3]int

	query = `SELECT T.tmp_name,
				    mail_title,
				    sender_name,
       			    p_name,
				    to_char(p_start_date, 'YYYY-MM-DD'),
				    to_char(p_end_date, 'YYYY-MM-DD'),
				    T.tag1,
				    T.tag2,
				    T.tag3,
				    COUNT(distinct pti.target_no) as Targets,
				    T.send_no,
				    COUNT(distinct ci.target_no) as Read,
				    COUNT(CASE WHEN ci.link_click_status THEN 1 END) as Connect,
				    COUNT(CASE WHEN ci.download_status THEN 1 END) as Infection
			FROM (SELECT p_no,
			  			 p_name,
						 tmp_name,
						 mail_title,
						 sender_name,
						 p_start_date,
						 p_end_date,
						 tag1,
						 tag2,
						 tag3,
						 send_no,
						 p.user_no
				  FROM project_info as p
						   LEFT JOIN template_info ti on p.tml_no = ti.tmp_no
		  WHERE p.user_no = $1  AND p.p_no = $2) AS T
			 LEFT JOIN project_target_info pti on T.user_no = pti.user_no AND T.p_no = pti.p_no
			 LEFT JOIN target_info ta on T.user_no = ta.user_no
			 LEFT JOIN count_info ci on ta.target_no = ci.target_no AND T.p_no = ci.project_no
		  GROUP BY T.tmp_name, mail_title, sender_name, p_name, p_start_date, p_end_date, T.tag1, T.tag2, T.tag3, T.send_no;`

	row := conn.QueryRow(query, num, pnum)

	// 프로젝트 상세에 필요한 정보들을 바인딩
	err = row.Scan(&pi2.TmpName, &pi2.MailTitle, &pi2.SenderEmail, &pi2.PName, &pi2.StartDate, &pi2.EndDate,
		&tags[0], &tags[1], &tags[2], &pi2.Targets, &pi2.SendNo, &pi2.Reading, &pi2.Connect, &pi2.Infection)
	if err != nil {
		return PInfo2{}, fmt.Errorf("%v", err)
	}

	Loop1:
		for i := 0; i < len(tags); i++ {
			if tags[i] == 0 {
				pi2.TagArray = append(pi2.TagArray, "")
				continue Loop1
			}

			for key, val := range Hashmap {
				if key == tags[i] {
					pi2.TagArray = append(pi2.TagArray, val)
					break
				}
			}
		}

	defer conn.Close()

	return pi2, nil
}

func GetDashboardInfo3(conn *sql.DB, num int) ([]PInfo3, error) {

	query := `SELECT p_no, ROW_NUMBER() over (ORDER BY p_no) as row_num, p_name, p_status, to_char(p_end_date, 'YYYY-MM-DD')
			  FROM project_info
			  WHERE user_no = $1`
	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var project []PInfo3
	pi3 := PInfo3{}
	for rows.Next() {
		err = rows.Scan(&pi3.PNo, &pi3.FakeNo, &pi3.PName, &pi3.PStatus, &pi3.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}

		project = append(project, pi3)
	}

	defer conn.Close()

	return project, nil
}