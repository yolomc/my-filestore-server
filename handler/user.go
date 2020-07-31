package handler

import (
	"fmt"
	"my-filestore-server/config"
	"my-filestore-server/db"
	"my-filestore-server/redis"
	"my-filestore-server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//SignupHandler 用户注册页面Get请求
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

//DoSignupHandler 处理注册post请求
func DoSignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	if len(username) < 3 || len(password) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusRegisterFailed,
		})
		return
	}

	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	ok := db.UserSignup(username, encodePwd)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册失败",
			"code": util.StatusRegisterFailed,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"code": util.StatusOK,
	})
}

//SigninHandler 用户登录页面get请求
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

//DoSigninHandler 处理登录post请求
func DoSigninHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	if len(username) < 3 || len(password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "请求参数无效",
			"code": util.StatusLoginFailed,
		})
		return
	}

	//检查是否存在于数据库
	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	userCheck := db.UserSignin(username, encodePwd)
	if !userCheck {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "用户名或密码错误",
			"code": util.StatusLoginFailed,
		})
		return
	}

	//生成访问凭证（token），并保如redis
	token := genToken(username)
	if err := redis.Set("token_"+username, token, config.TOKEN_TTL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "生成Token错误",
			"code": util.StatusLoginFailed,
		})
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
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

//UserInfoHandler 用户信息
func UserInfoHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "查询用户信息失败",
			"code": util.StatusServerError,
		})
		return
	}

	//返回用户数据
	resp := util.NewRespMsg(util.StatusOK, "OK", user)
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

//genToken 生成token (40位字符：md5+时间戳前8位)
func genToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return util.MD5([]byte(username+ts+config.TOKEN_SALT)) + ts[:8]
}
