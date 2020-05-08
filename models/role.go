package models

import "github.com/astaxie/beego/orm"

//角色
type Role struct {
	Id             int
	RoleName       string//角色名称
	Detail         string//备注
	ServerGroupIds string//服务器分组id
	TaskGroupIds   string//任务分组id
	CreateId       int//创建者id
	UpdateId       int//更新者id
	Status         int//状态  1-正常  0-删除
	CreateTime     int64//创建时间
	UpdateTime     int64//更新时间
}

func (role *Role) TableName() string {
	return TableName("uc_role")
}

func RoleGetList(page, pageSize int, filters ...interface{}) ([]*Role, int64) {
	//获得管理员表的句柄
	query := orm.NewOrm().QueryTable(TableName("uc_role"))
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

	list := make([]*Role, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//添加角色
func RoleAdd(role *Role) (int64, error) {
	id, err := orm.NewOrm().Insert(role)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//根据id查询角色
func RoleGetById(id int) (*Role, error) {
	r := new(Role)
	//查询
	err := orm.NewOrm().QueryTable(TableName("uc_role")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//更新
func (t *Role) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}