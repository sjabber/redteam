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

// 템플릿 api (Restful)
// - 악성메일 템플릿을 만들 수 있는 화면
// - 만들어진 악성메일 템플릿들을 조회, 수정, 삭제 할 수 있어야한다.
// - 먼저 리스트형태로 조회한다음
// - todo 클릭하면 편집할 수 있는 팝업이 뜨도록 만든다.

// 템플릿 등록하기
//func PostTemplateList(c *gin.Context) {
//	userID := c.GetString("email") // httpheader.go 의 AuthMiddleware 에 셋팅되어있음
//	db, _ := c.Get("db")           // httpheader.go 의 DBMiddleware에 셋팅되어있음.
//	conn := db.(sql.DB)
//
//	tmp := model.Template{}
//	c.ShouldBindJSON(&tmp)
//	err := tmp.Create(&conn, userID)
//	if err != nil {
//		log.Print(err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, tmp)
//}

// 템플릿 목록 가져오기
func GetTemplateList(c *gin.Context) {
	//num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tmp, err := model.ReadAll(&conn, num)
	if err != nil {
		// 템플릿을 읽어오는데 오류가 발생한 경우 500에러를 반환한다.
		log.Print("GetTemplate error occurred, account : ", c.Keys["email"])
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tmpls": tmp, "email": c.Keys["email"]})
}

func TemplateDetail(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tn := c.Query("template_no")
	tmpNo, _ := strconv.Atoi(tn)
	fmt.Print(tmpNo)

	//tmp_no := 1

	tmp, err := model.Detail(&conn, num, tmpNo)
	if err != nil {
		log.Println("Error fetching template details :", c.Keys["email"])
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"tmpls": tmp})
	}
}

// 템플릿 수정하기
func EditTemplate(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)
	err := tmp.Update(&conn, num)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// 템플릿 삭제하기
//func DeleteTemplateList(c *gin.Context) {
//	db, _ := c.Get("db")
//	num := c.Keys["number"].(int)
//	conn := db.(sql.DB)
//
//	tmp := model.Template{}
//	c.ShouldBindJSON(&tmp)
//	err := tmp.Delete(&conn, num)
//	if err != nil {
//		log.Print(err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//}

func DelTml(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)

	err := tmp.Delete(&conn, num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"isOK": true})
	}
}
