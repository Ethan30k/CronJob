package models

import "github.com/astaxie/beego/orm"

//任务日志
type TaskLog struct {
	Id          int
	TaskId      int//任务id创建时间
	Output      string//输出
	Error       string//错误信息
	Status      int//状态     0：正常  -1：错误  -2:超时
	ProcessTime int//任务执行时间
	CreateTime  int64//
}


func (tasklog *TaskLog) TableName() string {
	return TableName("task_log")
}

func TaskLogAdd(t *TaskLog) (int64, error) {
	return orm.NewOrm().Insert(t)
}


func TaskLogGetList(page, pageSize int, filters ...interface{}) ([]*TaskLog, int64) {
	//获得任务表的句柄
	query := orm.NewOrm().QueryTable(TableName("task_log"))
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

	list := make([]*TaskLog, 0)
	//计算偏移量
	offset := (page - 1) * pageSize
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//根据任务日志id查询任务日志
func TaskLogGetById(id int) (*TaskLog, error) {
	obj := &TaskLog{
		Id:id,
	}
	//查询
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}


//根据id删除日志
func TaskLogDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task_log")).Filter("id", id).Delete()
	return err
}