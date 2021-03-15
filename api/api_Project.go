package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
	"strconv"
)

func ProjectCreate(c *gin.Context) {
	// 계정번호
	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// 프로젝트 생성
	p := model.Project{}
	err := c.ShouldBindJSON(&p)
	err, code := p.ProjectCreate(&conn, num)
	if err != nil {
		c.JSON(code, gin.H{
			"status": code,
			"isOk":   0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"isOk":   1,
		})
		return
	}
}

func GetProject(c *gin.Context) {
	// 계정번호
	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// 프로젝트 조회
	projects, err := model.ReadProject(&conn, num)
	if err != nil {
		log.Println("GetProject error occurred, account :", c.Keys["email"])
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"isOk":   0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"project_list": projects,
		})
	}
}

func ModifyProject(c *gin.Context) {
	p := model.Project{}
	c.ShouldBindJSON(&p)

	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db")
	conn := db.(sql.DB)
	result, err := p.EndDateModify(&conn, num)

	if err != nil {
		log.Println("GetProject error occurred, account :", c.Keys["email"])
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      http.StatusInternalServerError,
			"isOk":        false,
			"dateCompare": result,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
	}

}

//func EndProjectList(c *gin.Context) {
//	// 계정정보
//	num := c.Keys["number"].(int)
//
//	// DB
//	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
//	conn := db.(sql.DB)
//
//	p := model.Project{}
//	c.ShouldBindJSON(&p)
//
//	err := p.EndProject(&conn, num)
//	if err != nil {
//		log.Println(err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{
//			"status" : http.StatusBadRequest,
//			"isOk": 0,
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"status": http.StatusOK,
//		"isOk": 1,
//	})
//}

func DeleteProject(c *gin.Context) {
	// 계정정보
	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	p := model.ProjectDelete{}
	c.ShouldBindJSON(&p)

	err := p.DeleteProject(&conn, num)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"isOk":   0,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"isOk":   1,
	})
}

func StartProjectList(c *gin.Context) {
	// 사용자 계정정보
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	p := model.ProjectStart2{}
	c.ShouldBindJSON(&p)

	// 프로젝트 상태변경
	err := p.StartProject(&conn, num)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"isOk":   0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"isOk":   1,
		})
	}
}

func GetTag(c *gin.Context) {
	// 계정정보
	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tags, err := model.GetTag(&conn, num)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"isOk":   0,
			"status": http.StatusBadRequest,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk":   1,
			"status": http.StatusOK,
			"tags":   tags, // 태그들
		})
	}

}

func ProjectDetail(c *gin.Context) {
	// 계정정보
	num := c.Keys["number"].(int)

	// DB
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// GET 메서드로 전달받은 템플릿 번호 --> tn에 저장
	tn := c.Query("template_no")
	pn := c.Query("project_no")
	tmpNo, _ := strconv.Atoi(tn)
	pNo, _ := strconv.Atoi(pn)
	fmt.Print(tmpNo)

	tmp, err := model.ProjectDetail(&conn, num, tmpNo, pNo)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"tmpls": tmp})
	}
}
