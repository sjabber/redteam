package model

import (
	"database/sql"
	"fmt"
	"regexp"
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
       to_char(p_start_date, 'YYYY-MM-DD HH24:MI'),
       to_char(p_end_date, 'YYYY-MM-DD HH24:MI'),
       T.tag1,
       T.tag2,
       T.tag3,
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
             p.tag3
      FROM project_info as p
               LEFT JOIN template_info ti on p.tml_no = ti.tmp_no
      WHERE p.user_no = $1
     ) AS T
         LEFT JOIN target_info ta on user_no = ta.user_no
WHERE T.tag1 > 0 and (T.tag1 = ta.tag1 or T.tag1 = ta.tag2 or T.tag1 = ta.tag3)
   or T.tag2 > 0 and (T.tag2 = ta.tag1 or T.tag2 = ta.tag2 or T.tag2 = ta.tag3)
   or T.tag3 > 0 and (T.tag3 = ta.tag1 or T.tag3 = ta.tag2 or T.tag3 = ta.tag3)
GROUP BY row_num, p_no, tmp_name, p_name, p_status, to_char(p_start_date, 'YYYY-MM-DD HH24:MI'),
         to_char(p_end_date, 'YYYY-MM-DD HH24:MI'), T.tag1, T.tag2, T.tag3
ORDER BY row_num;`

	rows, err := conn.Query(query, num)
	if err != nil {
		return nil, fmt.Errorf("There was an error reading the projects. ")
	}

	var tags [3]int
	var projects []Project // Project 구조체를 값으로 가지는 배열
	for rows.Next() {
		p := Project{}
		err = rows.Scan(&p.FakeNo, &p.PNo, &p.TemplateNo, &p.PName, &p.PStatus, &p.StartDate, &p.EndDate, &tags[0], &tags[1], &tags[2], &p.Targets)
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

func (p *Project) StartProject(conn *sql.DB, num int) error {
	_, err := conn.Exec(`UPDATE project_info
 								SET p_status = 1
 								WHERE user_no = $1 AND p_no = $2`,
		num, p.PNo)
	if err != nil {
		return fmt.Errorf("Error updating project status ")
	}

	//query := ``


	return nil
}