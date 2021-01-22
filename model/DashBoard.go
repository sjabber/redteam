package model

import (
	"database/sql"
	"fmt"
)

func GetDashboardInfo1(conn *sql.DB, num int) ([]Project ,error) {
	// 프로젝트 읽어오기전에 해시테이블에 태그정보 한번 넣고 시작한다.
	var query string

	query = `SELECT row_num,
       				p_no,
       				tmp_no,
       				tmp_name,
       				p_name,
       				to_char(p_start_date, 'YYYY-MM-DD'),
       				to_char(p_end_date, 'YYYY-MM-DD'),
				    T.tag1,
				    T.tag2,
				    T.tag3,
				    T.send_no,
				    COUNT(ta.target_no),
       				COUNT(ci.target_no),
       				COUNT(CASE WHEN ci.link_click_status THEN 1 END)
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
		GROUP BY row_num, p_no, tmp_no, tmp_name, p_name, to_char(p_start_date, 'YYYY-MM-DD'),
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
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TmlNo, &p.TemplateNo, &p.PName,
			&p.StartDate, &p.EndDate, &tags[0], &tags[1], &tags[2], &p.SendNo, &p.Targets, &p.Reading, &p.Infection)
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
