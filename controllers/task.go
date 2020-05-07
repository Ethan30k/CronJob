package controllers

import (
	"CronJob/jobs"
	"CronJob/models"
	"github.com/astaxie/beego"
	"github.com/gorhill/cronexpr"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	BaseController
}

//任务列表
func (this *TaskController) List() {
	this.Data["taskGroup"] = taskGroupLists(this.taskGroups, this.userId)
	this.Data["groupId"] = 0
	this.Data["pageTitle"] = "任务管理"
	this.display()
}

//显示任务列表具体内容
func (this *TaskController) Table() {
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
	arr := []int{0, 1, 2, 3}
	filters = append(filters, "status__in", arr)

	//接受任务分组id
	groupId, _ := this.GetInt("group_id", 0)

	if groupId == 0 {
		//当前不是超级管理员
		if this.userId != 1 {
			//通过逗号切割任务分组id
			groups := strings.Split(this.taskGroups, ",")
			//将字符串类型的任务分组转换为整型的任务分组
			groupsIds := make([]int, 0)
			for _, v := range groups {
				id, _ := strconv.Atoi(v)
				groupsIds = append(groupsIds, id)

			}
			filters = append(filters, "group_id__in", groupsIds)
		}
	} else {
		filters = append(filters, "group_id", groupId)
	}

	//接收分组名
	taskName := strings.TrimSpace(this.GetString("task_name"))
	if taskName != "" {
		filters = append(filters, "task_name__icontains", taskName)
	}

	//分页查询
	result, count := models.TaskGetList(page, this.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id

		//查询所有分组
		taskGroup := taskGroupLists(this.taskGroups, this.userId)

		groupName := "默认分组"
		if name, ok := taskGroup[v.GroupId]; ok {
			groupName = name
		}

		StatusText := []string{
			"<font color='red'><i class='fa fa-minus-square'></i></font>",
			"<font color='green'><i class='fa fa-check-square'></i></font>",
			"<font color='orange'><i class='fa fa-question-circle'></i></font>",
			"<font color='red'><i class='fa fa-times-circle'></i></font>",
		}

		//任务名称
		row["task_name"] = StatusText[v.Status] + groupName + "-" + v.TaskName
		row["description"] = v.Description
		//根据任务id获取当前任务执行对象
		e := jobs.GetEnteryById(v.Id)
		if e != nil {
			//next_time pre_time execute_times
			row["next_time"] = beego.Date(e.Next, "Y-m-d H:i:s")
			if e.Prev.Unix() > 0 {
				row["pre_time"] = beego.Date(e.Prev, "Y-m-d H:i:s")
			} else if v.PrevTime > 0 {
				row["pre_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			} else {
				row["pre_time"] = "-"
			}
		} else {
			row["next_time"] = "-"
			if v.PrevTime > 0 {
				row["pre_time"] = beego.Date(time.Unix(v.PrevTime, 0), "Y-m-d H:i:s")
			} else {
				row["pre_time"] = "-"
			}
		}

		row["execute_times"] = v.ExecuteTimes
		list[k] = row
	}
	this.ajaxList("成功", MSG_OK, count, list)
}

//任务详情
func (this *TaskController) Detail() {
	//获取任务id
	id, _ := this.GetInt("id")
	//根据任务id查询任务
	task, err := models.TaskGetById(id)
	if err != nil {
		return
	}
	this.Data["task"] = task

	TextStatus := []string{
		"<font color='red'><i class='fa fa-minus-square'></i> 暂停</font>",
		"<font color='green'><i class='fa fa-check-square'></i> 运行中</font>",
		"<font color='orange'><i class='fa fa-question-circle'></i> 待审核</font>",
		"<font color='red'><i class='fa fa-times-circle'></i> 审核失败</font>",
	}
	//-1：删除  0：停用  1：启用  3：不通过 2:待审核
	this.Data["TextStatus"] = TextStatus[task.Status]

	//GroupName
	//任务分组
	groupName := "默认分组"
	if task.GroupId > 0 {
		group, err := models.TaskGroupGetById(task.GroupId)
		if err == nil {
			groupName = group.GroupName
		}
	}
	this.Data["GroupName"] = groupName

	//serverName
	//服务器名称
	serverName := "本地服务器"
	if task.ServerId == 0 {
		serverName = "本地服务器"
	} else if task.ServerId > 0 {
		server, err := models.TaskSeverGetById(task.ServerId)
		if err == nil {
			serverName = server.ServerName
		}
	}
	this.Data["serverName"] = serverName

	//被通知人的id不是默认值并且不是空字符串
	if task.NotifyUserIds != "0" && task.NotifyUserIds != "" {
		this.Data["adminInfo"] = jobs.AllAdminInfo(task.NotifyUserIds)
	} else {
		this.Data["adminInfo"] = []*models.Admin{}
	}

	//CreateTime
	//任务创建时间
	this.Data["CreateTime"] = beego.Date(time.Unix(task.CreateTime, 0), "Y-m-d H:i:s")
	//CreateName
	//任务创建人
	createName := "未知"
	if task.CreateId > 0 {
		admin, err := models.AdminGetById(task.CreateId)
		if err == nil {
			createName = admin.RealName
		}
	}
	this.Data["CreateName"] = createName

	//UpdateTime
	//修改时间
	this.Data["UpdateTime"] = beego.Date(time.Unix(task.UpdateTime, 0), "Y-m-d H:i:s")

	//UpdateName
	//修改人名称
	updateName := "未知"
	if task.UpdateId > 0 {
		admin, err := models.AdminGetById(task.UpdateId)
		if err == nil {
			updateName = admin.RealName
		}
	}
	this.Data["UpdateName"] = updateName
	this.Data["pageTitle"] = "任务详情"

	this.display()
}

//测试
func (this *TaskController) AjaxRun() {
	//获取任务id
	id, _ := this.GetInt("id")
	//根据任务id查询任务
	task, err := models.TaskGetById(id)
	if err != nil {
		this.ajaxMsg("没有该任务，无法执行", MSG_ERR)
	}
	//根据task创建job
	job, err := jobs.NewJobFromTask(task)
	job.Run()
	this.ajaxMsg("", MSG_OK)
}

//编辑
func (this *TaskController) Edit() {
	//获取任务id
	id, _ := this.GetInt("id")
	//根据任务id查询任务
	task, err := models.TaskGetById(id)
	if err != nil {
		return
	}
	if task.Status == 1 {
		this.ajaxMsg("任务正在运行，无法编辑", MSG_ERR)
	}
	this.Data["task"] = task
	//任务分组
	this.Data["taskGroup"] = taskGroupLists(this.taskGroups, this.userId)
	//服务器分组
	this.Data["serverGroup"] = serverLists(this.serverGroups, this.userId)
	//管理员信息
	this.Data["adminInfo"] = jobs.AllAdminInfo("")
	//notify_user_ids
	//5,3,2
	var nodetifyUserIds []int
	if task.NotifyUserIds != "0" && task.NotifyUserIds != "" {
		notifyUserIdStr := strings.Split(task.NotifyUserIds, ",")
		for _, v := range notifyUserIdStr {
			i, _ := strconv.Atoi(v)
			nodetifyUserIds = append(nodetifyUserIds, i)
		}
	}
	this.Data["notify_user_ids"] = nodetifyUserIds
	this.Data["isAdmin"] = this.userId
	this.display()
}

func (this *TaskController) AjaxSave() {
	task_id, _ := this.GetInt("id")
	task, _ := models.TaskGetById(task_id)
	//task_name  group_id  server_id  description  concurrent
	// cron_spec  command  timeout  is_notify  notify_type  notify_user_ids
	task.TaskName = strings.TrimSpace(this.GetString("task_name"))
	task.GroupId, _ = this.GetInt("group_id")
	task.ServerId, _ = this.GetInt("server_id")
	task.Description = strings.TrimSpace(this.GetString("description"))
	task.Concurrent, _ = this.GetInt("concurrent")
	task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
	task.Command = strings.TrimSpace(this.GetString("command"))
	task.Timeout, _ = this.GetInt("timeout")
	task.IsNotify, _ = this.GetInt("is_notify")
	task.NotifyType, _ = this.GetInt("notify_type")
	task.NotifyUserIds = strings.TrimSpace(this.GetString("notify_user_ids"))
	//手动补全
	task.UpdateId = this.userId
	task.UpdateTime = time.Now().Unix()
	//如果是超级管理员，修改的任务无需审核
	if this.userId == 1 {
		task.Status = 0
	} else {
		task.Status = 2
	}

	//检查修改之后的命令中是否包含非法命令
	flag := checkCommand(task.Command)
	if flag {
		this.ajaxMsg("含有禁止命令!", MSG_ERR)
	}

	//判断时间表达式是否合法
	if _, err := cronexpr.Parse(task.CronSpec); err != nil {
		this.ajaxMsg("时间表达式无效!", MSG_ERR)
	}

	if err := task.Update(); err != nil {
		this.ajaxMsg("修改失败!", MSG_ERR)
	}

	this.ajaxMsg("", MSG_OK)
}

//检查是否包含含有禁用命令
//返回值：是否包含
//true：包含  false:不包含
func checkCommand(command string) bool {
	filters := make([]interface{}, 0)
	//查询正常状态的禁用命令
	filters = append(filters, "status", 0)
	//分页查询
	ban, _ := models.BanGetList(1, 10000, filters...)
	for _, v := range ban {
		if strings.Contains(v.Code, command) || strings.Contains(command, v.Code) {
			return true
		}
	}
	return false
}

//复制任务
func (this *TaskController) Copy() {
	id, _ := this.GetInt("id")
	//根据任务id查询任务
	task, err := models.TaskGetById(id)
	if err != nil {
		return
	}
	this.Data["task"] = task
	//任务分组
	this.Data["taskGroup"] = taskGroupLists(this.taskGroups, this.userId)
	//服务器分组
	this.Data["serverGroup"] = serverLists(this.serverGroups, this.userId)
	//管理员信息
	this.Data["adminInfo"] = jobs.AllAdminInfo("")
	this.Data["pageTitle"] = "任务复制"
	this.display()
}

//批量启动任务
func (this *TaskController) AjaxBatchStart() {
	//接收需要启动的任务的id
	ids := this.GetString("ids")
	//通过逗号进行切割
	idArr := strings.Split(ids, ",")
	//判断用户是否选择了任务
	if len(idArr) < 1 {
		this.ajaxMsg("请选择要操作的任务!", MSG_ERR)
	}
	for _, v := range idArr {
		//将任务id转换为整形
		id, _ := strconv.Atoi(v)
		//判断id是否非法
		if id < 1 {
			continue
		}
		//根据任务id查询任务
		if task, err := models.TaskGetById(id); err == nil {
			//根据task创建job
			job, err := jobs.NewJobFromTask(task)
			//判断是否创建成功
			if err == nil {
				//将当前job添加到切片中
				jobs.AddJob(task.CronSpec, job)
				//修改任务的状态为启用状态
				task.Status = 1
				task.Update("status")
			}
		}
	}

	this.ajaxMsg("", MSG_OK)
}

//批量暂停任务
func (this *TaskController) AjaxBatchPause() {
	//接收需要启动的任务的id
	ids := this.GetString("ids")
	//通过逗号进行切割
	idArr := strings.Split(ids, ",")
	//判断用户是否选择了任务
	if len(idArr) < 1 {
		this.ajaxMsg("请选择要操作的任务!", MSG_ERR)
	}
	for _, v := range idArr {
		//将任务id转换为整形
		id, _ := strconv.Atoi(v)
		//判断id是否非法
		if id < 1 {
			continue
		}
		//从切片删除job
		jobs.RemoveJob(id)
		//根据id查询任务
		if task, err := models.TaskGetById(id); err == nil {
			//任务的状态修改停用状态
			task.Status = 0
			//更新
			task.Update("status")
		}
	}
	this.ajaxMsg("", MSG_OK)
}