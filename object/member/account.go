package member

import (
	"bitsync/object"
	"time"
)

const STATUS_YES = 1 // 可用
const STATUS_NO = 2  // 停用

type Account struct {
	object.Base
	UID           int `orm:"column(uid)"`
	Account       string
	Password      string
	WechatOpenid string
	Status        int
	RegisterTime  time.Time
	LoginTime     time.Time
	LoginIp       string
	LastLoginTime time.Time
	LastLoginIp   string
}
