package sms

import "bitsync/object"

type SmsWallet struct {
	object.Base
	UID            int `orm:"column(uid)"`
	Balance        int
	PrepareConsume int
}
