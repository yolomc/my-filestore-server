package handler

import (
	"fmt"
	"io/ioutil"
	"math"
	"my-filestore-server/db"
	"my-filestore-server/redis"
	"my-filestore-server/util"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//MultipartUploadInfo 分块上传结构
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

//InitialMultipartUploadHandler 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.Data(http.StatusInternalServerError, "Application/json",
			util.NewRespMsg(util.StatusServerError, "参数不正确", nil).JSONBytes())
		return
	}

	//生成初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	//将初始化信息写入redis
	redis.HSet("MP_"+upInfo.UploadID,
		"filehash", upInfo.FileHash,
		"filesize", upInfo.FileSize,
		"chunkcount", upInfo.ChunkCount,
	)

	//将初始化数据返回到客户端
	c.Data(http.StatusOK, "application/json",
		util.NewRespMsg(util.StatusOK, "分块信息初始化成功", nil).JSONBytes())
}

// UploadPartHandler 上传文件分块
func UploadPartHandler(c *gin.Context) {

	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	//保存分块
	fpath := "/home/yolo/upload/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744) //创建目录
	file, err := os.Create(fpath)      //创建文件
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "创建文件失败",
			"code": util.StatusCreateFileError,
		})
		return
	}
	defer file.Close()

	//将Form的data 写入到file
	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		file.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 当前块信息写入redis
	redis.HSet("MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果到客户端
	c.Data(http.StatusOK, "application/json",
		util.NewRespMsg(util.StatusOK, "分块信息"+chunkIndex+"上传成功", nil).JSONBytes())
}

//CompleteUploadHandler 通知上传合并
func CompleteUploadHandler(c *gin.Context) {

	uploadid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	filehash := c.Request.FormValue("filehash")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	//通过uplaodid查询redis，并判断分块上传是否完成
	data, err := redis.HGetAll("MP_" + uploadid)
	if err != nil {
		c.Data(http.StatusInternalServerError, "application/json",
			util.NewRespMsg(util.StatusRedisGetError, "redis获取数据失败", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		field := string(data[i].([]byte))
		value := string(data[i+1].([]byte))
		switch {
		case field == "chunkcount":
			totalCount, _ = strconv.Atoi(value)
		case strings.HasPrefix(field, "chkidx_") && value == "1":
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.Data(http.StatusBadRequest, "application/json",
			util.NewRespMsg(util.StatusPartNotFull, "分块未全部上传", nil).JSONBytes())
		return
	}

	//合并分块文件
	partFileStorePath := "/home/yolo/upload/" + uploadid + "/" // 分块所在的目录
	fileStorePath := "/home/yolo/upload/" + filename           // 最后文件保存的路径
	if _, err := mergeAllPartFile(chunkCount, partFileStorePath, fileStorePath); err != nil {
		c.Data(http.StatusInternalServerError, "application/json",
			util.NewRespMsg(util.StatusMergeError, "分块归并失败", nil).JSONBytes())
		return
	}
	//删除redis缓存数据
	redis.Del("MP_" + uploadid)
	//删除分块文件
	os.RemoveAll(partFileStorePath)

	//更新唯一文件表和用户文件表
	if ok := db.OnFileUploadFinished(filehash, filename, int64(filesize), ""); !ok {
		c.Data(http.StatusInternalServerError, "application/json",
			util.NewRespMsg(util.StatusServerError, "数据写入失败", nil).JSONBytes())
		return
	}
	if ok := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize)); !ok {
		c.Data(http.StatusInternalServerError, "application/json",
			util.NewRespMsg(util.StatusServerError, "数据写入失败", nil).JSONBytes())
		return
	}

	c.Data(http.StatusOK, "application/json",
		util.NewRespMsg(util.StatusOK, "分块处理完毕", nil).JSONBytes())
}

//mergeAllPartFile 将分块文件合并成原文件，成功后删除分块文件
//chunkCount:分块数量 partFileStorePath 分块存储的路径 fileStorePath 文件最终地址
//参考：【golang 大文件分割 https://studygolang.com/articles/2687】
func mergeAllPartFile(chunkCount int, partFileStorePath, fileStorePath string) (bool, error) {
	fii, err := os.OpenFile(fileStorePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	for i := 1; i <= chunkCount; i++ {
		f, err := os.OpenFile(partFileStorePath+strconv.Itoa(int(i)), os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		fii.Write(b)
		f.Close()
	}
	fmt.Println(fileStorePath, " has been merge complete")
	return true, nil
}
