package main

import (
	"fmt"
	"net/http"

	"github.com/16Cloud/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucdHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("failed to start")
	}
}
