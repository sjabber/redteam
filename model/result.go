package model

import (
	"database/sql"
	"fmt"
	"strconv"
)

type ResultDetail struct {
	ProjectInfo        ProjectInfo
	ProjectSummary     ProjectSummary
	ResultPerStatistic ResultPerStatistic
}

type ProjectInfo struct {
	PNo          int    `json:"p_no"`          // 진짜 프로젝트 번호
	PName        string `json:"p_name"`        // 프로젝트 이름
	PDescription string `json:"p_description"` // 프로젝트 설명

	PStatusPercent  float32 `json:"p_status_percent"` // 정보 유출률 표기 (link 이후 / 총 훈련자 수)
	ProgressPercent float32 `json:"progress_percent"` // 전송률 (훈련 메일 전송 기반)

	StartDate string `json:"start_date"` // 프로젝트 시작일
	EndDate   string `json:"end_date"`   // 프로젝트 종료일

	SmtpId    string `json:"smtp_id"`     // 메일 전송 계정
	PMakeUser string `json:"p_make_user"` // 프로젝트 생성자

	TmlNo   int    `json:"tml_no"`   // 템플릿 번호
	TmlName string `json:"tml_name"` // 템플릿 이름

	TagArray []TagInfo `json:"tag_array"` // 프로젝트 등록 태그

	TargetCount   int `json:"target_count"`    // 훈련 대상자 수
	SendCount     int `json:"send_count"`      // 메일 보낸 횟수
	SendFailCount int `json:"send_fail_count"` // 발송 실패
	EmailRead     int `json:"email_read"`      //읽은 사람
	LinkClick     int `json:"link_click"`      // 링크 클릭 수
	Download      int `json:"download"`        // 감염비율

	PCreatedTime  string `json:"p_created_time"`  // 프로젝트 생성일
	PModifiedTime string `json:"p_modified_time"` // 프로젝트 최종 수정일
}

type TagInfo struct {
	TagName string `json:"tag_name"`
	TagNo   int    `json:"tag_no"`
}

type ProjectSummary struct {
	ResultOfAll             ResultOfAll
	ResultPerClassification ResultPerClassification
}

type ResultOfAll struct {
	InfectionCount int     `json:"infection_count"`
	TotalCount     int     `json:"total_count"`
	Rate           float32 `json:"rate"`
}

type ResultPerClassification struct {
	// C -> Classification
	CTag      ClassificationTag
	COrgan    ClassificationOrgan
	CPosition ClassificationPosition
}

type ClassificationTag struct {
	TagName        string  `json:"tag_name"`
	InfectionCount int     `json:"infection_count"`
	TotalCount     int     `json:"total_count"`
	Rate           float32 `json:"rate"`
}

type ClassificationOrgan struct {
	OrganName      string  `json:"organ_name"`
	InfectionCount int     `json:"infection_count"`
	TotalCount     int     `json:"total_count"`
	Rate           float32 `json:"rate"`
}

type ClassificationPosition struct {
	PositionName   string  `json:"position_name"`
	InfectionCount int     `json:"infection_count"`
	TotalCount     int     `json:"total_count"`
	Rate           float32 `json:"rate"`
}

type ResultPerStatistic struct {
	PerTag      []ClassificationTag
	PerOrgan    []ClassificationOrgan
	PerPosition []ClassificationPosition
}

func GetResultDetail(no string, userNo int, conn *sql.DB) (ResultDetail, error) {
	pNo, err := strconv.Atoi(no)
	if err != nil {
		return ResultDetail{}, err
	}
	rd := ResultDetail{
		ProjectInfo: ProjectInfo{PNo: pNo},
	}
	err = rd.getProjectInfo(userNo, conn)
	if err != nil {
		SugarLogger.Error(err.Error())
		return rd, err
	}
	err = rd.getProjectSummary(userNo, conn)
	if err != nil {
		SugarLogger.Error(err.Error())
		return rd, err
	}
	err = rd.getResultPerStatistic(userNo, conn)
	if err != nil {
		SugarLogger.Error(err.Error())
		return rd, err
	}

	defer conn.Close()

	return rd, nil
}

