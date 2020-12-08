package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/model"
)

func GetTag(c *gin.Context) {
	num := c.Keys["number"].(int)

	err := model.GetTag2(num)
	if err != nil {
		log.Println("GetTag error occurred during CreateProject, account :", c.Keys["email"])
	}

	c.JSON(http.StatusOK, gin.H{"tags" : err})
}