package main

import (
	"CronJob/jobs"
	"CronJob/models"
	_ "CronJob/routers"
	"github.com/astaxie/beego"
	"time"
)

func init() {
	models.Init(time.Now().Unix())
	jobs.InitJob()
}

func main() {
	beego.Run()
}
