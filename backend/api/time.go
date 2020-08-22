package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
	health check
*/

func Time(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]string{
		"Time": time.Now().String(),
	})

}
