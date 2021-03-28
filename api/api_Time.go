package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
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
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	defer conn.Close()

	var time string
	query := "SELECT now()"
	err := conn.QueryRow(query).Scan(&time)
	if err != nil {
		model.SugarLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "health check error")
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"Time": time,
	})
	return
}
