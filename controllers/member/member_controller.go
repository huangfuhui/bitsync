package member

import (
	"bitsync/controllers"
	"bitsync/logic"
	"bitsync/logic/Member"
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

}
