package task

import (
	"bitsync/util"
	"time"
	"github.com/astaxie/beego"
	"strings"
	"strconv"
	"bitsync/services"
)

type WarnTask struct {
}

func (task *WarnTask) Warn() {
	redis := util.Cli{}
	tplId, _ := beego.AppConfig.Int64("sms::tpl_price_warn")
	for {
		time.Sleep(time.Second)

		redis.Select(1)
		priceStr, _ := redis.Get("dragonex:eosusdt")
		if priceStr == "" {
			continue
		}

		dotIndex := strings.Index(priceStr, ".")
		newPrice, _ := strconv.ParseFloat(priceStr, 64)

		redis.Select(0)
		oldPriceStr, _ := redis.Get("13710227972:eosusdt")
		oldPrice, _ := strconv.ParseFloat(oldPriceStr, 64)

		if oldPriceStr == "" {
			redis.Set("13710227972:eosusdt", string([]rune(priceStr)[:dotIndex+1])+"0")
			continue
		}

		phones := []string{"13710227972"}
		if newPrice+0.5 <= oldPrice {
			updatePrice := strconv.FormatFloat(oldPrice-0.5, 'f', 1, 64)
			redis.Set("13710227972:eosusdt", updatePrice)
			beego.Info("【价格提醒】eos/usdt:" + updatePrice)

			for _, v := range phones {
				sms := services.SmsService{}
				err := sms.SendSingle("86", v, []string{"eos/usdt", updatePrice}, tplId)
				if err != nil {
					beego.Error(err)
				}
			}
		} else if newPrice-0.5 >= oldPrice {
			updatePrice := strconv.FormatFloat(oldPrice+0.5, 'f', 1, 64)
			redis.Set("13710227972:eosusdt", updatePrice)
			beego.Info("【价格提醒】eos/usdt:" + updatePrice)

			for _, v := range phones {
				sms := services.SmsService{}
				err := sms.SendSingle("86", v, []string{"eos/usdt", updatePrice}, tplId)
				if err != nil {
					beego.Error(err)
				}
			}
		}
	}
}
