package handler

import (
	"context"
	"fmt"
	"my-filestore-server/config"
	"my-filestore-server/db"
	"my-filestore-server/redis"
	"my-filestore-server/service/account/proto"
	"my-filestore-server/util"
	"time"
)

type User struct{}

func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	username := req.Username
	password := req.Password
	if len(username) < 3 || len(password) < 5 {
		resp.Code = util.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}

	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	ok := db.UserSignup(username, encodePwd)
	if !ok {
		resp.Code = util.StatusRegisterFailed
		resp.Message = "注册失败"
	} else {
		resp.Code = util.StatusOK
		resp.Message = "注册成功"
	}

	return nil
}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {
	username := req.Username
	password := req.Password
	if len(username) < 3 || len(password) < 5 {
		resp.Code = util.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}

	//检查是否存在于数据库
	encodePwd := util.Sha1([]byte(password + config.PWD_SALT))
	userCheck := db.UserSignin(username, encodePwd)
	if !userCheck {
		resp.Code = util.StatusLoginFailed
		resp.Message = "用户名或密码错误"
		return nil
	}

	//生成访问凭证（token），并存入redis
	token := genToken(username)
	if err := redis.Set("token_"+username, token, config.TOKEN_TTL); err != nil {
		resp.Code = util.StatusLoginFailed
		resp.Message = "生成Token错误"
		return nil
	}

	resp.Code = util.StatusOK
	resp.Message = "注册成功"
	resp.Token = token
	return nil
}

//genToken 生成token (40位字符：md5+时间戳前8位)
func genToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return util.MD5([]byte(username+ts+config.TOKEN_SALT)) + ts[:8]
}

func (*User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	username := req.Username
	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		resp.Code = util.StatusServerError
		resp.Message = "查询用户信息失败"
		return nil
	}

	resp.Code = util.StatusOK
	resp.Message = "查询用户信息成功"
	resp.SignupAt = user.SignupAt
	return nil
}
