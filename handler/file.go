package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"my-filestore-server/db"
	"my-filestore-server/meta"
	"my-filestore-server/util"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//UploadHandler 上传页面
func UploadHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/index.html")
}

//DoUploadHandler 上传文件处理
func DoUploadHandler(c *gin.Context) {

	//接收文件流并存储到本地
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取文件信息失败",
			"code": util.StatusFormReadError,
		})
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "/home/yolo/upload/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	//创建本地临时文件
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建文件失败",
			"code": util.StatusCreateFileError,
		})
		return
	}
	defer newFile.Close()

	//复制文件，计算大小
	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "复制文件失败",
			"code": util.StatusCopyFileError,
		})
		return
	}

	//计算hash值
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	meta.SaveFileMetaToDB(fileMeta)

	//更新用户文件表
	if ok := db.OnUserFileUploadFinished(c.Request.FormValue("username"), fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize); ok {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "更新tbl_user_file表失败",
			"code": util.StatusStoreToUserFileError,
		})
	}

	//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

}

//TryFastUploadHandler 尝试秒传
func TryFastUploadHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil || fileMeta == nil {
		c.Data(http.StatusBadRequest, "Application/json",
			util.NewRespMsg(util.StatusQueryFileError, "秒传失败，请访问普通上传接口", nil).JSONBytes())
		return
	}

	if ok := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize)); ok {
		c.Data(http.StatusOK, "Application/json",
			util.NewRespMsg(util.StatusOK, "秒传成功", nil).JSONBytes())
	} else {
		c.Data(http.StatusBadRequest, "Application/json",
			util.NewRespMsg(util.StatusFastUploadError, "秒传失败，请稍后重试", nil).JSONBytes())
	}

}

// GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {

	fileHash := c.Request.FormValue("filehash")
	fMeta := meta.GetFileMetaFromDB(fileHash)

	if fMeta != nil && fMeta.FileSha1 == fileHash {
		data, err := json.Marshal(&fMeta)
		if err != nil {
			c.Data(http.StatusInternalServerError, "Application/json",
				util.NewRespMsg(util.StatusServerError, "秒传失败，请稍后重试", nil).JSONBytes())
			return
		}
		c.Data(http.StatusOK, "Application/json", data)
	}
}

//FileQueryHandler 批量获取用户文件信息
func FileQueryHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	userFiles, err := db.QueryUserFileMetas(username, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "查询 user file表失败",
			"code": util.StatusQueryUserFilesError,
		})
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "数据格式有误",
			"code": util.StatusQueryUserFilesError,
		})
		return
	}
	c.Data(http.StatusOK, "Application/json", data)
}

// DownloadHandler 下载文件
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	fMeta := meta.GetFileMetaFromDB(fsha1)

	f, err := os.Open(fMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "文件打开失败",
			"code": util.StatusFileOpenError,
		})
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "文件读取失败",
			"code": util.StatusFileReadError,
		})
		return
	}

	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Header("content-disposition", "attachment; filename=\""+fMeta.FileName+"\"")
	// write data to client
	c.Data(http.StatusOK, "application/octect-stream", data)
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
