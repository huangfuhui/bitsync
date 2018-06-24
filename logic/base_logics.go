package logic

import (
	"bitsync/controllers"
	"bitsync/error"
)

type BaseLogic struct {
	controllers.BaseController
	error.ApiError
}

// 获取参数
func (logic *BaseLogic) Get(param string) string {
	return logic.GetString(param)
}
