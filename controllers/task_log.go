package controllers

import (
	"CronJob/models"
	"strconv"
	"github.com/astaxie/beego"
	"time"
	"CronJob/libs"
	"CronJob/jobs"
	"strings"
)

//任务日志
type TaskLogController struct {
	BaseController
}

func (this *TaskLogController) List() {
	//task_id=2
	taskId, err := this.GetInt("task_id")
	if err != nil {
		return
	}
	//根据任务id查询任务
	task, err := models.TaskGetById(taskId)
	if err != nil {
		return
	}

	this.Data["task_id"] = taskId
	this.Data["pageTitle"] = "日志管理-" + task.TaskName + "(#" + strconv.Itoa(task.Id) + ")"
	this.display()
}

//task_id=10&page=1&limit=10
func (this *TaskLogController) Table() {
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
	//获取任务id
	taskId, err := this.GetInt("task_id")
	if err != nil {
		return
	}

	status, err := this.GetInt("status")

	filters := make([]interface{}, 0)
	filters = append(filters, "task_id", taskId)

	//接收没有出现错误，并且不是查询所有
	if err == nil && status != 9 {
		filters = append(filters, "status", status)
	}
	//查询
	result, count := models.TaskLogGetList(page, this.pageSize, filters...)
	/*
	id
	task_id
	start_time
	process_time
	output_size
	status
	*/
	//创建切片，用于存储日志，其中每一个map对应的就是一条日志
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] =  v.Id
		row["task_id"] = v.TaskId
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = v.ProcessTime / 1000
		//正常执行
		if v.Status == 0 {
			row["output_size"] = libs.SizeFormat(len([]byte(v.Output)))
		}else {//异常执行
			row["output_size"] = libs.SizeFormat(len([]byte(v.Error)))
		}
		TextStatus := []string{
			"<font color='orange'><i class='fa fa-question-circle'></i> 超时</font>",
			"<font color='red'><i class='fa fa-times-circle'></i> 错误</font>",
			"<font color='green'><i class='fa fa-check-square'></i> 正常</font>",
		}
		//0：正常  -1：错误  -2:超时
		row["status"] = TextStatus[v.Status+2]
		list[k] = row
	}
	this.ajaxList("成功", MSG_OK, count, list)
}


func (this *TaskLogController) Detail() {
	//接收日志id
	id, _ := this.GetInt("id")
	tasklog, err := models.TaskLogGetById(id)
	if err != nil {
		return
	}

	//根据任务id查询任务
	task, err := models.TaskGetById(tasklog.TaskId)
	if err != nil {
		return
	}
	this.Data["task"] = task
	row := make(map[string]interface{})
	//id  start_time process_time  output_size  status
	//output  error
	row["id"] = tasklog.Id
	row["start_time"] = beego.Date(time.Unix(tasklog.CreateTime, 0), "Y-m-d H:i:s")
	row["process_time"] = tasklog.ProcessTime / 1000
	//正常执行
	if tasklog.Status == 0 {
		row["output_size"] = libs.SizeFormat(len([]byte(tasklog.Output)))
	}else {//异常执行
		row["output_size"] = libs.SizeFormat(len([]byte(tasklog.Error)))
	}

	TextStatus := []string{
		"<font color='orange'><i class='fa fa-question-circle'></i> 超时</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 错误</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 正常</font>",
	}
	//0：正常  -1：错误  -2:超时
	row["status"] = TextStatus[tasklog.Status+2]
	row["output"] = tasklog.Output
	row["error"] = tasklog.Error
	this.Data["taskLog"] = row

	//0：正常  -1：错误  -2:超时TextStatus
	this.Data["TextStatus"] = TextStatus[tasklog.Status+2]


	//GroupName
	groupName := "默认分组"
	if task.GroupId > 0 {
		//根据任务分组id查询任务分组
		group, err := models.TaskGroupGetById(task.GroupId)
		if err == nil {
			groupName = group.GroupName
		}
	}
	this.Data["GroupName"] = groupName

	//serverName
	serverName := "本地服务器"
	if task.ServerId > 0 {
		//根据服务器id查询服务器
		server, err := models.TaskSeverGetById(task.ServerId)
		if err == nil {
			serverName = server.ServerName
		}
	}
	this.Data["serverName"] = serverName

	//adminInfo
	if task.NotifyUserIds != "" && task.NotifyUserIds != "0" {
		this.Data["adminInfo"] = jobs.AllAdminInfo(task.NotifyUserIds)
	}else {
		this.Data["adminInfo"] = []*models.Admin{}
	}

	//CreateTime
	this.Data["CreateTime"] = beego.Date(time.Unix(task.CreateTime, 0), "Y-m-d H:i:s")

	//CreateName
	createName := "未知"
	if task.CreateId > 0 {
		//根据id查询管理员
		admin, err := models.AdminGetById(task.CreateId)
		if err == nil {
			createName = admin.RealName
		}
	}
	this.Data["CreateName"] = createName

	//UpdateTime
	this.Data["UpdateTime"] = beego.Date(time.Unix(task.UpdateTime, 0), "Y-m-d H:i:s")

	//UpdateName
	updateName := "未知"
	if task.UpdateId > 0 {
		//根据id查询管理员
		admin, err := models.AdminGetById(task.UpdateId)
		if err == nil {
			updateName = admin.RealName
		}
	}
	this.Data["UpdateName"] = updateName

	this.Data["pageTitle"] = "日志详情" + "(#" + strconv.Itoa(id) + ")"
	this.display()
}

//批量删除日志
func (this *TaskLogController) AjaxDel() {
	ids := this.GetString("ids")
	idArr := strings.Split(ids, ",")

	if len(idArr) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	for _, v := range idArr {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		models.TaskLogDelById(id)
	}
	this.ajaxMsg("", MSG_OK)
}


















