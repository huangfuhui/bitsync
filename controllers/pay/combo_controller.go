package pay

import (
	"bitsync/controllers"
	"bitsync/logic/pay"
)

type ComboController struct {
	controllers.BaseController
}

// 获取套餐信息
func (c *ComboController) Get() {
	l := pay.ComboLogic{}
	res := l.Get()

	c.Output(res)
}
