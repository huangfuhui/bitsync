package account

import (
	"bitsync/logic"
	"regexp"
)

type AccountLogic struct {
	logic.BaseLogic
}

// 注册
func (l *AccountLogic) Register(handset, password, pin string) {
	match := false
	match, _ = regexp.MatchString("^.{8,15}$", password)
	if !match {
		l.Warn("密码长度必须是8-15个字符")
	}
	match, _ = regexp.MatchString("^.*[a-zA-Z].*$", password)
	if !match {
		l.Warn("密码至少包含一个字母")
	}
	match, _ = regexp.MatchString("^.*[0-9].*$", password)
	if !match {
		l.Warn("密码至少包含一个数字")
	}
}
