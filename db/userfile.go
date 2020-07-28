package db

import (
	"fmt"
	"my-filestore-server/db/mysql"
	"strings"
)

//UserFile 用户文件表结构
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

//OnUserFileUploadFinished 生成用户文件表数据
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	var sb strings.Builder
	sb.WriteString("insert ignore into tbl_user_file (")
	sb.WriteString("user_name,file_sha1,file_name,file_size")
	sb.WriteString(") values (?,?,?,?)")
	stmt, err := mysql.DBConn().Prepare(sb.String())
	if err != nil {
		fmt.Println("Error of Prepare:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize)
	if err != nil {
		fmt.Println("Error of Execute:" + err.Error())
		return false
	}

	// rf, err := res.RowsAffected()
	// if err != nil {
	// 	fmt.Println("Error of RowsAffected:" + err.Error())
	// 	return false
	// }
	// if rf == 0 {
	// 	return false
	// }
	return true
}

//QueryUserFileMetas 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	var sb strings.Builder
	sb.WriteString("select file_sha1,file_name,file_size,upload_at,last_update ")
	sb.WriteString("from tbl_user_file where user_name=? limit ? ")
	stmt, err := mysql.DBConn().Prepare(sb.String())
	if err != nil {
		fmt.Println("Error of Prepare:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println("Error of Query:" + err.Error())
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println("Error of Scan:" + err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	// if userCount == 0 {
	// 	fmt.Println("username not found:" + username)
	// 	return false
	// }
	return userFiles, nil
}
