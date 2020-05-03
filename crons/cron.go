package crons

import (
	"github.com/gorhill/cronexpr"
	"sort"
	"time"
)

//true:删除    false:不删除
type RemoveCheckFunc func(e *Entry) bool

//用于调度多个Cron任务
type Cron struct {
	//任务列表
	entries []*Entry
	//停止的通道
	stop chan struct{}
	//添加通道
	add chan *Entry
	//删除通道
	remove chan RemoveCheckFunc
	//复制通道
	snapshot chan []*Entry
	//表示Cron是否正在运行
	running bool
}

//每种定时任务的启动方式可能不一样，所以定义一个接口，当我们创建定时任务的时候，可以明确该定时任务该如何运行，可以到那时在实现该方法
type Job interface {
	Run()
}

//Cron定时任务
type Entry struct {
	//时间表达式
	Schedule cronexpr.Expression
	//下一次执行时间
	Next time.Time
	//上一次执行时间
	Prev time.Time

	Job
}

type byTime []*Entry

//求长度
func (s byTime)Len() int {
	return len(s)
}

//交换
func (s byTime) Swap(i,j int) {
	s[i], s[j] = s[j], s[i]
}

//比较大小
func (s byTime)Less(i,j int) bool {
	if s[i].Next.IsZero(){
		return false
	}
	if s[j].Next.IsZero(){
		return true
	}
	return s[i].Next.Before(s[j].Next)
}

func (c *Cron)Start()  {
	c.running = true
	go c.Run()
}

//调度协程
func (c *Cron)Run()  {
	//获取当前时间
	now := time.Now().Local()
	//遍历每一个任务
	for _,entry := range c.entries{
		//根据时间表达式和当前时间计算每一个任务下一次执行的时间
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		//根据下一次的运行时间对任务进行排序
		sort.Sort(byTime(c.entries))

		//获取最近一次要执行的时间
		var effective time.Time
		if len(c.entries) == 0 || c.entries[0].Next.IsZero(){
			effective = now.AddDate(10, 0 ,0)
		}else{
			effective = c.entries[0].Next
		}

		select {
		//最近的一次任务时间到达
		case now = <-time.After(effective.Sub(now)):
			for _,e := range c.entries{
				if e.Next!=effective{
					break
				}
				go e.Job.Run()
				e.Prev=e.Next
				e.Next=e.Schedule.Next(effective)
			}
		//添加
		case newEntry:=<-c.add:
			c.entries = append(c.entries, newEntry)
			newEntry.Next = newEntry.Schedule.Next(now)
		//删除
		case cb:=<-c.remove:
			//创建切片，用于存储未删除的任务
			newEntries := make([]*Entry, 0)
			//遍历原始切片
			for _,e :=range c.entries{
				//删除失败，将该任务追加到新的切片中
				if !cb(e){
					newEntries = append(newEntries, e)
				}
			}
			c.entries = newEntries
		//复制
		case <- c.snapshot:
			c.snapshot<-c.entrySnapshot()
		//停止
		case <-c.stop:
			return
		}
		now = time.Now().Local()
	}
}

//spec:时间表达式
//cmd：任务
func (c *Cron)AddJob(spec string, cmd Job) error {
	//将原生时间表达式转换为Expression
	schedule, err := cronexpr.Parse(spec)
	//处理
	if err != nil{
		return err
	}
	c.Schedule(*schedule, cmd)
	return nil
}

func (c *Cron)Schedule(schedule cronexpr.Expression, cmd Job)  {
	entry:= &Entry{
		Schedule:schedule,
		Job:cmd,
	}
	//判断当前entries是否被别的协程占用
	if !c.running{
		c.entries = append(c.entries, entry)
	}
	c.add <-entry
}

//复制任务列表
func (c *Cron)entrySnapshot() []*Entry {
	entries := []*Entry{}
	//遍历原始切片
	for _,e := range c.entries{
		//错误写法，e的指针在两个切片中都存在
		//entries = append(entries, e)
		entries = append(entries, &Entry{
			Schedule:e.Schedule,
			Next:e.Next,
			Prev:e.Prev,
			Job:e.Job,
		})
	}
	return entries
}



























