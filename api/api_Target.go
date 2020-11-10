package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"redteam/model"
	"strconv"
)

func RegTarget(c *gin.Context) {
	// interface{} -> int 로 형변환하여 num 에 저장한다.
	// num (계정번호) => 등록한 정보를 관리자 번호로 관리하기 위해 사용함.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)

	err := target.CreateTarget(&conn, num)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"target_registration_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"target_registration_error, register_account": c.Keys["email"]})
	}
}

func GetTarget(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	target, tag, err := model.ReadTarget(num) //DB에 저장된 대상들을 읽어오는 메서드
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"Target read error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"targets": target, "tags": tag, "register_account": c.Keys["email"]})
}

// 헤더를 먼저 정의한다음 파일 다운로드 메서드를 이용하여 파일을 다운로드받도록 한다.
func DeleteTarget(c *gin.Context) {

	// num (계정번호) => 해당 계정에 속한 정보들만 삭제할 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db")
	conn := db.(sql.DB)

	//JSON 이 아닌 배열로 받아온다.

	target := model.TargetNumber{}
	c.ShouldBindJSON(&target)

	err := target.DeleteTarget(&conn, num)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusPaymentRequired, gin.H{"target_deleting_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"delete_success, deleting_account": c.Keys["email"]})
	}
}

// 훈련대상들을 입력할 형식이 담긴 엑셀파일을 클라이언트가 다운받을때 사용하는 api
func DownloadExcel(c *gin.Context) {
	header := c.Writer.Header()
	header["content-type"] = []string{"application/vnd.ms-excel"}
	header["content-disposition"] = []string{"attachment; filename=" + "Target.xlsx"}

	// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (클라이언트에게 줄 엑셀파일을 보관해둘 디렉토리 경로로 수정)
	// 현재는 프로젝트파일의 Spreadsheet 폴더에 보관해둔다.
	destFile := "./Spreadsheet/sample.xlsx"
	file, err := os.Open(destFile)
	if err != nil {
		log.Print(err.Error())
		c.String(http.StatusMethodNotAllowed, "%v", err)
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}

// 업로드한 엑셀파일의 형식에 맞게 작성한 경우 DB에 일괄등록한다.
func ImportTargets(c *gin.Context) {
	// 단일 파일 전송
	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err.Error())
		c.String(http.StatusNotAcceptable, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	// num (int) -> str (string) 변환
	str := strconv.Itoa(num)

	// 파일을 구체적인 장소로 업로드한다. (서버에 파일을 저장할 장소)
	filename := filepath.Base(file.Filename)

	// todo 2 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (클라이언트로부터 다운로드받을 파일을 하나 만든다.)
	// 현재는 컴퓨터의 다운로드파일로 업로드 받는다.
	uploadPath := "./Spreadsheet/" + filename + str
	log.Println(filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		log.Print(err.Error())
		c.String(http.StatusNotAcceptable, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	} else {
		c.String(http.StatusOK, fmt.Sprintf("Status : Posted, File name : %s", filename+str))
	} // 파일 전송이 완료됨.

	/////////////////아래 코드들부터 전송받은 파일을 읽어 DB에 등록한다.////////////////////////////
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)

	// ImportTargets 메세지로 해당 파일을 읽어서 DB에 저장한다.
	err = target.ImportTargets(&conn, uploadPath, num)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"Batch registration error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Batch registration success": c.Keys["email"]})
	}

	// DB에 등록이 완료되어 필요없어진 파일을 삭제하는 코드
	// todo 2 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 2는 전부 같은 경로로 수정)
	err2 := os.Remove("./Spreadsheet/" + filename + str)
	if err2 != nil {
		panic(err2) //현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴
	}
}

//사용자가 등록한 대상들을 엑셀파일로 추출한다.
func ExportTarget(c *gin.Context) {
	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	header := c.Writer.Header()
	header["content-type"] = []string{"application/vnd.ms-excel"}
	header["content-disposition"] = []string{"attachment; filename=" + "Registered_Targets.xlsx"}

	// 해당 계정으로 등록된 훈련대상들의 파일을 생성한다.
	err := model.ExportTargets(num) // 클라이언트에게 전달해줄 엑셀파일을 생성하여 아래 코드에서 사용한다.
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusProxyAuthRequired, gin.H{"Export error ": err.Error()})
	}

	str := strconv.Itoa(num)
	// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (클라이언트에게 줄 엑셀파일을 보관해둘 디렉토리 경로로 수정)
	// 현재는 프로젝트파일의 Spreadsheet 폴더에 보관해둔다.
	destFile := "./Spreadsheet/Registered_Targets" + str + ".xlsx"
	file, err := os.Open(destFile)
	if err != nil {
		c.String(http.StatusOK, "%v", err)
		return
	}
	io.Copy(c.Writer, file)
	file.Close()

	// 사용자가 파일을 다운로드받으면 생성한 파일은 다시 지운다.
	// todo 3 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 3은 전부 같은 경로로 수정)
	err3 := os.Remove("./Spreadsheet/Registered_Targets" + str + ".xlsx")
	if err3 != nil {
		//현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴
		panic(err3)
	}
	os.Clearenv()
}

func RegTag(c *gin.Context) {
	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	tag := model.Tag{}
	c.ShouldBindJSON(&tag)
	err := tag.CreateTag(&conn)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusRequestTimeout, gin.H{"target_registration_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"registering_success, register_account": c.Keys["email"]})
	}

}

func DeleteTag(c *gin.Context) {
	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	tag := model.Tag{}
	c.ShouldBindJSON(&tag)

	err := tag.DeleteTag(&conn)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusConflict, gin.H{"tag_deleting_error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"delete_success, deleting_account": c.Keys["email"]})
	}
}
