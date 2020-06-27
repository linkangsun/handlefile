package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	fileStorePath  = "/Users/rick/UploadRepo/"
	defaultBufSize = 1024
)

func main() {

	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/getlist", getFileList)
	logrus.Printf("Listening: 8080\n")
	logrus.Fatal(http.ListenAndServe(":8080", nil))

}

var uploadFile = func(w http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(32 << 20)
	fhs, ok := req.MultipartForm.File["uploadfile"]
	if !ok {
		fmt.Fprintf(w, "key not ok")
		return
	}

	if len(fhs) == 0 {
		fmt.Fprintf(w, "len(fhs) == 0")
		return
	}

	for _, fh := range fhs {
		src, err := fh.Open()
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		dst, _ := os.Create(fileStorePath + fh.Filename)

		buf := make([]byte, defaultBufSize)
		io.CopyBuffer(dst, src, buf)
		logrus.Printf("loop crement")
	}
	fmt.Fprintf(w, "done")
}

var downloadFile = func(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	for _, filename := range q["file"] {
		file, _ := os.Open(fileStorePath + filename)
		defer file.Close()
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		io.Copy(w, file)
	}
}

var getFileList = func(w http.ResponseWriter, req *http.Request) {
	filesInfo, _ := ioutil.ReadDir(fileStorePath)
	fileList := []string{}
	for _, fileInfo := range filesInfo {
		fileList = append(fileList, fileInfo.Name())
	}
	retJSON, _ := json.Marshal(fileList)
	io.Copy(w, bytes.NewReader(retJSON))
}
