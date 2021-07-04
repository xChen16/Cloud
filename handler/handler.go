package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/16Cloud/meta"
	"github.com/16Cloud/util"
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

			fileMeta := meta.FileMeta{
				FileName: head.Filename,
				Location: "/tmp/" + head.Filename,
				UploadAt: time.Now().Format("2021-0701 08:01:01"),
			}
			newFile, err := os.Create(fileMeta.Location)
			if err != nil {
				fmt.Printf("failed to create file err:%s\n", err.Error())

			}
			defer newFile.Close()
			fileMeta.FileSize, err = io.Copy(newFile, file)
			if err != nil {
				fmt.Printf("failed to save data into file err:%s", err.Error())
				return
			}
			newFile.Seek(0, 0)
			fileMeta.FileSha1 = util.FileSha1(newFile)
			meta.UpdateFileMeta(fileMeta)
			http.Redirect(w, r, "file/upload/suc", http.StatusFound)
		}

	}
}

func UploaSucdHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 获取元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	//fm, _ := meta.GetFileMetaDB(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}
