package logic

import "bitsync/controllers"

type BaseLogic struct {
	controllers.BaseController
}

// 获取参数
func (logic *BaseLogic) Get(param string) string {
	return logic.GetString(param)
}
