package main

import (
	"CronJob/jobs"
	"CronJob/models"
	_ "CronJob/routers"
	"github.com/astaxie/beego"
)

func init() {
	models.Init()
	jobs.InitJob()
}

func main() {
	beego.Run()
}
