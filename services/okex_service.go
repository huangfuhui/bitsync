package services

import (
	"github.com/astaxie/beego"
	"net/url"
	"github.com/gorilla/websocket"
	"strings"
	"github.com/bitly/go-simplejson"
	"bitsync/util"
)

type OkexService struct {
}

func (service *OkexService) WatchOkex() {
	okexScheme := beego.AppConfig.String("okex::ws_scheme")
	okexUrl := beego.AppConfig.String("okex::ws_url")
	okexPath := beego.AppConfig.String("okex::ws_path")

	conUrl := url.URL{Scheme: okexScheme, Host: okexUrl, Path: okexPath}

	// 监听价格变动和更新本地价格信息
	prices := make(chan string, 1024)
	go func() {
		priceValidTime := beego.AppConfig.String("watch::price_valid_time")
		for {
			select {
			case priceSlice := <-prices:
				priceSli := strings.Split(priceSlice, ":")

				key := "okex:" + priceSli[0]
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
			beego.Error("【okex】dial: " + err.Error())
			continue
		} else {
			beego.Info("【okex】websocket通信建立.")
		}

		// 2.批量订阅价格
		beego.Info("【okex】开始价格订阅.")
		var subStrSli []string
		usdtPairs := beego.AppConfig.String("okex::usdt_pair")
		usdtSlice := strings.Split(usdtPairs, ",")
		for _, v := range usdtSlice {
			subStrSli = append(subStrSli, "{'event':'addChannel','channel':'ok_sub_spot_"+v+"_usdt_ticker'}")
		}

		err = con.WriteMessage(websocket.TextMessage, []byte("["+strings.Join(subStrSli, ",")+"]"))
		if err != nil {
			beego.Error(err.Error())
			beego.Info("【okex】websocket通信关闭.")
			con.Close()
			return
		}
		beego.Info("【okex】价格订阅成功.")

		// 3.解析价格
		for {
			_, jsonData, err := con.ReadMessage()
			if err != nil {
				beego.Error(err)
				beego.Info("【okex】websocket通信关闭.")
				con.Close()
				break
			}
			beego.Debug("【okex】", string(jsonData))

			data, err := simplejson.NewJson(jsonData)
			if err != nil {
				beego.Error(err)
				continue
			}

			channel, _ := data.GetIndex(0).Get("channel").String()
			price, _ := data.GetIndex(0).Get("data").Get("last").String()
			if channel == "" || price == "" {
				continue
			}

			channelSli := strings.Split(channel, "_")
			symbolPair := channelSli[3] + channelSli[4]

			prices <- symbolPair + ":" + price
		}
		beego.Info("【okex】尝试重连websocket.")
	}
}
