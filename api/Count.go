package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"redteam/model"
	"strconv"
)

//	1->true
//	t->true
//	T->true
//	TRUE->true
//	true->true
//	True->true
//	0->false
//	f->false
//	F->false
//	FALSE->false
//	false->false
//	False->false
func CountTarget(c *gin.Context) {
	// 포멧 http://localhost:5000/api/CountTarget?tNo=2&pNo=2&email=true&link=false&download=True
	tNo, err := strconv.Atoi(c.Query("tNo"))
	pNo, err := strconv.Atoi(c.Query("pNo"))
	emailReadStatus, err := strconv.ParseBool(c.Query("email"))
	linkClickStatus, err := strconv.ParseBool(c.Query("link"))
	downloadStatus, err := strconv.ParseBool(c.Query("download"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(sql.DB)
	counter := model.CounterModel{
		TargetNo:        tNo,
		ProjectNo:       pNo,
		EmailReadStatus: emailReadStatus,
		LinkClickStatus: linkClickStatus,
		DownloadStatus:  downloadStatus,
	}

	err = counter.UpdateCount(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
