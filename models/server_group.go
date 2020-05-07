package models

import "github.com/astaxie/beego/orm"

//服务器分组
type ServerGroup struct {
	Id          int
	GroupName   string//组名
	Description string//说明
	Status      int//1-正常   0-删除
	CreateTime  int64//创建时间
	UpdateTime  int64//更新时间
	CreateId    int//创建者id
	UpdateId    int//更新者id
}

func (servergroup *ServerGroup) TableName() string {
	return TableName("task_server_group")
}

func ServerGroupGetList(page, pageSize int, filters ...interface{}) ([]*ServerGroup, int64) {
	//获得任务表的句柄
	query := orm.NewOrm().QueryTable(TableName("task_server_group"))
	//判断是否存在过滤条件
	if len(filters) > 0 {
		//获取过滤条件的长度
		l := len(filters)
		//遍历过滤条件
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()

	list := make([]*ServerGroup, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//根据任务分组id查询任务分组
func ServerGroupGetById(id int) (*ServerGroup, error) {
	obj := &ServerGroup{
		Id:id,
	}
	//查询
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}