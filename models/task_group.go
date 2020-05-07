package models

import "github.com/astaxie/beego/orm"

//任务分组
type Group struct {
	Id          int
	GroupName   string//组名
	Description string//描述
	CreateId    int//创建者id
	CreateTime  int64//创建时间
	UpdateId    int//更新者id
	UpdateTime  int64//更新时间
	Status      int//状态   1-正常  0-删除
}

func (group *Group) TableName() string {
	return TableName("task_group")
}

func GroupGetList(page, pageSize int, filters ...interface{}) ([]*Group, int64) {
	//获得管理员表的句柄
	query := orm.NewOrm().QueryTable(TableName("task_group"))
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

	list := make([]*Group, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//根据任务分组id查询任务分组
func TaskGroupGetById(id int) (*ServerGroup, error) {
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