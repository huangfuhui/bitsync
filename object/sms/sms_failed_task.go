package sms

import (
	"bitsync/object"
)

type SmsFailedTask struct {
	object.Base
	UID          int `orm:"column(uid)"`
	SmsTaskId    int
	FailedReason string
}
