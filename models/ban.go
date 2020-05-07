package models

import "github.com/astaxie/beego/orm"

//禁止的命令
type Ban struct {
	Id         int
	Code       string//命令
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
	Status     int//状态  0-正常  1-删除
}

func (ban *Ban) TableName() string {
	return TableName("task_ban")
}

func BanGetList(page, pageSize int, filters ...interface{}) ([]*Ban, int64) {
	//获得管理员表的句柄
	query := orm.NewOrm().QueryTable(TableName("task_ban"))
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

	list := make([]*Ban, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}
