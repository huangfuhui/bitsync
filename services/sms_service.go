package services

import (
	"github.com/astaxie/beego"
	"net/url"
	"strings"
	"bitsync/util"
	"strconv"
	"net/http"
	"time"
	"crypto/sha256"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"errors"
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

type SingleTplRes struct {
	Result int64  `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int64  `json:"fee"`
	Sid    string `json:"sid"`
}

// 单发短信
func (service *SmsService) SendSingle(nationCode string, mobile string, params []string) error {
	appSign := beego.AppConfig.String("sms::app_sign")
	appId := beego.AppConfig.String("sms::app_id")
	appKey := beego.AppConfig.String("sms::app_key")
	apiUrlConf := beego.AppConfig.String("sms::api_send_single")
	tplId, _ := beego.AppConfig.Int64("sms::tpl_price_warn")

	random := util.Random{}.Rand(1000, 9999)
	times := time.Now().Unix()

	// 拼接URL
	apiUrlSli := strings.Split(apiUrlConf, "@")
	schemeHost := strings.Split(apiUrlSli[1], "://")
	apiUrl := url.URL{
		Scheme: schemeHost[0],
		Host:   schemeHost[1],
	}
	param := apiUrl.Query()
	param.Add("sdkappid", appId)
	param.Add("random", strconv.FormatInt(random, 10))
	apiUrl.RawQuery = param.Encode()

	// 计算签名
	sig := service.sig(appKey, strconv.FormatInt(random, 10), strconv.FormatInt(times, 10), mobile)

	// 组装json请求数据
	request := SingleTpl{
		Ext:    "",
		Extend: "",
		Params: params,
		Sig:    sig,
		Sign:   appSign,
		Tel: Tel{
			Mobile:     mobile,
			Nationcode: nationCode,
		},
		Time:  times,
		TplId: tplId,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(strings.ToUpper(apiUrlSli[0]), apiUrl.String(), strings.NewReader(string(jsonRequest)))
	if err != nil {
		return err
	}

	// 发起请求
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析结果
	res, err := simplejson.NewJson(response)
	if err != nil {
		return err
	}
	result, _ := res.Get("result").Int64()
	errMsg, _ := res.Get("errmsg").String()
	if result != 0 {
		return errors.New(errMsg)
	}

	return nil
}

// 计算签名
func (service *SmsService) sig(appKey, random, times, mobile string) string {
	value := url.Values{}
	value.Add("appkey", appKey)
	value.Add("random", random)
	value.Add("time", times)
	value.Add("mobile", mobile)

	sh := sha256.New()
	sh.Write([]byte(value.Encode()))
	return fmt.Sprintf("%x", sh.Sum(nil))
}
