package account

import (
	"bitsync/logic"
	"regexp"
	"bitsync/services"
	"bitsync/models/member"
	"github.com/astaxie/beego"
	"crypto/md5"
	"io"
	"fmt"
)

type AccountLogic struct {
	logic.BaseLogic
}

// 注册
func (l *AccountLogic) Register(handset, password, pin string) (UID int) {
	account := member.AccountModel{}
	exists := account.Exists(handset)
	if exists > 0 {
		l.Warn("账号已经存在")
		return
	}

	// 校验密码强度
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", password)
	if !match {
		l.BadRequest("密码长度必须是8-15个字符")
		return
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个字母")
		return
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", password)
	if !match {
		l.BadRequest("密码至少包含一个数字")
		return
	}

	// 校验验证码
	sms := services.PinService{}
	_, err := sms.Validate(services.PIN_REGISTER, handset, pin)
	if err != nil {
		l.Warn(err.Error())
		return
	}

	// 密码加密
	salt := beego.AppConfig.String("secretkey")
	w := md5.New()
	io.WriteString(w, salt+password)
	password = fmt.Sprintf("%x", w.Sum(nil))

	UID, err = account.NewAccount(handset, password, "")
	if err != nil {
		beego.Error(err)
		l.ServerError()
		return
	}

	return
}
