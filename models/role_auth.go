package models

import (
	"bytes"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type RoleAuth struct {
	AuthId int `orm:"pk"`
	RoleId int64
}

func (roleauth *RoleAuth) TableName() string {
	return TableName("uc_role_auth")
}


//根据角色id获取权限id
func RoleAuthGetByIds(RoleIds string) (Authids string,err error) {
	//获得中间表的句柄
	query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
	//通过逗号切割RoleIds
	ids := strings.Split(RoleIds, ",")
	list := make([]*RoleAuth, 0)
	//查询
	_, err = query.Filter("role_id__in", ids).All(&list, "auth_id")
	if err !=nil{
		return "", err
	}
	b := bytes.Buffer{}
	//遍历list，将auth_id拼接成字符串
	for _,v := range list{
		if v.AuthId!=0&& v.AuthId!=1{
			b.WriteString(strconv.Itoa(v.AuthId))
			b.WriteString(",")
		}
	}

	//出去最右侧的逗号
	Authids = strings.TrimRight(b.String(), ",")
	return Authids, nil
}

//批量插入
func RoleAuthBatchAdd(ras *[]RoleAuth) (int64, error) {
	return orm.NewOrm().InsertMulti(len(*ras), ras)
}