package controllers

import (
	"CronJob/models"
	"strconv"
	"strings"
	"time"
)

type RoleController struct {
	BaseController
}

func (this *RoleController)List()  {
	this.Data["pageTitle"] = "角色管理"
	this.display()
}

//加载角色管理中的内容部分
func (this *RoleController)Table()  {
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

	//角色名称
	roleName := strings.TrimSpace(this.GetString("roleName"))


	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if roleName != "" {
		filters = append(filters, "role_name__icontains", roleName)
	}
	//分页查询
	result, count := models.RoleGetList(page, this.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		row["detail"] = v.Detail
		list[k] = row
	}
	this.ajaxList("成功", MSG_OK, count, list)
}

//跳转到角色新增页面
func (this *RoleController) Add() {
	this.Data["zTree"] = true
	this.Data["serverGroup"] = serverLists(this.serverGroups, this.userId)
	this.Data["taskGroup"] = taskGroupLists(this.taskGroups, this.userId)
	this.Data["pageTitle"] = "新增角色"
	this.display()
}

func (this *RoleController) AjaxSave() {
	// var data = {'role_name':role_name,'detail':detail,'nodes_data':nodes_data,'server_group_ids':server_group_ids,'task_group_ids':task_group_ids}
	//var data = {'role_name':role_name,'detail':detail,'nodes_data':nodes_data,'server_group_ids':server_group_ids,'task_group_ids':task_group_ids,'id':id,}
	role := new(models.Role)
	role.RoleName = strings.TrimSpace(this.GetString("role_name"))
	role.Detail = strings.TrimSpace(this.GetString("detail"))

	// 2
	//  2  3  4
	auths := strings.TrimSpace(this.GetString("nodes_data"))

	role.ServerGroupIds = strings.TrimSpace(this.GetString("server_group_ids"))
	role.TaskGroupIds = strings.TrimSpace(this.GetString("task_group_ids"))


	role.UpdateId = this.userId
	role.UpdateTime = time.Now().Unix()
	role.Status = 1


	//接收id
	role_id, _ := this.GetInt("id")

	//保存
	if role_id == 0 {
		role.CreateId = this.userId
		role.CreateTime = time.Now().Unix()
		if id, err := models.RoleAdd(role); err != nil {
			this.ajaxMsg("添加失败!", MSG_ERR)
		}else {
			authsArr := strings.Split(auths, ",")
			ras := make([]models.RoleAuth, 0)
			for _, v := range authsArr {
				ra := models.RoleAuth{}
				aid, _ := strconv.Atoi(v)
				ra.AuthId = aid
				ra.RoleId = id
				ras = append(ras, ra)
			}
			if len(ras) > 0 {
				models.RoleAuthBatchAdd(&ras)
			}
		}
		this.ajaxMsg("", MSG_OK)
	}

	//修改
	role.Id = role_id
	if err := role.Update(); err != nil {
		this.ajaxMsg("更新失败!", MSG_ERR)
	}else {
		//根据角色id删除中间表中的内容
		models.RoleAuthDelete(role_id)

		ras := make([]models.RoleAuth, 0)
		auths = strings.TrimRight(auths, ",")
		authsArr := strings.Split(auths, ",")
		for _, v := range authsArr {
			ra := models.RoleAuth{}
			authId, _ := strconv.Atoi(v)
			ra.AuthId = authId
			ra.RoleId = int64(role_id)
			ras = append(ras, ra)
		}

		if len(ras) > 0 {
			models.RoleAuthBatchAdd(&ras)
		}
	}
	this.ajaxMsg("", MSG_OK)
}

/*
  {{range $k, $v := .serverGroup}}
	<input type="checkbox" name="server_group_id" lay-filter="server_group_id" title="{{$v.GroupName}}" value="{{$v.GroupId}}" {{range $ks,$vs:=$.server_group_ids}} {{if eq $v.GroupId $vs}}checked{{end}}{{end}} lay-skin="primary">
{{end}}
     1          2     3     4
本地服务器组   A组   B组   C组

server_group_ids:  2   4
*/
func (this *RoleController) Edit() {
	this.Data["zTree"] = true//引入zTree的css文件
	this.Data["pageTitle"] = "编辑角色"

	//服务器分组
	this.Data["serverGroup"] = serverLists(this.serverGroups, this.userId)
	//任务分组
	this.Data["taskGroup"] = taskGroupLists(this.taskGroups, this.userId)

	id, err := this.GetInt("id")
	if err != nil || id <= 0 {
		return
	}
	//根据id查询角色
	role, err := models.RoleGetById(id)
	if err != nil {
		return
	}
	row := make(map[string]interface{})
	//server_group_ids  task_group_ids  role_name  id  detail
	row["server_group_ids"] = role.ServerGroupIds
	row["task_group_ids"] = role.TaskGroupIds
	row["role_name"] = role.RoleName
	row["id"] = role.Id
	row["detail"] = role.Detail
	this.Data["role"] = row

	//服务器分组id
	//server_group_ids
	//10,2
	serverGroupIdsArr := strings.Split(role.ServerGroupIds, ",")
	serverGroupIds := make([]int, 0)
	for _, v := range serverGroupIdsArr {
		id, _ := strconv.Atoi(v)
		serverGroupIds = append(serverGroupIds, id)
	}
	//10  2
	this.Data["server_group_ids"] = serverGroupIds

	//任务分组id
	//task_group_ids
	taskGroupIdsArr := strings.Split(role.TaskGroupIds, ",")
	taskGroupIds := make([]int, 0)
	for _, v := range taskGroupIdsArr {
		id, _ := strconv.Atoi(v)
		taskGroupIds = append(taskGroupIds, id)
	}
	//10  2
	this.Data["task_group_ids"] = taskGroupIds

	//auth
	//被编辑人所对应的权限id
	roleAuth, _ := models.RoleAuthGetById(id)
	authId := make([]int, 0)
	for _, v := range roleAuth {
		authId = append(authId, v.AuthId)
	}
	this.Data["auth"] = authId
	this.display()
}

//删除角色
func (this *RoleController) AjaxDel() {
	//接收角色id
	role_id, _ := this.GetInt("id")
	//查询角色
	role, _ := models.RoleGetById(role_id)
	role.Status = 0//0是被删除状态
	role.UpdateId = this.userId
	role.UpdateTime = time.Now().Unix()
	//更新
	if err := role.Update(); err != nil {
		this.ajaxMsg("删除失败!", MSG_ERR)
	}
	//根据角色id删除中间表中的内容
	models.RoleAuthDelete(role_id)
	this.ajaxMsg("", MSG_OK)
}