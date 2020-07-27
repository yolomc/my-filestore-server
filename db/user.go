package db

import (
	"fmt"
	"my-filestore-server/db/mysql"
)

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

//UserSignup 用户注册
func UserSignup(username string, pwd string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (user_name,user_pwd) values (?,?)")
	if err != nil {
		fmt.Println("Error of Prepare:" + err.Error())
		return false
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, pwd)
	if err != nil {
		fmt.Println("Error of Execute:" + err.Error())
		return false
	}

	rf, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error of RowsAffected:" + err.Error())
		return false
	}
	if rf == 0 {
		return false
	}
	return true
}

//UserSignin 根据用户名和密码判断用户是否存在
func UserSignin(username string, encodePwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select count(1) from tbl_user where user_name=? and user_pwd=? limit 1")
	if err != nil {
		fmt.Println("Error of Prepare:" + err.Error())
		return false
	}
	defer stmt.Close()

	var userCount int
	err = stmt.QueryRow(username, encodePwd).Scan(&userCount)
	if err != nil {
		fmt.Println("Error of QueryRow/Scan:" + err.Error())
		return false
	}
	if userCount == 0 {
		fmt.Println("username not found:" + username)
		return false
	}
	return true
}

//GetUserInfo 获取UserInfo
func GetUserInfo(username string) (*User, error) {
	stmt, err := mysql.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Error of Prepare:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	// 执行查询的操作
	user := &User{}
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println("Error of QueryRow/Scan:" + err.Error())
		return nil, err
	}
	return user, nil
}
