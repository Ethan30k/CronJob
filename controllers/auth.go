package controllers

import (
	"CronJob/models"
	"strings"
	"time"
)

type AuthController struct {
	BaseController
}

//获取树形菜单的所有节点
func (this *AuthController) GetNodes() {
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	//分页查询
	result, count := models.AuthGetList(1, 1000, filters...)
	/*
	name, //节点显示的文本
	open, //节点是否展开
	id,  //节点的标识属性，对应的是启用简单数据格式时idKey对应的属性名，并不一定是id,如果setting中定义的idKey:"zId",那么此处就是zId
	pId, //节点parentId属性，命名规则同id
  }
	*/
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["name"] = v.AuthName
		row["open"] = true
		row["id"] = v.Id
		row["pId"] = v.Pid
		list[k] = row
	}
	this.ajaxList("成功", MSG_OK, count, list)
}

func (this *AuthController) List() {
	this.Data["zTree"] = true
	this.Data["pageTitle"] = "权限因子"
	this.display()
}

/*
auth_name: 任务列表
pid: 31
auth_url: /task/aaa
icon:
sort: 13
is_show: 1
id: 0

auth_name: 编辑11111111
pid: 31
auth_url: /task/edit
icon:
sort: 100
is_show: 0
id: 37
*/
func (this *AuthController) AjaxSave() {
	auth := new(models.Auth)
	auth.AuthName = strings.TrimSpace(this.GetString("auth_name"))
	auth.Pid, _ = this.GetInt("pid")
	auth.AuthUrl = strings.TrimSpace(this.GetString("auth_url"))
	auth.Icon = strings.TrimSpace(this.GetString("icon"))
	auth.Sort, _ = this.GetInt("sort")
	auth.IsShow, _ = this.GetInt("is_show")
	id, _ := this.GetInt("id")

	auth.UpdateId = this.userId
	auth.UpdateTime =  time.Now().Unix()
	auth.Status = 1
	//添加
	if id == 0 {
		auth.CreateId = this.userId
		auth.CreateTime =  time.Now().Unix()
		if _, err := models.AuthAdd(auth); err != nil {
			this.ajaxMsg("添加失败!", MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}
	auth.Id = id
	if err := auth.Update(); err != nil {
		this.ajaxMsg("修改失败!", MSG_ERR)
	}
	this.ajaxMsg("", MSG_OK)
}
/*
auth_url
icon
sort
is_show
*/
//获取某一个节点
func (this *AuthController) GetNode() {
	id, _ := this.GetInt("id")
	result, _ := models.AuthGetById(id)
	row := make(map[string]interface{})
	row["auth_url"] = result.AuthUrl
	row["icon"] = result.Icon
	row["sort"] = result.Sort
	row["is_show"] = result.IsShow
	this.ajaxList("成功", MSG_OK, 0, row)
}

//删除权限
func (this *AuthController) AjaxDel() {
	id, _ := this.GetInt("id")
	auth, _ := models.AuthGetById(id)
	auth.Status = 0//0是被删除状态
	auth.UpdateId = this.userId
	auth.UpdateTime = time.Now().Unix()
	if err := auth.Update(); err != nil {
		this.ajaxMsg("删除失败!", MSG_ERR)
	}
	this.ajaxMsg("", MSG_OK)
}