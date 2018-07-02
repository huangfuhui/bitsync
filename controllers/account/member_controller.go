package account

import (
	"bitsync/controllers"
	"bitsync/logic"
	"bitsync/logic/account"
)

type MemberController struct {
	controllers.BaseController
}

// 获取会员信息
func (c *MemberController) Get() {
	l := account.MemberLogic{logic.BaseLogic{c.BaseController}}
	res := l.Get()

	c.Output(res)
}

// 更新会员信息
func (c *MemberController) Update() {

}
