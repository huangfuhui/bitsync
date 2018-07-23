package services

import (
	"github.com/astaxie/beego"
	"net/url"
	"github.com/gorilla/websocket"
	"strings"
	"github.com/bitly/go-simplejson"
	"bitsync/util"
	"time"
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
					beego.Error("【okex】", err)
				}
			}
		}
	}()

	for {
		// 1.建立websocket通信
		con, _, err := websocket.DefaultDialer.Dial(conUrl.String(), nil)
		if err != nil {
			beego.Error("【okex】dial: ", err)
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
			beego.Error("【okex】", err)
			beego.Info("【okex】websocket通信关闭.")
			con.Close()
			return
		}
		beego.Info("【okex】价格订阅成功.")

		timer := time.NewTimer(time.Second * 10)
		timeTag := time.Now().Unix()

		// 3.解析价格
		for {
			select {
			case <-timer.C:
				go func() {
					err = con.WriteMessage(websocket.TextMessage, []byte(`{'event':'ping'}`))
					if err != nil {
						beego.Error("【okex】", err)
						beego.Info("【okex】websocket通信关闭.")
						con.Close()
					}
					timeTag = time.Now().Unix()
				}()
			default:
				if time.Now().Unix()-timeTag > 30 {
					beego.Info("【okex】心跳检活超时, 尝试重连.")
					con.Close()
					goto retryConnect
				}

				_, jsonData, err := con.ReadMessage()
				if err != nil {
					beego.Error("【okex】", err)
					beego.Info("【okex】websocket通信关闭.")
					con.Close()
					break
				}
				beego.Debug("【okex】", string(jsonData))

				data, err := simplejson.NewJson(jsonData)
				if err != nil {
					beego.Error("【okex】", err)
					goto retryConnect
				}

				// 检测心跳
				pong, _ := data.Get("event").String()
				if pong == "pong" {
					timer = time.NewTimer(time.Second * 10)
					beego.Debug("【okex】心跳检活成功.")
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
		}
	retryConnect:
		beego.Info("【okex】尝试重连websocket.")
	}
}
