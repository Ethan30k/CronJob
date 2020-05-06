package controllers

import (
	"CronJob/libs"
	"CronJob/models"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type LoginController struct {
	BaseController
}

//登录（跳转登录页面）
func (this *LoginController)Login()  {
	this.Data["siteName"]=beego.AppConfig.String("")
	this.TplName = "login/login.html"
}

func (this *LoginController)LoginIn(){
	//接收用户名
	username := strings.TrimSpace(this.GetString("username"))
	//接收密码
	password := strings.TrimSpace(this.GetString("password"))
	//判断用户名和密码是否为空
	if username!="" && password!=""{
		//根据用户名查询管理员
		user, err := models.AdminGetByName(username)
		//密码=MD5(password+salt)
		//查询出错，或者密码不正确
		if err !=nil||user.Password != libs.Md5([]byte(password+user.Salt)){
			this.ajaxMsg("账号或密码错误",MSG_ERR)
		}else if user.Status == 0 {
			this.ajaxMsg("该账户已禁用",MSG_ERR)
		}else {
			//为userid赋值
			this.userId = user.Id
			this.user = user
			//用户id+MD5(登陆ip+密码+密码盐)
			authkey := libs.Md5([]byte(this.getClientIp()+"|"+user.Password+user.Salt))
			this.Ctx.SetCookie("auth", strconv.Itoa(user.Id) +"|"+authkey, 60*60*24*7)
			this.ajaxMsg("登录成功！",MSG_OK)
		}
	}
}