func (rd *ResultDetail) getProjectInfo(userNo int, conn *sql.DB) error {
	queryPi := `
select p_name
     , p_description
     , to_char(p_start_date, 'YYYY-MM-DD')             as p_start_date
     , to_char(p_end_date, 'YYYY-MM-DD')               as p_end_date
     , si.smtp_id
     , ui.user_email
     , ti.tmp_no
     , ti.tmp_name
     , to_char(pi.created_time, 'YYYY-MM-DD HH24:MI')  as created_time
     , to_char(pi.modified_time, 'YYYY-MM-DD HH24:MI') as modified_time
     , pi.send_no
     , t1.tag_name
     , t1.tag_no
     , t2.tag_name
     , t2.tag_no
     , t3.tag_name
     , t3.tag_no
from public.project_info pi
         left join smtp_info si on pi.user_no = si.user_no
         left join user_info ui on pi.user_no = ui.user_no
         left join template_info ti on pi.tml_no = ti.tmp_no
         left join tag_info t1 on t1.tag_no = pi.tag1
         left join tag_info t2 on pi.tag2 = t2.tag_no
         left join tag_info t3 on pi.tag3 = t3.tag_no
where pi.p_no = $1
  and pi.user_no = $2
`
	resultPi := conn.QueryRow(queryPi, rd.ProjectInfo.PNo, userNo)

	// todo ㅜㅜ 후지다 .... 설계가 잘못된듯 ..
	tag1 := TagInfo{}
	tag2 := TagInfo{}
	tag3 := TagInfo{}
	err := resultPi.Scan(
		&rd.ProjectInfo.PName,
		&rd.ProjectInfo.PDescription,
		&rd.ProjectInfo.StartDate,
		&rd.ProjectInfo.EndDate,
		&rd.ProjectInfo.SmtpId,
		&rd.ProjectInfo.PMakeUser,
		&rd.ProjectInfo.TmlNo,
		&rd.ProjectInfo.TmlName,
		&rd.ProjectInfo.PCreatedTime,
		&rd.ProjectInfo.PModifiedTime,
		&rd.ProjectInfo.SendCount,
		&tag1.TagName,
		&tag1.TagNo,
		&tag2.TagName,
		&tag2.TagNo,
		&tag3.TagName,
		&tag3.TagNo,
	)
	if err != nil {
		return err
	}
	if tag1.TagNo != 0 {
		rd.ProjectInfo.TagArray = append(rd.ProjectInfo.TagArray, tag1)
	}
	if tag2.TagNo != 0 {
		rd.ProjectInfo.TagArray = append(rd.ProjectInfo.TagArray, tag2)
	}
	if tag3.TagNo != 0 {
		rd.ProjectInfo.TagArray = append(rd.ProjectInfo.TagArray, tag3)
	}
	queryCount := `
select count(email_read_status) as read
     , count(link_click_status) as click
     , count(download_status)   as download
from project_info pi
         left join count_info ci on pi.p_no = ci.project_no
where (ci.email_read_status = true
    or ci.link_click_status = true
    or ci.download_status = true)
  and pi.p_no = $1
  and pi.user_no = $2
group by p_no
`
	resultCi := conn.QueryRow(queryCount, rd.ProjectInfo.PNo, userNo)
	err = resultCi.Scan(
		&rd.ProjectInfo.EmailRead,
		&rd.ProjectInfo.LinkClick,
		&rd.ProjectInfo.Download)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		SugarLogger.Error(err.Error())
		return err
	}

	return nil
}

