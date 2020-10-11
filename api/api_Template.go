package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
)

// 템플릿 api (Restful)
// - 악성메일 템플릿을 만들 수 있는 화면
// - 만들어진 악성메일 템플릿들을 조회, 수정, 삭제 할 수 있어야한다.
// - 먼저 리스트형태로 조회한다음
// - 클릭하면 편집할 수 있는 팝업이 뜨도록 만든다.
// 저번에 만드 jwt 토큰을 활용하도록 한다.

func PostTemplateList(c *gin.Context) {
	userID := c.GetString("email") // httpheader.go 의 AuthMiddleware 에 셋팅되어있음
	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware에 셋팅되어있음.
	conn := db.(sql.DB)

	tmp := model.Template{}
	c.ShouldBindJSON(&tmp)
	err := tmp.Create(&conn, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"템플릿 생성 오류": err.Error()})
		return
	}

	// 제대로 생성되면 statusOk 신호 (200번) 을 보낸다.
	c.JSON(http.StatusOK, tmp)
}

func GetTemplateList(c *gin.Context) {
	tmps, err := model.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"템플릿 목록" : tmps})

}