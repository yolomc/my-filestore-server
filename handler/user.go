package handler

import (
	"fmt"
	"io/ioutil"
	"my-filestore-server/config"
	"my-filestore-server/db"
	"my-filestore-server/redis"
	"my-filestore-server/util"
	"net/http"
	"time"
)

//SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	ok := db.UserSignup(username, encodePwd)
	if !ok {
		w.Write([]byte("FAILED"))
	}
	w.Write([]byte("SUCCESS"))
}

//SigninHandler 用户登录
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	//验证用户名和密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	userCheck := db.UserSignin(username, encodePwd)
	if !userCheck {
		w.Write([]byte("FAILED"))
		return
	}

	//生成访问凭证（token），并保如redis
	token := GenToken(username)
	if err := redis.Set("token_"+username, token, config.TOKEN_TTL); err != nil {
		w.Write([]byte("FAILED"))
		return
	}

	//登录成功，跳转到主页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/home",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

//UserInfoHandler 用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	//验证token
	//token := r.Form.Get("token")
	// if !IsTokenValid(username, token) {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//返回用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//GenToken 生成token (40位字符：md5+时间戳前8位)
func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return util.MD5([]byte(username+ts+config.TOKEN_SALT)) + ts[:8]
}
