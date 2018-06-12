package services

import (
	"github.com/astaxie/beego"
	"net/url"
	"strings"
	"math/rand"
	"time"
)

type SmsService struct {
}

type SingleTpl struct {
	Ext    string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    Tel      `json:"tel"`
	Time   int64    `json:"time"`
	TplId  int64    `json:"tpl_id"`
}

type Tel struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}

// 单发短信
func (service *SmsService) SendSingle(nationCode string, mobile string, params []string, tplId int64) {
	appId := beego.AppConfig.String("sms::app_id")
	appKey := beego.AppConfig.String("sms::app_key")
	apiUrlConf := beego.AppConfig.String("sms::api_send_single")
	tplId, _ := beego.AppConfig.Int64("sms::tpl_price_warn")

	randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63n(1000)

	// 拼接URL
	apiUrlSli := strings.Split(apiUrlConf, "@")
	schemeHost := strings.Split(apiUrlSli[1], "://")
	apiUrl := url.URL{
		Scheme: schemeHost[0],
		Host:   schemeHost[1],
	}
	param := apiUrl.Query()
	param.Add("sdkappid", appId)
	param.Add("random", "")
	apiUrl.RawQuery = param.Encode()
}
