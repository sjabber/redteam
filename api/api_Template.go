package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

// 템플릿 api (Restful)
// - 악성메일 템플릿을 만들 수 있는 화면
// - 만들어진 악성메일 템플릿들을 조회, 수정, 삭제 할 수 있어야한다.
// - 먼저 리스트형태로 조회한다음
// - todo 클릭하면 편집할 수 있는 팝업이 뜨도록 만든다.


// 템플릿 등록하기
func PostTemplateList(c *gin.Context) {
	userID := c.GetString("email") // httpheader.go 의 AuthMiddleware 에 셋팅되어있음
	db, _ := c.Get("db")           // httpheader.go 의 DBMiddleware에 셋팅되어있음.
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)
	err := tmp.Create(&conn, userID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Template creation error": err.Error()})
		return
	}

	// 제대로 생성되면 statusOk 신호 (200번) 을 보낸다.
	c.JSON(http.StatusOK, tmp)
}

// 템플릿 목록 가져오기
func GetTemplateList(c *gin.Context) {
	tmp, err, num := model.ReadAll()
	if err != nil {
		log.Print(err.Error())
		//case 400 : 템플릿을 DB 로부터 읽어오는데 오류가 발생.
		//case 401 : 읽어온 정보를 바인딩하는데 오류가 발생.
	}
	c.JSON(num, gin.H{"tmpls": tmp, "email": c.Keys["email"]})
}

// 템플릿 수정하기
func PutTemplateList(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)
	err := tmp.Update(&conn)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Template modification error": err.Error()})
		return
	}
}

// 템플릿 삭제하기
func DeleteTemplateList(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)
	err := tmp.Update(&conn)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Template deletion error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Deleted successfully, account deleted": c.Keys["email"]})
	}
}
