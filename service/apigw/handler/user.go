package handler

import (
	"context"
	"log"
	"my-filestore-server/service/account/proto"
	"my-filestore-server/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

var (
	userCli proto.UserService
)

func init() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})

	service := micro.NewService(micro.Registry(reg))
	//初始化service，解析命令行参数等
	service.Init()
	//初始化一个rpcClient
	userCli = proto.NewUserService("go.micro.service.user", service.Client())
}

//SignupHandler 用户注册页面Get请求
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

//DoSignupHandler 处理注册post请求
func DoSignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &proto.ReqSignup{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
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

	resp, err := userCli.Signin(context.TODO(), &proto.ReqSignin{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	//登录成功，跳转到主页
	respMsg := util.RespMsg{
		Code: resp.Code,
		Msg:  resp.Message,
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    resp.Token,
		},
	}
	c.Data(http.StatusOK, "application/json", respMsg.JSONBytes())
}

//UserInfoHandler 用户信息
func UserInfoHandler(c *gin.Context) {

	username := c.Request.FormValue("username")

	resp, err := userCli.UserInfo(context.TODO(), &proto.ReqUserInfo{Username: username})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	//返回用户数据
	c.Data(http.StatusOK, "application/json",
		util.NewRespMsg(resp.Code, resp.Message, gin.H{
			"Username":   username,
			"CreateTime": resp.SignupAt,
		}).JSONBytes())
}
