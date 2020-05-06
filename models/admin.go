package models

import "github.com/astaxie/beego/orm"

//管理员
type Admin struct {
	Id         int
	LoginName  string//登录名称
	RealName   string//真实姓名
	Password   string//密码
	RoleIds    string//角色id字符串
	Phone      string//联系电话
	Email      string//邮箱
	Salt       string//密码盐
	LastLogin  int64//最后登录时间
	LastIp     string//最后登录的ip
	Status     int//状态 1-正常  0-禁用
	CreateId   int//创建者id
	UpdateId   int//更新者id
	CreateTime int64//创建时间
	UpdateTime int64//更新时间
}


func (admin *Admin) TableName() string {
	return TableName("uc_admin")
}

func AdminGetList(page, pageSize int, filters ...interface{}) ([]*Admin, int64) {
	//获得管理员表的句柄
	query := orm.NewOrm().QueryTable(TableName("uc_admin"))
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

	list := make([]*Admin, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//根据用户名查询管理员
func AdminGetByName(loginName string) (*Admin, error) {
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("login_name", loginName).One(admin)
	if err != nil{
		return nil,err
	}
	return admin, nil
}

//根据id查询管理员
func AdminGetById(id int) (*Admin , error){
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("id", id).One(admin)
	if err !=nil{
		return nil, err
	}
	return admin,nil
}