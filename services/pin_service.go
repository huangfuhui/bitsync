package services

import (
	"bitsync/util"
	"github.com/astaxie/beego"
	"errors"
	"strings"
	"strconv"
	"time"
)

const BUSINESS_CODE_PIN = "1001" // 登录验证码

var (
	ERR_EMPTY_BUSINESS_CODE = errors.New("业务码不能为空")
	ERR_EMPTY_HANDSET       = errors.New("手机号码不能为空")
	ERR_LIMIT_REQUEST       = errors.New("请求频率过高")
)

type PinService struct {
}

// 发送验证码
func (service *PinService) Send(businessCode, handset string) (string, error) {
	if businessCode == "" {
		return "", ERR_EMPTY_BUSINESS_CODE
	} else if handset == "" {
		return "", ERR_EMPTY_HANDSET
	}

	db, _ := beego.AppConfig.Int("redis_db_pin")
	redis := util.Cli{}
	redis.Select(db)

	key := "pin:" + businessCode + ":" + handset
	listLength, _ := redis.Llen(key)
	lastRequest, _ := redis.Lindex(key, "0")

	// 上一次请求的信息
	var lastRequestTime int64 = 0
	if lastRequest != "" {
		lastRequestRes := strings.Split(lastRequest, "|")
		lastRequestTime, _ = strconv.ParseInt(lastRequestRes[0], 0, 64)
	}

	// 最旧的请求信息
	var firstRequestTime int64 = 0
	if listLength > 0 {
		firstRequest, _ := redis.Lindex(key, strconv.FormatInt(int64(listLength-1), 10))
		firstRequestRes := strings.Split(firstRequest, "|")
		firstRequestTime, _ = strconv.ParseInt(firstRequestRes[0], 0, 64)
	}

	// 频率限制，每60s允许一次请求，每24小时允许10次请求
	now := time.Now().Unix()
	if now-lastRequestTime < 60 {
		return "", ERR_LIMIT_REQUEST
	}
	if listLength > 10 && lastRequestTime-firstRequestTime < 3600*24 {
		return "", ERR_LIMIT_REQUEST
	}

	// 发送验证码
	sms := SmsService{}
	random := util.Random{}
	pin := strconv.FormatInt(random.Rand(1000, 9999), 10)
	tplId, _ := beego.AppConfig.Int64("tpl_verify_code")
	err := sms.SendSingle("86", handset, []string{pin}, tplId)
	if err != nil {
		return "", err
	}

	// 记录请求信息
	value := strconv.FormatInt(now, 10) + "|" + pin
	redis.Lpush(key, value)
	redis.Ltrim(key, "0", "9")
	// 更新验证码记录有效期为一天
	redis.SetEx(key, "86400")

	return pin, nil
}

func (service *PinService) Validate() {

}
