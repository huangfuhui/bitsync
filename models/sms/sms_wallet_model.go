package sms

import (
	"bitsync/models"
	"bitsync/object/sms"
	"github.com/astaxie/beego/orm"
)

type SmsWallet struct {
	models.BaseModel
}

// 查询用户短信钱包信息
func (m *SmsWallet) Balance(UID int) (wallet sms.SmsWallet, err error) {
	wallet = sms.SmsWallet{UID: UID}
	_, _, err = orm.NewOrm().ReadOrCreate(&wallet, "UID")
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}
