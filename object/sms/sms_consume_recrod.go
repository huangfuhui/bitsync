package sms

import (
	"bitsync/object"
	"time"
)

type SmsConsumeRecord struct {
	object.Base
	UID         int `orm:"column(uid)"`
	Handset     string
	Amount      int
	SmsContent  string
	Status      int
	ConsumeTime time.Time
}
