package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"my-filestore-server/db"
	"my-filestore-server/meta"
	"my-filestore-server/util"
	"net/http"
	"os"
	"strconv"
	"time"
)

//UploadHandler 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		//返回html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	case http.MethodPost:
		//接收文件流并存储到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data: %s", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/home/yolo/upload/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file: %s", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data to file: %s", err.Error())
			return
		}

		//计算hash值
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.SaveFileMetaToDB(fileMeta)

		//更新用户文件表
		r.ParseForm()
		if ok := db.OnUserFileUploadFinished(r.Form.Get("username"), fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize); ok {
			http.Redirect(w, r, "/home", http.StatusFound)
		} else {
			w.Write([]byte("Upload Failed."))
		}

		//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}
}

//TryFastUploadHandler 尝试秒传
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fileMeta == nil {
		w.Write(util.NewRespMsg(-1, "秒传失败，请访问普通上传接口", nil).JSONBytes())
		return
	}

	if ok := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize)); ok {
		w.Write(util.NewRespMsg(0, "秒传成功", nil).JSONBytes())
	} else {
		w.Write(util.NewRespMsg(-2, "秒传失败，请稍后重试", nil).JSONBytes())
	}

}

// UploadSucHandler 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	fMeta := meta.GetFileMetaFromDB(fileHash)

	if fMeta != nil && fMeta.FileSha1 == fileHash {
		data, err := json.Marshal(&fMeta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

//FileQueryHandler 批量获取用户文件信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	limit, _ := strconv.Atoi(r.Form.Get("limit"))
	userFiles, err := db.QueryUserFileMetas(username, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMetaFromDB(fsha1)

	f, err := os.Open(fMeta.Location)
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
	w.Header().Set("content-disposition", "attachment;filename=\""+fMeta.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler 更新文件名称（只更新显示名称，不修改磁盘上的原文件）
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMetaFromDB(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMetaToDB(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//FileDeleteHandler 删除文件元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMetaFromDB(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMetaFromDB(fileSha1)

	w.WriteHeader(http.StatusOK)
}
