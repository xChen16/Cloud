package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	//下载
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data err:%s", err)
			defer file.Close()
			newFile, err := os.Create("/tmp/" + head.Filename)
			if err != nil {
				fmt.Printf("failed to create file err:%s\n", err.Error())

			}
			defer newFile.Close()
			_, err = io.Copy(newFile, file)
			if err != nil {
				fmt.Printf("failed to save data into file err:%s", err.Error())
				return
			}
			http.Redirect(w, r, "file/upload/suc", http.StatusFound)
		}

	}
}

func UploaSucdHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
