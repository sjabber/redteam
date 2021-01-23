package model

import (
	"database/sql"
)

type PInfo1 struct {
	Targets   int `json:"targets"`
	Scheduled int `json:"scheduled"`
	Ongoing   int `json:"ongoing"`
	Closed    int `json:"closed"`
}



func GetDashboardInfo1(conn *sql.DB, num int) (PInfo1, error) {
	// 프로젝트 읽어오기전에 해시테이블에 태그정보 한번 넣고 시작한다.
	var query string

	pi := PInfo1{}

	query = `SELECT count(DISTINCT target_no), ready, ing, closed
			 FROM target_info as ti
         				LEFT JOIN project_info pi on pi.user_no = ti.user_no
         				LEFT JOIN (SELECT ready, ing, closed, user_no
                    				FROM (SELECT COUNT(case when p_status = 0 then 1 end) as ready,
                                 				COUNT(case when p_status = 1 then 1 end) as ing,
                                 				COUNT(case when p_status = 2 then 1 end) as closed,
                                 				user_no
                          				  FROM project_info
                    					  WHERE user_no = $1
                          				  GROUP BY user_no) as pi1) pi2 on pi2.user_no = ti.user_no
			 WHERE ti.user_no = $1 and ti.tag1 > 0 and (ti.tag1 = pi.tag1 or ti.tag1 = pi.tag2 or ti.tag1 = pi.tag3)
   				or ti.tag2 > 0 and (ti.tag2 = pi.tag1 or ti.tag2 = pi.tag2 or ti.tag2 = pi.tag3)
   				or ti.tag3 > 0 and (ti.tag3 = pi.tag1 or ti.tag3 = pi.tag2 or ti.tag3 = pi.tag3)
			 GROUP BY ready, ing, closed;`
	rows := conn.QueryRow(query, num)
	err := rows.Scan(&pi.Targets, &pi.Scheduled, &pi.Ongoing, &pi.Closed)
	if err != nil {
		return pi, err
	}

	return pi, nil
}

func GetDashboardInfo2(conn *sql.DB, num int) () {


}