func (rd *ResultDetail) getProjectSummary(userNo int, conn *sql.DB) error {
	query := `
select distinct count(target_no) as targetTotal
from target_info ti
where (ti.tag1 in (%s)
   or ti.tag2 in (%s)
   or ti.tag3 in (%s))
 and user_no = $1
`
	inCondition := ""
	for idx, no := range rd.ProjectInfo.TagArray {
		i := strconv.Itoa(no.TagNo)
		if (idx + 1) != len(rd.ProjectInfo.TagArray) {
			inCondition += i + ","
		} else {
			inCondition += i
		}
	}
	q := fmt.Sprintf(query, inCondition, inCondition, inCondition)
	tt := conn.QueryRow(q, userNo)
	err := tt.Scan(&rd.ProjectInfo.TargetCount)
	if err != nil {
		SugarLogger.Error(err.Error())
		return err
	}
	// 전송 실패 카운트
	rd.ProjectInfo.SendFailCount = rd.ProjectInfo.TargetCount - rd.ProjectInfo.SendCount

	rd.ProjectInfo.ProgressPercent = float32(float64(rd.ProjectInfo.SendCount)/
		float64(rd.ProjectInfo.TargetCount)) * 100

	infectionQuery := `
select count(*)
from count_info
where project_no = $1
  and (link_click_status = true or download_status = true)
`
	iQ := conn.QueryRow(infectionQuery, rd.ProjectInfo.PNo)
	err = iQ.Scan(&rd.ProjectSummary.ResultOfAll.InfectionCount)
	if err == nil {
		rd.ProjectInfo.PStatusPercent = float32(
			float64(rd.ProjectSummary.ResultOfAll.InfectionCount)/
				float64(rd.ProjectInfo.TargetCount)) * 100
	}

	// 감염되었다는 판단을 link, download 부터
	rd.ProjectSummary.ResultOfAll.TotalCount = rd.ProjectInfo.TargetCount
	rd.ProjectSummary.ResultOfAll.Rate = rd.ProjectInfo.PStatusPercent

	perPositionQuery := `
select ti.target_position, count(*) as count
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and (link_click_status = true or download_status = true)
group by ti.target_position
order by count desc
limit 1
`
	pPQ := conn.QueryRow(perPositionQuery, rd.ProjectInfo.PNo)
	_ = pPQ.Scan(&rd.ProjectSummary.ResultPerClassification.CPosition.PositionName,
		&rd.ProjectSummary.ResultPerClassification.CPosition.InfectionCount)

	perPositionTotalQuery := `
select count(*)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and target_position = $2
group by target_position
`
	oPTQ := conn.QueryRow(perPositionTotalQuery, rd.ProjectInfo.PNo,
		rd.ProjectSummary.ResultPerClassification.CPosition.PositionName)
	err = oPTQ.Scan(&rd.ProjectSummary.ResultPerClassification.CPosition.TotalCount)
	if err == nil {
		rd.ProjectSummary.ResultPerClassification.CPosition.Rate = float32(float64(
			rd.ProjectSummary.ResultPerClassification.CPosition.InfectionCount)/
			float64(rd.ProjectSummary.ResultPerClassification.CPosition.TotalCount)) * 100
	}

	perOrganQuery := `
select ti.target_organize, count(*) as count
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and (link_click_status = true or download_status = true)
group by ti.target_organize
order by count desc
limit 1
`
	pOQ := conn.QueryRow(perOrganQuery, rd.ProjectInfo.PNo)
	err = pOQ.Scan(&rd.ProjectSummary.ResultPerClassification.COrgan.OrganName,
		&rd.ProjectSummary.ResultPerClassification.COrgan.InfectionCount)

	perOrganTotal := `
select count(*)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and target_organize = $2
group by target_organize
`
	pOT := conn.QueryRow(perOrganTotal, rd.ProjectInfo.PNo,
		rd.ProjectSummary.ResultPerClassification.COrgan.OrganName)
	err = pOT.Scan(&rd.ProjectSummary.ResultPerClassification.COrgan.TotalCount)

	if err == nil {
		rd.ProjectSummary.ResultPerClassification.COrgan.Rate = float32(float64(
			rd.ProjectSummary.ResultPerClassification.COrgan.InfectionCount)/
			float64(rd.ProjectSummary.ResultPerClassification.COrgan.TotalCount)) * 100
	}


	return nil
}

