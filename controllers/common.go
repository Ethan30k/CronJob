package controllers

import (
	"CronJob/libs"
	"CronJob/models"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

const (
	//成功
	MSG_OK = 0
	//失败
	MSG_ERR = -1
)

type BaseController struct {
	beego.Controller
	userId         int
	user           *models.Admin
	controllerName string
	actionName     string

	serverGroups string //服务器分组id
}

func (this *BaseController) Prepare() {
	//获取控制器名称和方法名称  IndexController_Index.html
	//IndexController  LinkController TagController
	controllerName, actionName := this.GetControllerAndAction()
	//去除控制器名称尾部的Controller并将结果转换为小写
	this.controllerName = strings.ToLower(controllerName[:len(controllerName)-10])
	//将方法名称转换为小写
	this.actionName = strings.ToLower(actionName)

	this.Auth()
}

//获取服务器分组id
func (this *BaseController) dataAuth(user *models.Admin) {
	//不能等于默认值和空字符串
	if user.RoleIds == "0" || user.RoleIds == "" {
		return
	}
	Filters := make([]interface{}, 0)

	Filters = append(Filters, "status", 1)

	//通过逗号进行切割
	RoleIdsArr := strings.Split(user.RoleIds, ",")
	//将字符串类型的RoleIdsArr转换为整型
	RoleIds := make([]int, 0)
	for _, v := range RoleIdsArr {
		id, _ := strconv.Atoi(v)
		RoleIds = append(RoleIds, id)
	}

	Filters = append(Filters, "id__in", RoleIds)
	Result, _ := models.RoleGetList(1, 1000, Filters...)
	serverGroup := ""
	//拼接服务器
	for _, v := range Result {
		serverGroup += v.ServerGroupIds + ","
	}
	this.serverGroups = strings.TrimRight(serverGroup, ",")

}

//验证用户是否登录了
func (this *BaseController) Auth() {
	//获取cookie并通过|切割
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")

	if len(arr) == 2 {
		idstr, authkey := arr[0], arr[1]
		//将字符串id转换为整型
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			//根据id查询管理员
			user, err := models.AdminGetById(userId)
			if err == nil && authkey == libs.Md5([]byte(this.getClientIp()+"|"+user.Password+user.Salt)) {
				this.userId = user.Id
				this.user = user
				//加载菜单
				this.AdminAuth()
				//加载servergroupIds
				this.dataAuth(user)
			}
		}
	} else {
		if this.controllerName != "login" && this.actionName != "login" {
			this.redirect(beego.URLFor("LoginController.Login"))
		}
	}
}

func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
}

func (this *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["message"] = msg
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

//获取用户登录ip
func (this *BaseController) getClientIp() string {
	array := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return array[0]
}

//加载左侧菜单栏
func (this *BaseController) AdminAuth() {
	//创建切片存储过滤条件
	filters := make([]interface{}, 0)
	//查询状态正常的权限
	filters = append(filters, "status", 1)
	//判断是否是超级管理员，如果为1这是超级管理员
	if this.userId != 1 {
		//根据角色id获取权限id
		adminAuthIds, _ := models.RoleAuthGetByIds(this.user.RoleIds)
		//将字符串adminAuthIds转换为切片
		adminAuthIdArr := strings.Split(adminAuthIds, ",")
		filters = append(filters, "id__in", adminAuthIdArr)
	}

	//分页查询
	result, _ := models.AuthGetList(1, 1000, filters...)

	//一级菜单
	list := make([]map[string]interface{}, len(result))

	//二级菜单
	list2 := make([]map[string]interface{}, len(result))
	i, j := 0, 0
	for _, v := range result {
		//创建map,用于存储每一记录
		row := make(map[string]interface{})
		//判断父级id是否为1并且isshow是否为1
		//一级菜单
		if v.Pid == 1 && v.IsShow == 1 {
			row["Icon"] = v.Icon
			row["AuthName"] = v.AuthName
			row["Id"] = v.Id
			list[i] = row
			i++
		}
		//二级菜单
		if v.Pid != 1 && v.IsShow == 1 {
			row["Pid"] = v.Pid
			row["AuthUrl"] = v.AuthUrl
			row["Icon"] = v.Icon
			row["AuthName"] = v.AuthName
			row["Id"] = v.Id
			list2[j] = row
			j++
		}
	}
	//不要忘记切割
	this.Data["SideMenu1"] = list[:i]
	this.Data["SideMenu2"] = list2[:j]
}

//authStr:服务器分组
//adminId:管理员id
func serverGroupLists(authStr string, adminId int) (sgl map[int]string) {
	Filters := make([]interface{}, 0)
	Filters = append(Filters, "status", 1)
	//服务器分组id不为空并且不是超级管理员
	if authStr != "" && adminId != 1 {
		//通过逗号进行切割
		serverGroupIdsArr := strings.Split(authStr, ",")
		serverGroupIds := make([]int, 0)
		for _, v := range serverGroupIdsArr {
			id, _ := strconv.Atoi(v)
			serverGroupIds = append(serverGroupIds, id)
		}
		Filters = append(Filters, "id__in", serverGroupIds)
	}
	//分页查询
	groupResult, n := models.ServerGroupGetList(1, 1000, Filters...)
	sgl = make(map[int]string, n)
	//遍历服务器分组切片，将其中值封装到map中
	for _,gv := range groupResult{
		sgl[gv.Id] = gv.GroupName
	}
	return sgl
}
