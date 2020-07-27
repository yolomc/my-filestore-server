package db

import (
	"database/sql"
	"fmt"
	"my-filestore-server/db/mysql"
	"strings"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnFileUploadFinished 文件上传完成后保存文件信息
func OnFileUploadFinished(fileHash string, fileName string, fileSize int64, fileAddr string) bool {

	var sb strings.Builder
	sb.WriteString("insert ignore into tbl_file (")
	sb.WriteString("file_sha1,file_name,file_size,file_addr,status")
	sb.WriteString(") values (?,?,?,?,1)")
	stmt, err := mysql.DBConn().Prepare(sb.String())
	if err != nil {
		fmt.Println("Failed to prepare statement:" + err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Println("Failed to execute statement:" + err.Error())
		return false
	}

	rf, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error of RowsAffected:" + err.Error())
		return false
	}

	if rf <= 0 {
		fmt.Printf("File with hash[%s] hash been uploaded before.", fileHash)
	}
	return true
}

// GetFileMeta 从数据库获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {

	var sb strings.Builder
	sb.WriteString("select file_sha1,file_name,file_size,file_addr from tbl_file ")
	sb.WriteString("where file_sha1=? and status=1 limit 1 ")
	stmt, err := mysql.DBConn().Prepare(sb.String())
	if err != nil {
		fmt.Println("Failed to prepare statement:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	tf := &TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tf.FileHash, &tf.FileName, &tf.FileSize, &tf.FileAddr)
	if err != nil {
		fmt.Println("Error of stmt.QueryRow:" + err.Error())
		return nil, err
	}
	return tf, nil
}

//UpdateFileMeta 更新 file_name
func UpdateFileMeta(fileHash string, fileName string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set file_name=? where file_sha1=?")
	if err != nil {
		fmt.Println("Failed to prepare statement:" + err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(fileName, fileHash)
	if err != nil {
		fmt.Println("Failed to execute statement:" + err.Error())
		return false
	}

	rf, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error of RowsAffected:" + err.Error())
		return false
	}

	if rf <= 0 {
		fmt.Printf("File with hash[%s] hash been uploaded before.", fileHash)
	}
	return true
}

func DeleteFileMeta(fileHash string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_file set status=0 where file_sha1=?")
	if err != nil {
		fmt.Println("Failed to prepare statement:" + err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(fileHash)
	if err != nil {
		fmt.Println("Failed to execute statement:" + err.Error())
		return false
	}

	rf, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error of RowsAffected:" + err.Error())
		return false
	}

	if rf <= 0 {
		fmt.Printf("File with hash[%s] hash been uploaded before.", fileHash)
	}
	return true
}