func (rd *ResultDetail) getResultPerStatistic(userNo int, conn *sql.DB) error {
	// 조직별 통계
	perOrganTotalQuery := `
select target_organize, count(*)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
group by target_organize
`
	pOTQ, err := conn.Query(perOrganTotalQuery, rd.ProjectInfo.PNo)
	if err != nil {
		SugarLogger.Error(err.Error())
		return err
	}
	for pOTQ.Next() {
		tmpOrgan := ClassificationOrgan{}
		_ = pOTQ.Scan(&tmpOrgan.OrganName, &tmpOrgan.TotalCount)
		rd.ResultPerStatistic.PerOrgan = append(rd.ResultPerStatistic.PerOrgan, tmpOrgan)
	}
	perOrganQuery := `
select ti.target_organize, count(*) as count
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and (link_click_status = true or download_status = true)
group by ti.target_organize
order by count desc
`
	pOQ, _ := conn.Query(perOrganQuery, rd.ProjectInfo.PNo)
	for pOQ.Next() {
		tmpOrgan := ClassificationOrgan{}
		_ = pOQ.Scan(&tmpOrgan.OrganName, &tmpOrgan.InfectionCount)
		for idx, organ := range rd.ResultPerStatistic.PerOrgan {
			if organ.OrganName == tmpOrgan.OrganName {
				rd.ResultPerStatistic.PerOrgan[idx].InfectionCount =
					tmpOrgan.InfectionCount
				rd.ResultPerStatistic.PerOrgan[idx].Rate = float32(
					float64(rd.ResultPerStatistic.PerOrgan[idx].InfectionCount)/
						float64(rd.ResultPerStatistic.PerOrgan[idx].TotalCount)) * 100
			}
		}
	}

	// 직급별 통계
	perPositionTotalQuery := `
select target_position, count(*)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
group by target_position
`
	pPTQ, err := conn.Query(perPositionTotalQuery, rd.ProjectInfo.PNo)
	if err != nil {
		SugarLogger.Error(err.Error())
		return err
	}
	for pPTQ.Next() {
		tmpPosition := ClassificationPosition{}
		_ = pPTQ.Scan(&tmpPosition.PositionName, &tmpPosition.TotalCount)
		rd.ResultPerStatistic.PerPosition = append(rd.ResultPerStatistic.PerPosition, tmpPosition)
	}
	perPositionQuery := `
select ti.target_position, count(*) as count
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where project_no = $1
  and (link_click_status = true or download_status = true)
group by ti.target_position
order by count desc
`
	pPQ, _ := conn.Query(perPositionQuery, rd.ProjectInfo.PNo)
	for pPQ.Next() {
		tmpPosition := ClassificationPosition{}
		_ = pPQ.Scan(&tmpPosition.PositionName, &tmpPosition.InfectionCount)
		for idx, organ := range rd.ResultPerStatistic.PerPosition {
			if organ.PositionName == tmpPosition.PositionName {
				rd.ResultPerStatistic.PerPosition[idx].InfectionCount =
					tmpPosition.InfectionCount
				rd.ResultPerStatistic.PerPosition[idx].Rate = float32(
					float64(rd.ResultPerStatistic.PerPosition[idx].InfectionCount)/
						float64(rd.ResultPerStatistic.PerPosition[idx].TotalCount)) * 100
			}
		}
	}

	perTagTotalQuery := `
select count(target_name)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where ti.user_no = $1
  and (ti.tag1 = $2
    or ti.tag2 = $2
    or ti.tag3 = $2)
`
	perTagQuery := `
select count(target_name)
from target_info ti
         left join count_info ci on ti.target_no = ci.target_no
where ti.user_no = $1
  and (link_click_status = true or download_status = true)
  and (ti.tag1 = $2
    or ti.tag2 = $2
    or ti.tag3 = $2)
`
	maxTag := ClassificationTag{}
	for idx, v := range rd.ProjectInfo.TagArray {
		tmpTag := ClassificationTag{TagName: v.TagName}
		rs1 := conn.QueryRow(perTagTotalQuery, userNo, v.TagNo)
		_ = rs1.Scan(&tmpTag.TotalCount)
		rs2 := conn.QueryRow(perTagQuery, userNo, v.TagNo)
		_ = rs2.Scan(&tmpTag.InfectionCount)
		tmpTag.Rate = float32(float64(tmpTag.InfectionCount)/
			float64(tmpTag.TotalCount)) * 100
		// rate 가 가장 높은 태그 정보
		if idx == 0 {
			maxTag = tmpTag
		} else if maxTag.Rate < tmpTag.Rate {
			maxTag = tmpTag
		}
		rd.ResultPerStatistic.PerTag = append(rd.ResultPerStatistic.PerTag, tmpTag)
	}
	// 가장 rate 가 높은 태그 정보
	rd.ProjectSummary.ResultPerClassification.CTag = maxTag


	return nil
}
