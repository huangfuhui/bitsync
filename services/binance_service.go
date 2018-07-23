package services

import (
	"github.com/astaxie/beego"
	"net/url"
	"strings"
	"bitsync/util"
	"github.com/gorilla/websocket"
	"github.com/bitly/go-simplejson"
)

type BinanceService struct {
}

// 建立连接，并监控和解析通信数据
func (service *BinanceService) WatchBinance() {
	binanceScheme := beego.AppConfig.String("binance::ws_scheme")
	binanceUrl := beego.AppConfig.String("binance::ws_url")
	binancePath := beego.AppConfig.String("binance::ws_path")

	conUrl := url.URL{Scheme: binanceScheme, Host: binanceUrl, Path: binancePath}

	prices := make(chan string, 1024)
	go func() {
		priceValidTime := beego.AppConfig.String("watch::price_valid_time")
		for {
			select {
			case priceSlice := <-prices:
				priceSli := strings.Split(priceSlice, ":")

				key := "binance:" + strings.ToLower(priceSli[0])
				value := priceSli[1]

				con := util.Redis.Con()
				err := util.Redis.SetEx(con, key, value, priceValidTime)
				if err != nil {
					beego.Error(err)
				}
			}
		}
	}()

	for {
		// 1.建立websocket通信
		con, _, err := websocket.DefaultDialer.Dial(conUrl.String(), nil)
		if err != nil {
			beego.Error("【binance】dial: " + err.Error())
			continue
		} else {
			beego.Info("【binance】websocket通信建立.")
		}

		// 2.解析和更新本地价格信息
		for {
			_, response, err := con.ReadMessage()
			if err != nil {
				beego.Error("【binance】", err)
				beego.Info("【binance】websocket通信关闭.")
				con.Close()
				break
			}
			jsonData, err := simplejson.NewJson(response)
			if err != nil {
				beego.Error(err)
				beego.Info("【binance】websocket通信关闭.")
				con.Close()
				break
			} else if jsonData == nil {
				continue
			}

			resSli, _ := jsonData.Array()
			for k, _ := range resSli {
				symbol, _ := jsonData.GetIndex(k).Get("s").String()
				price, _ := jsonData.GetIndex(k).Get("c").String()

				if symbol != "" && price != "" {
					prices <- symbol + ":" + price
				}
			}
		}
		beego.Info("【binance】尝试重连websocket.")
	}
}
