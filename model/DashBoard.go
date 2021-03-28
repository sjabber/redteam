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
	TagArray    [3]string `json:"tag_no"`
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
		SugarLogger.Error(err.Error())
		return PInfo1{}, fmt.Errorf("%v", err)
	}

	defer conn.Close()

	return pi1, nil
}

func GetDashboardInfo2(conn *sql.DB, num int, PNum int) (PInfo2, error) {

	var query string

	// 객체생성
	pi2 := PInfo2{}

	query = `SELECT ti.tmp_name,
       ti.mail_title,
       sender_email,
       p_name,
       to_char(p_start_date, 'YYYY-MM-DD'),
       to_char(p_end_date, 'YYYY-MM-DD'),
       COALESCE(tag_name1, '')                          as tag_name1,
       COALESCE(tag_name2, '')                          as tag_name2,
       COALESCE(tag_name3, '')                          as tag_name3,
       target_count                                     as Targets,
       T.send_no,
       COUNT(distinct ci.target_no)                     as Read,
       COUNT(CASE WHEN ci.link_click_status THEN 1 END) as Connect,
       COUNT(CASE WHEN ci.download_status THEN 1 END)   as Infection
FROM (SELECT target_count,
             p.p_no,
             sender_email,
             p_name,
             tml_no,
             p_start_date,
             p_end_date,
             tag1,
             tag2,
             tag3,
             send_no,
             p.user_no
      FROM project_info as p
               LEFT JOIN (SELECT COUNT(target_no) AS target_count, user_no, p_no
                          FROM project_target_info
                          GROUP BY user_no, p_no) AS pti ON pti.user_no = p.user_no AND pti.p_no = p.p_no
      WHERE p.user_no = $1
        AND p.p_no = $2) AS T
         LEFT JOIN (SELECT tag_name as tag_name1, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti1 on ti1.tag_no = T.tag1
         LEFT JOIN (SELECT tag_name as tag_name2, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti2 on ti2.tag_no = T.tag2
         LEFT JOIN (SELECT tag_name as tag_name3, user_no, tag_no
                    FROM tag_info
                    WHERE user_no = $1) ti3 on ti3.tag_no = T.tag3
         LEFT JOIN template_info ti on T.tml_no = ti.tmp_no
         LEFT JOIN target_info ta on T.user_no = ta.user_no
         LEFT JOIN count_info ci on ta.target_no = ci.target_no AND T.p_no = ci.project_no
GROUP BY ti.tmp_name, ti.mail_title, sender_email, p_name, to_char(p_start_date, 'YYYY-MM-DD'),
         to_char(p_end_date, 'YYYY-MM-DD'), COALESCE(tag_name1, ''), COALESCE(tag_name2, ''), COALESCE(tag_name3, ''),
         target_count, T.send_no;`

	row := conn.QueryRow(query, num, PNum)

	// 프로젝트 상세에 필요한 정보들을 바인딩
	err := row.Scan(&pi2.TmpName, &pi2.MailTitle, &pi2.SenderEmail, &pi2.PName, &pi2.StartDate, &pi2.EndDate,
					&pi2.TagArray[0], &pi2.TagArray[1], &pi2.TagArray[2], &pi2.Targets, &pi2.SendNo, &pi2.Reading,
					&pi2.Connect, &pi2.Infection)
	if err != nil {
		SugarLogger.Error(err.Error())
		return PInfo2{}, fmt.Errorf("%v", err)
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
		SugarLogger.Error(err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	var project []PInfo3
	pi3 := PInfo3{}
	for rows.Next() {
		err = rows.Scan(&pi3.PNo, &pi3.FakeNo, &pi3.PName, &pi3.PStatus, &pi3.EndDate)
		if err != nil {
			SugarLogger.Error(err.Error())
			return nil, fmt.Errorf("%v", err)
		}

		project = append(project, pi3)
	}

	defer conn.Close()

	return project, nil
}