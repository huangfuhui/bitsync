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
	ok := v.Validate(&c.BaseController, account.Register{
		handset,
		password,
		pin,
	})
	if !ok {
		return
	}

	l := accountLogic.AccountLogic{logic.BaseLogic{c.BaseController}}
	res := l.Register(handset, password, pin)

	c.Output(res)
}

// 发送注册验证码
func (c *AccountController) RegisterPin() {
	handset := c.GetString("handset")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.RegisterPIN{
		handset,
	})
	if !ok {
		return
	}

	l := accountLogic.AccountLogic{logic.BaseLogic{c.BaseController}}
	l.RegisterPin(handset)

	c.Output("")
}

// 登录
func (c *AccountController) Login() {
	handset := c.GetString("handset")
	password := c.GetString("password")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.Login{
		handset,
		password,
	})
	if !ok {
		return
	}

	l := accountLogic.AccountLogic{logic.BaseLogic{c.BaseController}}
	res := l.Login(handset, password)

	c.Output(res)
}
