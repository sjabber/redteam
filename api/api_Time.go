package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

/*
	health check
*/

func Time(c *gin.Context) {
	/*
		DB와 커넥션을 맺어 select now()로 디비의 시간을 가져온다.
		그러면 DB와 api 서버 모두의 health check가 가능하다
	*/
	db, err := model.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "db connection error")
		log.Print("Health check error occurred, account : ", Account)
		return
	}
	defer db.Close()

	var time string
	query := "SELECT now()"
	err = db.QueryRow(query).Scan(&time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "health check error")
		log.Print("Health check error occurred, account : ", Account)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"Time": time,
	})
	return
}
