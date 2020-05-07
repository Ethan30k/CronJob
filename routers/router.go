package routers

import (
	"CronJob/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//登录
    beego.Router("/", &controllers.LoginController{},"*:Login")
    beego.Router("/login_in", &controllers.LoginController{},"*:LoginIn")

    //首页
    beego.Router("/home", &controllers.HomeController{},"*:Index")
    beego.Router("/home/start", &controllers.HomeController{},"*:Start")

    //任务
    beego.AutoRouter(&controllers.TaskController{})
    //任务日志
    beego.AutoRouter(&controllers.TaskLogController{})

	//group/list
	beego.AutoRouter(&controllers.GroupController{})
}
