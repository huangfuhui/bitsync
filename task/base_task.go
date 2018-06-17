package task

import "github.com/astaxie/beego"

type BaseTask struct {
}

func (task *BaseTask) Execute() {
	beego.Info("【任务调度】开始分发任务.")

	// 价格提醒
	go func(){
		priceWarn := WarnTask{}
		priceWarn.Warn()
	}()


	beego.Info("【任务调度】任务分发完成.")
}
