package member

import (
	"bitsync/controllers"
	"bitsync/logic"
	"bitsync/logic/Member"
	"bitsync/validator"
	"bitsync/validator/account"
)

type MemberController struct {
	controllers.BaseController
}

// 获取会员信息
func (c *MemberController) Get() {
	l := Member.MemberLogic{logic.BaseLogic{c.BaseController}}
	res := l.Get()

	c.Output(res)
}

// 更新会员信息
func (c *MemberController) Update() {
	name := c.GetString("name")
	email := c.GetString("email")
	sex, _ := c.GetInt("sex", 1)
	birthday := c.GetString("birthday")

	v := validator.BaseValidator{}
	ok := v.Validate(&c.BaseController, account.Update{
		name,
		email,
		sex,
		birthday,
	})
	if !ok {
		return
	}

	l := Member.MemberLogic{logic.BaseLogic{c.BaseController}}
	l.Update(name, email, sex, birthday)

	c.Output("")
}
