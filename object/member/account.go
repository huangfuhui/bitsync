package member

import (
	"bitsync/object"
	"time"
)

type Account struct {
	object.Base
	UID           int
	Account       string
	Password      string
	WechatOpeonid string
	Status        int
	RegisterTime  time.Time
	LastLoginTime time.Time
	LastLoginIp   time.Time
}
