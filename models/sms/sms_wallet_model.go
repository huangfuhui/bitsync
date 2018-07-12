package sms

import (
	"bitsync/models"
	"bitsync/object/sms"
	"github.com/astaxie/beego/orm"
)

type SmsWalletModel struct {
	models.BaseModel
}

// 查询用户短信钱包信息
func (m *SmsWalletModel) Wallet(UID int) (wallet sms.SmsWallet, err error) {
	wallet = sms.SmsWallet{UID: UID}
	_, _, err = orm.NewOrm().ReadOrCreate(&wallet, "UID")
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

// 查询用户短信钱余额
func (m *SmsWalletModel) Balance(UID int) (balance int, err error) {
	wallet, err := m.Wallet(UID)

	return wallet.Balance, err
}
