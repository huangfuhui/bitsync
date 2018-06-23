package member

import (
	"bitsync/object"
	"time"
)

const STATUS_YES = 1 // 可用
const STATUS_NO = 2  // 停用

type Account struct {
	object.Base
	UID           int
	Account       string
	Password      string
	WechatOpeonid string
	Status        int
	RegisterTime  time.Time
	LoginTime     time.Time
	LoginIp       time.Time
	LastLoginTime time.Time
	LastLoginIp   time.Time
}
