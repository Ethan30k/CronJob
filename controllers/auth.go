package controllers

import (
	"CronJob/models"
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
