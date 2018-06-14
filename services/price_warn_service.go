package services

import (
	"bitsync/util"
	"strconv"
	"github.com/astaxie/beego"
	"math"
)

type PriceWarnService struct {
}

func (service *PriceWarnService) Warn() {
	for {
		currentPrice, _ := util.Redis.Get(util.Redis.Con(), "huobi:eosusdt")
		if currentPrice == "" {
			continue
		}
		currentFloat, _ := strconv.ParseFloat(currentPrice, 64)
		current := int(math.Ceil(currentFloat/5) * 5)

		lastPrice, _ := util.Redis.Get(util.Redis.Con(), "price_warn:13710227972")
		newPrice := strconv.FormatInt(int64(current), 10)
		if lastPrice != "" {
			last, _ := strconv.Atoi(lastPrice)

			if last != current {
				// warn := SmsService{}
				// warn.SendSingle("86", "13710227972", []string{"EOS", currentPrice})
				beego.Info("send msg: " + currentPrice)
			} else {
				newPrice = lastPrice
			}
		}

		util.Redis.Set(util.Redis.Con(), "price_warn:13710227972", newPrice)
	}
}
