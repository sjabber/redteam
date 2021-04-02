package model

import (
	"database/sql"
)

type CounterModel struct {
	TargetNo        int  `json:"target_no"`
	ProjectNo       int  `json:"project_no"`
	EmailReadStatus bool `json:"email_read_status"`
	LinkClickStatus bool `json:"link_click_status"`
	DownloadStatus  bool `json:"download_status"`
}

func (cm *CounterModel) UpdateCount(conn *sql.DB) error {
	// 1차적으로 프로젝트의 날짜가 종료시점 이전인지 검사한다.
	endTimeQuery := `
select p_end_date
from project_info
where p_no = $1
  and p_end_date > now()
`
	end := ""
	err := conn.QueryRow(endTimeQuery, cm.ProjectNo).Scan(&end)
	if err != nil {
		// 카운트용 api이기 때문에 굳이 error를 표출하지 않음
		SugarLogger.Error(err.Error())
		return nil
	}

	// 이후에 count_info 테이블에 해당하는 값이 존재하는지 조회하고
	existQuery := `
select target_no,
       project_no,
       email_read_status,
       link_click_status,
       download_status
from count_info ci
         left join project_info pi on ci.project_no = pi.p_no
where ci.project_no = $1
  and ci.target_no = $2
	`
	resultCM := CounterModel{}
	exist := conn.QueryRow(existQuery, cm.ProjectNo, cm.TargetNo)

	// 존재하지 않으면 값을 삽입
	err = exist.Scan(&resultCM.TargetNo, &resultCM.ProjectNo, &resultCM.EmailReadStatus, &resultCM.LinkClickStatus, &resultCM.DownloadStatus)
	if err == sql.ErrNoRows {
		_, err := conn.Exec(`
insert into count_info (target_no,
						project_no,
						email_read_status,
						link_click_status,
						download_status)
values ($1, $2, $3, $4, $5)
`, cm.TargetNo, cm.ProjectNo, cm.EmailReadStatus, cm.LinkClickStatus, cm.DownloadStatus)
		if err != nil {
			SugarLogger.Error(err.Error())
			return err
		}
	} else { // 존재하면 값을 갱신해준다. resultCM은 현재 DB에 있는 값을 가지고, cm은 Get 메소드로부터 받은 값을 가지고있는다.
		// 한번 true 된 값은 변하지 않는다.
		if resultCM.EmailReadStatus == true {
			cm.EmailReadStatus = true
		}
		if resultCM.DownloadStatus == true {
			cm.DownloadStatus = true
		}
		if resultCM.LinkClickStatus == true {
			cm.LinkClickStatus = true
		}
		_, err = conn.Exec(`
update count_info
set email_read_status = $1,
    link_click_status = $2,
    download_status = $3,
    modified_time = now()
where target_no = $4
  and project_no = $5
`, cm.EmailReadStatus, cm.LinkClickStatus, cm.DownloadStatus, cm.TargetNo, cm.ProjectNo)
	}
	if err != nil {
		SugarLogger.Error(err.Error())
		return err
	}

	defer conn.Close()

	return nil
}
