package main

import (
	"github.com/astaxie/beego"
	_ "bitsync/models"
	"github.com/astaxie/beego/logs"
	"bitsync/services"
	_ "bitsync/util"
	_ "bitsync/routers"
)

func init() {
	runMode := beego.AppConfig.String("runmode")

	if runMode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}

	beego.SetLogFuncCall(true)
	beego.SetLogger(logs.AdapterFile, `{"filename":"log/bitsync.log","daily":true,"maxdays":7}`)

	// 启动价格监控
	go services.Watch()
}

func main() {
	beego.Info("初始化完成,启动应用...")
	beego.Run()
}
