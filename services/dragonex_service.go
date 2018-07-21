package services

import (
	"github.com/astaxie/beego"
	"net/http"
	"strings"
	"net/url"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"strconv"
	"bitsync/util"
)

type DragonexService struct {
}

var symbols = make(map[string]int64)

// 与龙交所建立连接，查询价格信息
func (service *DragonexService) WatchDragonex() {
	dragonexUrl := beego.AppConfig.String("dragonex::http_url")
	dragonexScheme := beego.AppConfig.String("dragonex::http_scheme")
	priceValidTime := beego.AppConfig.String("watch::price_valid_time")

	initRes := service.initSymbol(dragonexUrl, dragonexScheme)
	if !initRes {
		beego.Info("【龙交所】初始化交易对信息失败.")
		return
	}

	apiMarketReal := beego.AppConfig.String("dragonex::api_market_real")
	apiMarketRealSli := strings.Split(apiMarketReal, "@")

	defer beego.Info("【龙交所】价格信息同步关闭.")

	usdtPair := beego.AppConfig.String("dragonex::usdt_pair")

	beego.Info("【龙交所】价格信息同步开始.")
	for {
		priceMap := make(map[string]string)
		for k, v := range symbols {
			// 过滤不需要的货币
			symbolSli := strings.Split(k, "_")
			if !strings.Contains(usdtPair, symbolSli[0]) && symbolSli[1] != "usdt" {
				continue
			}

			client := &http.Client{}
			conUrl := url.URL{Scheme: dragonexScheme, Host: dragonexUrl, Path: apiMarketRealSli[1]}
			conValue := url.Values{}

			// 拼装请求的URL
			conValue.Add("symbol_id", strconv.FormatInt(v, 10))
			requestUrl := conUrl.String() + "?" + conValue.Encode()

			req, err := http.NewRequest(strings.ToUpper(apiMarketRealSli[0]), requestUrl, nil)

			if err != nil {
				beego.Error(err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				beego.Error(err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				beego.Error(err)
				return
			} else if http.StatusOK != resp.StatusCode {
				beego.Error(string(body))
				continue
			}

			// 解析价格
			symbolKey := strings.Replace(k, "_", "", -1)
			jsonData, _ := simplejson.NewJson(body)
			priceResMap, _ := jsonData.Get("data").GetIndex(0).Map()
			if len(priceResMap) == 0 {
				beego.Debug("【龙交所】价格解析失败,symbol:" + k)
				priceMap[symbolKey] = "0"
				continue
			}
			price, err := jsonData.Get("data").GetIndex(0).Get("close_price").String()
			if err != nil {
				beego.Error(err)
				return
			}
			priceMap[symbolKey] = price
		}

		// 更新本地价格信息
		for symbol, price := range priceMap {
			key := "dragonex:" + symbol
			con := util.Redis.Con();
			err := util.Redis.SetEx(con, key, price, priceValidTime)
			if err != nil {
				beego.Error(err)
				return
			}
		}

		beego.Debug("【龙交所】成功完成一次价格同步.")
	}
}

// 初始化交易对信息
func (service *DragonexService) initSymbol(dragonexUrl, dragonexScheme string) bool {
	client := &http.Client{}

	// 查询龙交所交易对信息
	apiAllSymbol := beego.AppConfig.String("dragonex::api_all_symbol")
	apiAllSymbolSli := strings.Split(apiAllSymbol, "@")
	conUrl := url.URL{Scheme: dragonexScheme, Host: dragonexUrl, Path: apiAllSymbolSli[1]}
	req, err := http.NewRequest(strings.ToUpper(apiAllSymbolSli[0]), conUrl.String(), nil)
	if err != nil {
		beego.Error(err)
		return false
	}

	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		beego.Error(err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || http.StatusOK != resp.StatusCode {
		beego.Info("【龙交所】查询货币信息失败.")
		beego.Error(err)
		return false
	}
	beego.Info("【龙交所】成功查询货币信息.")

	symbolData, err := simplejson.NewJson([]byte(body))
	if err != nil {
		beego.Error(err)
		return false
	}
	symbolArr, err := symbolData.Get("data").Array()
	if err != nil {
		beego.Error(err)
		return false
	}
	symbolQuantity := len(symbolArr)
	for i := 0; i < symbolQuantity; i++ {
		symbolName, err := symbolData.Get("data").GetIndex(i).Get("symbol").String()
		if err != nil {
			beego.Error(err)
			return false
		}
		symbolId, err := symbolData.Get("data").GetIndex(i).Get("symbol_id").Int64()
		if err != nil {
			beego.Error(err)
			return false
		}
		symbols[symbolName] = symbolId
	}

	return true
}
