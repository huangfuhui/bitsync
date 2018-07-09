package sms

import "bitsync/object"

type SmsWallet struct {
	object.Base
	UID            int `orm:"column(UID)"`
	Balance        int
	PrepareConsume int
}
