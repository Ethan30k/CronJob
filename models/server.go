package models

import "github.com/astaxie/beego/orm"

//任务服务器
type TaskServer struct {
	Id            int
	GroupId       int//服务器分组id
	ServerName    string//服务器名称
	ServerAccount string//账户名称
	ServerOuterIp string//服务器外网ip
	ServerIp      string//服务器内网ip
	Port          int//端口
	Password      string//服务器密码
	Type          int//登录类型 0：密码登录
	Detail        string//备注
	CreateTime    int64//创建时间
	UpdateTime    int64//更新时间
	Status        int//状态：0-正常  1-删除
}


func (server *TaskServer) TableName() string {
	return TableName("task_server")
}

//根据id查询服务器
func TaskSeverGetById(id int) (*TaskServer, error) {
	obj := &TaskServer{
		Id:id,
	}
	//查询
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}