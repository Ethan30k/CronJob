package controllers

import (
	"strings"
	"strconv"
	"CronJob/models"
	"github.com/astaxie/beego"
	"time"
)

type GroupController struct {
	BaseController
}

//任务分组列表
func (this *GroupController) List() {
	this.Data["pageTitle"] = "任务分组管理"
	this.display()
}

//任务分组列表
func (this *GroupController) Table() {
	//获取当前页码
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}

	//获取每页显示的数量
	limit, err := this.GetInt("limit")
	if err != nil {
		limit = 30
	}
	this.pageSize = limit

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)

	groupName := strings.TrimSpace(this.GetString("groupName"))
	filters = append(filters, "group_name__icontains", groupName)

	if this.userId != 1 {
		groups := strings.Split(this.taskGroups, ",")
		groupsIdArr := make([]int, 0)
		for _, v := range groups {
			id, _ := strconv.Atoi(v)
			groupsIdArr = append(groupsIdArr, id)
		}
		//1  4  6
		filters = append(filters, "id__in", groupsIdArr)
	}
	//分页查询
	result, count := models.GroupGetList(page, this.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	/*
	id
	group_name
	description
	create_time
	update_time
	*/
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] =  v.Id
		row["group_name"] = v.GroupName
		row["description"] = v.Description
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	this.ajaxList("成功", MSG_OK, count, list)

}


//编辑任务分组
func (this *GroupController) Edit() {
	//group_name   description  id
	id, _ := this.GetInt("id", 0)
	group, _ := models.TaskGroupGetById(id)
	row := make(map[string]interface{})
	row["group_name"] = group.GroupName
	row["description"] = group.Description
	row["id"] = group.Id

	this.Data["group"] = row
	this.Data["hideTop"] = true
	this.display()
}

func (this *GroupController) AjaxSave() {
	/*
	group_name: 资源C组
	description: 资源C组
	id: 3
	*/
	id, _ := this.GetInt("id")
	group, err := models.TaskGroupGetById(id)
	if err != nil {
		return
	}

	group.GroupName = strings.TrimSpace(this.GetString("group_name"))
	group.Description = strings.TrimSpace(this.GetString("description"))
	group.Id = id
	group.UpdateId = this.userId
	group.UpdateTime = time.Now().Unix()
	if err := group.Update(); err != nil {
		this.ajaxMsg("更新失败!", MSG_ERR)
	}
	this.ajaxMsg("", MSG_OK)
}