package account

import (
	"bitsync/logic"
	"regexp"
	"bitsync/services"
	"bitsync/models/member"
	"github.com/astaxie/beego"
)

type AccountLogic struct {
	logic.BaseLogic
}

// 注册
func (l *AccountLogic) Register(handset, password, pin string) (UID int) {
	// 校验密码强度
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", password)
	if !match {
		l.BadRequest("密码长度必须是8-15个字符")
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个字母")
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个数字")
	}

	// 校验验证码
	sms := services.PinService{}
	_, err := sms.Validate(services.PIN_REGISTER, handset, pin)
	if err != nil {
		l.Warn(err.Error())
	}

	account := member.AccountModel{}
	UID, err = account.NewAccount(handset, password, "")
	if err != nil {
		beego.Error(err)
		l.ServerError()
	}

	return
}
