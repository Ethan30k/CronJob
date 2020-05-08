package models

import "github.com/astaxie/beego/orm"

//权限
type Auth struct {
	Id         int//主键
	Pid        int//父级id
	AuthName   string//权限名称
	AuthUrl    string//url地址
	Sort       int//排序
	Icon       string
	IsShow     int//是否隐藏  0：不显示  1：显示
	UserId     int//操作者id
	CreateId   int//创建者id
	UpdateId   int//更新者id
	Status     int//状态 1-正常    0-删除
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
}

func (auth *Auth) TableName() string {
	return TableName("uc_auth")
}

func AuthGetList(page, pageSize int, filters ...interface{}) ([]*Auth, int64) {
	//获得任务表的句柄
	query := orm.NewOrm().QueryTable(TableName("uc_auth"))
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

	list := make([]*Auth, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("pid", "sort").Limit(pageSize, offset).All(&list)
	return list, total
}

//添加权限
func AuthAdd(auth *Auth) (int64, error) {
	return orm.NewOrm().Insert(auth)
}

//根据任务id获取任务
func AuthGetById(id int) (*Auth, error) {
	auth := &Auth{
		Id:id,
	}
	err := orm.NewOrm().Read(auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}


func (auth *Auth) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(auth, fields...); err != nil {
		return err
	}
	return nil
}
