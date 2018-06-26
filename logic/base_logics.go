package logic

import (
	"bitsync/controllers"
)

type BaseLogic struct {
	controllers.BaseController
}

// 查询请求参数值
func (l *BaseLogic) Param(name string) string {
	return l.GetString(name)
}
