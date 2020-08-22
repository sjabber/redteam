package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"redteam/backend/model"
)

func UserRegister(c *gin.Context) {
	user := model.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Register()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//db, _ := c.Get("db")
	//conn := db.(sql.DB)
	//err = user.Register(&conn)
	//if err != nil {
	//	fmt.Println("Error in user.Register()")
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	token, err := user.GetAuthToken()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.Email})
}
