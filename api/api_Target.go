package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"io"
	"net/http"
	"os"
	"redteam/model"
	"strconv"
	"time"
)

func RegTarget(c *gin.Context) {
	// interface{} -> int 로 형변환하여 num 에 저장한다.
	// num (계정번호) => 등록한 정보를 관리자 번호로 관리하기 위해 사용함.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	target := model.Target{}
	c.ShouldBindJSON(&target)

	errCode, err := target.CreateTarget(&conn, num)
	if err != nil {
		if errCode == 500 {
			model.SugarLogger.Errorf("%v", err.Error())
		}

		c.JSON(errCode, gin.H{
			"isOk": false,
		})
		return
	} else {
		// errCode == status.Ok (200)
		c.Status(errCode)
	}
	return
}

func GetTarget(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	// URL 에 포함된 page 수를 page 변수에 int 로 형변환 후 바인딩.
	pg := c.Query("page")
	page, _ := strconv.Atoi(pg)

	targets, total, pages, err := model.ReadTarget(&conn, num, page) //DB에 저장된 대상들을 읽어오는 메서드
	if err != nil {
		model.SugarLogger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk":    1,
			"status":  http.StatusOK,
			"targets": targets, // 대상 20개
			//"tags": model.GetTag(num), // 태그들
			"total": total, // 대상의 총 갯수
			"pages": pages, // 총 페이지 수
			"page":  page,  // 클릭한 페이지가 몇페이지인지
		})
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
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
		model.SugarLogger.Error(err.Error())
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}

// 업로드한 엑셀파일의 형식에 맞게 작성한 경우 DB에 일괄등록한다.
func ImportTargets(c *gin.Context) {
	// 단일 파일 전송
	f, err := c.FormFile("file")
	if err != nil {
		model.SugarLogger.Error(err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("get form error: %s", err.Error()))
		return
	}

	// 읽을 수 있는 형태로 파일헤더를 가공한다.
	file, err := f.Open()
	if err != nil {
		model.SugarLogger.Error(err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("get form error: %s", err.Error()))
		return
	}
	defer file.Close()

	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	// 아래 코드들부터 전송받은 파일을 읽어 DB에 등록한다.
	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// 바인딩
	target := model.Target{}
	c.ShouldBindJSON(&target)

	// ImportTargets 메세지로 해당 파일을 읽어서 DB에 저장한다.
	errCode, err := target.ImportTargets(&conn, num, file)
	if err != nil {
		model.SugarLogger.Info(err.Error())
		c.JSON(errCode, gin.H{
			"isOk": false,
		})
	} else {
		c.JSON(errCode, gin.H{
			"isOk": true,
		})
	}

}

//사용자가 등록한 대상들을 엑셀파일로 추출한다.
func ExportTarget(c *gin.Context) {
	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db")
	conn := db.(sql.DB)

	// URL 에 포함된 tag 번호를 tagNumber 변수에 int 로 형변환 후 바인딩.
	pg := c.Query("tag_no")
	tagNumber, _ := strconv.Atoi(pg)

	// 해당 계정으로 등록된 훈련대상들의 파일을 생성한다.
	buffer, err := model.ExportTargets(&conn, num, tagNumber) // 클라이언트에게 전달해줄 엑셀파일을 생성하여 아래 코드에서 사용한다.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk ": false,
		})
		return
	}

	// buffer 에 담긴 파일을 전송한다.
	downloadName := time.Now().UTC().Format("Targets-20060102150405.xlsx")
	c.Header("Content-Description", "File Transfer")
	c.Header(cors.ExposeHeadersKey, "Content-Disposition")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
}

func RegTag(c *gin.Context) {
	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	tag := model.Tag{}
	c.ShouldBindJSON(&tag)
	err, errCode := tag.CreateTag(&conn, num)
	if err != nil {
		c.JSON(errCode, gin.H{
			"isOk": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
	}

}

func DeleteTag(c *gin.Context) {
	// num (계정번호) => 해당 계정에 속한 정보들만 추출할 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	tag := model.Tag{}
	c.ShouldBindJSON(&tag)

	err := tag.DeleteTag(&conn, num)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk": true,
		})
	}
}

func Search(c *gin.Context) {
	// num (계정번호) => 해당 계정으로 등록한 정보들만 볼 수 있다.
	num := c.Keys["number"].(int)

	db, _ := c.Get("db") // httpheader.go 의 DBMiddleware 에 셋팅되어있음.
	conn := db.(sql.DB)

	// URL 에 포함된 page 수를 page 변수에 int 로 형변환 후 바인딩.
	pg := c.Query("page")
	searchDivision := c.Query("search_division")
	searchText := c.Query("search_text")
	page, _ := strconv.Atoi(pg)

	targets, total, pages, err := model.SearchTarget(&conn, num, page, searchDivision, searchText) //DB에 저장된 대상들을 읽어오는 메서드
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isOk": false,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"isOk":    true,
			"status":  http.StatusOK,
			"targets": targets,           // 대상 20개
			"total":   total,             // 대상의 총 갯수
			"pages":   pages,             // 총 페이지 수
			"page":    page,              // 클릭한 페이지가 몇페이지인지
		})
	}
}
