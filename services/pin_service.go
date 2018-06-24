package services

import (
	"bitsync/util"
	"github.com/astaxie/beego"
)

const BUSINESS_CODE_PIN = "1000" // 登录验证码

type PinService struct {
}

func (service *PinService) Send(businessCode, handset string) {
	db, _ := beego.AppConfig.Int("redis_db_pin")
	redis := util.Cli{}
	redis.Select(db)

	key := "pin:" + businessCode + ":" + handset
	redis.Get(key)
}

func (service *PinService) Validate() {

}
