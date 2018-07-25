package market

import (
	"bitsync/controllers"
	"bitsync/logic/market"
)

type ExchangeController struct {
	controllers.BaseController
}

// 获取所有交易所
func (c *ExchangeController) List() {
	l := market.ExchangeLogic{}
	res := l.List()

	c.Output(res)
}
