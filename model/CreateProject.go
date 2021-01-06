package model

type Project struct {
	ProjectName		string	`json:"project_name"` // 프로젝트 이름
	ProjectDesc		string	`json:"project_desc"` // 프로젝트 설명
	ProjectStart	string	`json:"project_start"`// 프로젝트 시작일
	ProjectEnd		string	`json:"project_end"`  // 프로젝트 종료일
	ProjectTemplate	string	`json:"project_template"`

}

func
