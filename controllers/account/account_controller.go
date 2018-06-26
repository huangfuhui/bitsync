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
