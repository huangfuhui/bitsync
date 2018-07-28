package sms

import (
	"bitsync/models"
	"bitsync/object/sms"
	"github.com/astaxie/beego/orm"
)

type SmsFailedTaskModel struct {
	models.BaseModel
}

// 新增一条失败记录
func (m *SmsFailedTaskModel) Add(UID, smsTaskId int, failedReason string) error {
	failed := sms.SmsFailedTask{
		UID:          UID,
		SmsTaskId:    smsTaskId,
		FailedReason: failedReason,
	}

	_, err := orm.NewOrm().Insert(&failed)
	return err
}
