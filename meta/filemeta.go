package meta

import "my-filestore-server/db"

//FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

//var fileMetas map[string]FileMeta

// func init() {
// 	fileMetas = make(map[string]FileMeta, 0)
// }

// UpdateFileMeta 新增/更新 文件元信息
// func UpdateFileMeta(fmeta FileMeta) {
// 	fileMetas[fmeta.FileSha1] = fmeta
// }

// SaveFileMetaToDB 保存文件元信息到数据库
func SaveFileMetaToDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//UpdateFileMetaToDB 更新文件名
func UpdateFileMetaToDB(fmeta *FileMeta) bool {
	return db.UpdateFileMeta(fmeta.FileSha1, fmeta.FileName)
}

// GetFileMeta 通过sha1值获取文件元信息
// func GetFileMeta(fileSha1 string) FileMeta {
// 	return fileMetas[fileSha1]
// }

// GetFileMetaFromDB 从数据库获取文件元信息
func GetFileMetaFromDB(fileSha1 string) *FileMeta {
	tf, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return nil
	}
	return &FileMeta{
		FileSha1: tf.FileHash,
		FileName: tf.FileName.String,
		FileSize: tf.FileSize.Int64,
		Location: tf.FileAddr.String,
	}
}

// RemoveFileMeta 移除文件元信息
// func RemoveFileMeta(fileSha1 string) {
// 	delete(fileMetas, fileSha1)
// }

//RemoveFileMetaFromDB 移除文件元信息
func RemoveFileMetaFromDB(fileSha1 string) bool {
	return db.DeleteFileMeta(fileSha1)
}
