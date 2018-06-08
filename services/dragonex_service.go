package services

import (
	"github.com/astaxie/beego"
	"net/http"
	"strings"
	"net/url"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"strconv"
)

var symbols = make(map[string]int64)

// 与龙交所建立连接，查询价格信息
func WatchDragonex() {
	dragonexUrl := beego.AppConfig.String("dragonex::http_url")
	dragonexScheme := beego.AppConfig.String("dragonex::http_scheme")

	initSymbol(dragonexUrl, dragonexScheme)

	apiMarketReal := beego.AppConfig.String("dragonex::api_market_real")
	apiMarketRealSli := strings.Split(apiMarketReal, "@")

	for {
		for _, v := range symbols {
			client := &http.Client{}
			conUrl := url.URL{Scheme: dragonexScheme, Host: dragonexUrl, Path: apiMarketRealSli[1]}

			req, err := http.NewRequest(strings.ToUpper(apiMarketRealSli[0]), conUrl.String(), nil)

			if err != nil {
				beego.Error(err)
				return
			}

			req.Form.Add("symbol_id", strconv.FormatInt(v, 10))

			resp, err := client.Do(req)
			if err != nil {
				resp.Body.Close()
				beego.Error(err)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil || http.StatusOK != resp.StatusCode {
				beego.Info("【龙交所】查询货币信息失败.")
				beego.Error(err)
				return
			}

			beego.Info(string(body))
		}
	}
}

// 初始化交易对信息
func initSymbol(dragonexUrl, dragonexScheme string) {
	client := &http.Client{}

	// 查询龙交所交易对信息
	apiAllSymbol := beego.AppConfig.String("dragonex::api_all_symbol")
	apiAllSymbolSli := strings.Split(apiAllSymbol, "@")
	conUrl := url.URL{Scheme: dragonexScheme, Host: dragonexUrl, Path: apiAllSymbolSli[1]}
	req, err := http.NewRequest(strings.ToUpper(apiAllSymbolSli[0]), conUrl.String(), nil)
	if err != nil {
		beego.Error(err)
		return
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || http.StatusOK != resp.StatusCode {
		beego.Info("【龙交所】查询货币信息失败.")
		beego.Error(err)
		return
	}
	beego.Info("【龙交所】成功查询货币信息.")
	beego.Debug("【龙交所】" + string(body))

	// TODO:更新本地交易对信息
	symbolData, err := simplejson.NewJson([]byte(body))
	if err != nil {
		beego.Error(err)
		return
	}
	symbolArr, err := symbolData.Get("data").Array()
	if err != nil {
		beego.Error(err)
		return
	}
	symbolQuantity := len(symbolArr)
	for i := 0; i < symbolQuantity; i++ {
		symbolName, err := symbolData.Get("data").GetIndex(i).Get("symbol").String()
		if err != nil {
			beego.Error(err)
			return
		}
		symbolId, err := symbolData.Get("data").GetIndex(i).Get("symbol_id").Int64()
		if err != nil {
			beego.Error(err)
			return
		}
		symbols[symbolName] = symbolId
	}
}
