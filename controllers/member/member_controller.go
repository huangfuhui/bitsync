package member

import (
	"bitsync/controllers"
	"bitsync/logic"
	memberLogic "bitsync/logic/member"
	"bitsync/validator"
	"bitsync/validator/member"
)

type MemberController struct {
	controllers.BaseController
}

// 获取会员信息
func (c *MemberController) Get() {
	l := memberLogic.MemberLogic{logic.BaseLogic{c.BaseController}}
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
	ok := v.Validate(&c.BaseController, member.Update{
		name,
		email,
		sex,
		birthday,
	})
	if !ok {
		return
	}

	l := memberLogic.MemberLogic{logic.BaseLogic{c.BaseController}}
	l.Update(name, email, sex, birthday)

	c.Output("")
}
