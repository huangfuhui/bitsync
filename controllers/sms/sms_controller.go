package sms

import (
	"bitsync/controllers"
	"bitsync/logic/sms"
	"bitsync/logic"
)

type IndexController struct {
	controllers.BaseController
}

// 查询钱包余额
func (c *IndexController) Wallet() {
	l := sms.SmsWalletLogic{logic.BaseLogic{c.BaseController}}
	res := l.Get()

	c.Output(res)
}
