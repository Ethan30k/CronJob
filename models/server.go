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

func TaskServerGetList(page, pageSize int, filters ...interface{}) ([]*TaskServer, int64) {
	//获得管理员表的句柄
	query := orm.NewOrm().QueryTable(TableName("task_server"))
	//判断是否存在过滤条件
	if len(filters) > 0 {
		//获取过滤条件的长度
		l := len(filters)
		//遍历过滤条件
		for k := 0; k < l; k += 2 {
			//fmt.Printf("filters[%d] = %v", k, filters[k])
			//fmt.Printf("filters[%d] = %v", k+1, filters[k+1])
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()

	list := make([]*TaskServer, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}
