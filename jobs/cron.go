package jobs

import (
	"CronJob/crons"
	"sync"
	"github.com/astaxie/beego"
)

var (
	mainCron *crons.Cron
	lock sync.Mutex
)

func init() {
	mainCron = crons.New()
	mainCron.Start()
}

//添加job
func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	//判断是否获取entrie
	if GetEnteryById(job.id) != nil {
		return false
	}

	//添加job
	err := mainCron.AddJob(spec, job)

	if err != nil {
		beego.Error("AddJob: ", err.Error())
		return false
	}
	return true
}
//根据id查找entries
func GetEnteryById(id int) *crons.Entry {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return e
			}
		}
	}
	return nil
}

func GetEntries(size int) []*crons.Entry {
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}