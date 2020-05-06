package controllers

import (
	"CronJob/jobs"
	"CronJob/models"
	"github.com/astaxie/beego"
	"time"
)

type HomeController struct {
	BaseController
}

//加载首页菜单
func (this *HomeController) Index() {
	this.TplName = "public/main.html"
}

//加载首页中间部分
func (this *HomeController) Start() {
	//--------------最近执行的任务，最近执行成功的任务数量，最近执行失败的任务数量--------------
	//查询最近执行的20条任务
	logs, _ := models.TaskLogGetList(1, 20)

	recentLogs := make([]map[string]interface{}, len(logs))

	failJob := 0
	okJob := 0

	//遍历日志
	for k, v := range logs {
		row := make(map[string]interface{})
		//task_name  start_time  id   status
		//根据任务id获取任务
		task, err := models.TaskGetById(v.TaskId)

		if err != nil {
			row["task_name"] = ""
			recentLogs[k] = row
			continue
		}
		row["task_name"] = task.TaskName
		//将时间戳转换为年月日时分秒的格式
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:m:s")
		row["id"] = task.Id
		row["status"] = v.Status
		recentLogs[k] = row

		//任务执行失败
		if v.Status != 0 {
			failJob ++
		} else {
			okJob ++
		}
	}
	this.Data["recentLogs"] = recentLogs
	this.Data["failJob"] = failJob
	this.Data["okJob"] = okJob

	//服务器分组信息
	//[服务器分组id]服务器分组名称
	groups_map := serverGroupLists(this.serverGroups, this.userId)
	//获取到最近即将执行的20条任务
	entries := jobs.GetEntries(20)
	//创建切片，用于存储即将执行的任务
	jobList := make([]map[string]interface{}, len(entries))
	//遍历即将执行的任务的切片
	for k,v := range entries{
		row:= make(map[string]interface{})
		job:= v.Job.(*jobs.Job)
		//根据任务id获取任务
		task, _ := models.TaskGetById(job.GetId())
		row["task_id"] = job.GetId()
		row["task_name"] = task.TaskName
		row["task_group"] = groups_map[task.GroupId]
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
	}
	//即将执行的任务的长度
	startJob:=len(jobList)
	this.Data["startJob"] = startJob
	this.Data["jobs"] = jobList

	this.Layout="public/layout.html"
	this.TplName=this.controllerName +"/"+this.actionName+".html"
}

