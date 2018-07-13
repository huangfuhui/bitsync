package member

import (
	"bitsync/controllers"
	"bitsync/validator"
	"bitsync/validator/account"
	"bitsync/logic"
	"bitsync/logic/Member"
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

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
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

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
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

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
	res := l.Login(handset, password)

	c.Output(res)
}

// 修改密码
func (c *AccountController) ModifyPassword() {
	oldPwd := c.GetString("oldPwd")
	newPwd := c.GetString("newPwd")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.ModifyPassword{
		oldPwd,
		newPwd,
	})
	if !ok {
		return
	}

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
	l.ModifyPassword(oldPwd, newPwd)

	c.Output("")
}

// 发送重置密码操作验证码
func (c *AccountController) PasswordPin() {
	handset := c.GetString("handset")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.RegisterPIN{
		handset,
	})
	if !ok {
		return
	}

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
	l.PasswordPin(handset)

	c.Output("")
}

// 重置密码
func (c *AccountController) ResetPassword() {
	handset := c.GetString("handset")
	pin := c.GetString("pin")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.ResetPassword{
		handset,
		pin,
	})
	if !ok {
		return
	}

	l := Member.AccountLogic{logic.BaseLogic{c.BaseController}}
	newPwd, _ := l.ResetPassword(handset, pin)

	c.Output(map[string]string{"new_password": newPwd})
}
