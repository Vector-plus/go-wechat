package upload

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"wechat/utils"
	"wechat/utils/re"
)

func UPloadfile(writer http.ResponseWriter, request *http.Request) (int, string) {
	srcFile, head, err := request.FormFile("file")
	if err != nil {
		fmt.Println("read file error", err)
		return re.ERROR_FILE_READ, ""
	}
	suffix := ".png"
	name := head.Filename
	tem := strings.Split(name, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create(utils.FILEDIRE + fileName)
	if err != nil {
		fmt.Println("create file error", err)
		return re.ERROR_FILE_READ, ""
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("create file error", err)
		return re.ERROR_FILE_READ, ""
	}
	url := utils.FILEDIRE + fileName
	return re.SUCCSE, url
}

// func DownloadFile (writer http.ResponseWriter, request *http.Request, path string) int{

// 	return re.SUCCSE
// }
