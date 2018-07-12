package sms

import (
	"bitsync/object"
	"time"
)

const (
	TYPE_THRESOLD_VALUE = 1

	STATUS_WAIT    = 0
	STATUS_SUCCESS = 1
	STATUS_FAIL    = 2
	STATUS_CANCEL  = 3
)

type SmsTask struct {
	object.Base
	UID        int `orm:"column(uid)"`
	TaskId     int
	Type       int
	Status     int
	FinishTime time.Time
}
