package services

import (
	"strings"
	"github.com/bitly/go-simplejson"
	"github.com/astaxie/beego"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"bitsync/util"
)

type BithumbService struct {
}

// 建立连接，查询价格信息
func (service *BithumbService) WatchBithumb() {
	bithumbUrl := beego.AppConfig.String("bithumb::http_url")
	bithumbScheme := beego.AppConfig.String("bithumb::http_scheme")
	priceValidTime := beego.AppConfig.String("watch::price_valid_time")

	apiTicker := beego.AppConfig.String("bithumb::api_ticker")
	apiTickerSli := strings.Split(apiTicker, "@")

	defer beego.Info("【bithumb】价格信息同步关闭.")

	krwPair := beego.AppConfig.String("bithumb::krw_pair")
	krwPairSli := strings.Split(krwPair, ",")

	beego.Info("【bithumb】价格信息同步开始.")
	for {
		time.Sleep(time.Millisecond * 200)

		priceMap := make(map[string]string)
		client := &http.Client{}
		conUrl := url.URL{Scheme: bithumbScheme, Host: bithumbUrl, Path: apiTickerSli[1]}

		req, err := http.NewRequest(strings.ToUpper(apiTickerSli[0]), conUrl.String(), nil)
		if err != nil {
			beego.Error("【bithumb】", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			beego.Error("【bithumb】", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			beego.Error("【bithumb】", err)
			return
		} else if http.StatusOK != resp.StatusCode {
			beego.Error("【bithumb】", string(body))
			continue
		}

		jsonData, _ := simplejson.NewJson(body)
		status, _ := jsonData.Get("status").String()
		if status != "0000" {
			beego.Error("【bithumb】", string(body))
			continue
		}

		// 解析价格
		for _, v := range krwPairSli {
			price, _ := jsonData.Get("data").Get(strings.ToUpper(v)).Get("buy_price").String()
			if price != "" {
				priceMap[v+"krw"] = price
			}
		}

		// 更新本地价格信息
		for symbol, price := range priceMap {
			key := "bithumb:" + symbol
			con := util.Redis.Con();
			err := util.Redis.SetEx(con, key, price, priceValidTime)
			if err != nil {
				beego.Error("【bithumb】", err)
				return
			}
		}

		beego.Debug("【bithumb】成功完成一次价格同步.")
	}
}
