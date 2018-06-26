package account

import (
	"bitsync/controllers"
	"bitsync/validator"
	"bitsync/validator/account"
	accountLogic "bitsync/logic/account"
	"bitsync/logic"
)

type AccountController struct {
	controllers.BaseController
}

// 用户注册
func (c *AccountController) Register() {
	handset := c.GetString("handset")
	password := c.GetString("password")
	pin := c.GetString("pin")

	v := validator.BaseValidator{}
	v.Validate(&c.BaseController, account.Register{
		handset,
		password,
		pin,
	})

	l := accountLogic.AccountLogic{logic.BaseLogic{c.BaseController}}
	l.Register(handset, password, pin)

	c.Output("")
}

// 发送注册验证码
func (c *AccountController) RegisterPin() {
	handset := c.GetString("handset")

	v := validator.BaseValidator{}
	v.Validate(&c.BaseController, account.RegisterPIN{
		handset,
	})

	l := accountLogic.AccountLogic{logic.BaseLogic{c.BaseController}}
	l.RegisterPin(handset)

	c.Output("")
}
