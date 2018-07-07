package main

import (
	"github.com/astaxie/beego"
	_ "bitsync/models"
	"github.com/astaxie/beego/logs"
	_ "bitsync/util"
	_ "bitsync/routers"
	"bitsync/controllers"
	"bitsync/task"
)

func init() {
	runMode := beego.AppConfig.String("runmode")

	if runMode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}

	beego.SetLogFuncCall(true)
	beego.SetLogger(logs.AdapterFile, `{"filename":"log/bitsync.log","daily":true,"maxdays":7}`)

	// 任务调度
	t := task.BaseTask{}
	t.Execute()
}

func main() {
	beego.ErrorController(&controllers.ErrorController{})

	beego.Info("初始化完成,启动应用...")
	beego.Run()
}
