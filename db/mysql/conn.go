package mysql

import (
	"database/sql"
	"fmt"
	"os"

	// 导入mysql驱动
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	if err := db.Ping(); err != nil {
		fmt.Println("Failed to connect to mysql:" + err.Error())
		os.Exit(1)
	}
}

//DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
