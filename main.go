package main

import (
	_ "CronJob/models"
	_ "CronJob/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
