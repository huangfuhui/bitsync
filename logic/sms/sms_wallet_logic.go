package sms

import (
	"bitsync/logic"
	"bitsync/models/sms"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/jinzhu/now"
)

type SmsWalletLogic struct {
	logic.BaseLogic
}

// 查询钱包余额
func (l *SmsWalletLogic) Get() orm.Params {
	UID := l.GetUID()

	walletModel := sms.SmsWalletModel{}
	wallet, err := walletModel.Wallet(UID)
	if err != nil {
		beego.Error(err)
		l.Warn("查询失败")
		return orm.Params{}
	}

	recordModel := sms.SmsConsumeRecordModel{}
	today, err := recordModel.ConsumeQuantity(UID, now.BeginningOfDay(), now.EndOfDay())
	if err != nil {
		beego.Error(err)
		l.Warn("查询失败")
		return orm.Params{}
	}
	currentMonth, err := recordModel.ConsumeQuantity(UID, now.BeginningOfMonth(), now.EndOfMonth())
	if err != nil {
		beego.Error(err)
		l.Warn("查询失败")
		return orm.Params{}
	}

	return orm.Params{
		"balance":              wallet.Balance,
		"prepare_consume":      wallet.PrepareConsume,
		"today_consume":        today,
		"current_mont_consume": currentMonth,
	}
}
