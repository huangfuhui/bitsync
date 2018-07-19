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

type GateService struct {
}

// 建立连接，并监控和解析通信数据
func (service *GateService) WatchGate() {
	gateScheme := beego.AppConfig.String("gate::ws_scheme")
	gateUrl := beego.AppConfig.String("gate::ws_url")
	gatePath := beego.AppConfig.String("gate::ws_path")

	conUrl := url.URL{Scheme: gateScheme, Host: gateUrl, Path: gatePath}

	// 监听价格变动和更新本地价格信息
	prices := make(chan string, 1024)
	go func() {
		priceValidTime := beego.AppConfig.String("watch::price_valid_time")
		for {
			select {
			case priceSlice := <-prices:
				priceSli := strings.Split(priceSlice, ":")

				key := "gate:" + priceSli[0]
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
			beego.Error("【gate】dial: " + err.Error())
			continue
		} else {
			beego.Info("【gate】websocket通信建立.")
		}

		// 2.价格信息订阅
		beego.Info("【gate】开始价格订阅.")
		usdtPairs := beego.AppConfig.String("gate::usdt_pair")
		usdtSlice := strings.Split(usdtPairs, ",")
		var subPairs []string
		for _, v := range usdtSlice {
			subPairs = append(subPairs, `"`+strings.ToUpper(v)+`_USDT"`)
		}
		err = con.WriteMessage(websocket.TextMessage, []byte(`{"id":888,"method":"ticker.subscribe", "params":[`+strings.Join(subPairs, ",")+`]}`))
		if err != nil {
			beego.Error(err.Error())
			beego.Info("【gate】websocket通信关闭.")
			con.Close()
			return
		}

		// 3.解析订阅结果
		_, subData, err := con.ReadMessage()
		if err != nil {
			beego.Error(err)
			return
		}
		subRes, err := simplejson.NewJson(subData)
		if err != nil {
			beego.Error(err)
			return
		}
		subStatus, _ := subRes.Get("result").Get("status").String()
		if subStatus != "success" {
			beego.Error("【gate】", string(subData))
			continue
		}
		beego.Info("【gate】价格订阅成功.")

		// 心跳检测
		timer := time.NewTimer(time.Second * 2)
		timeTag := time.Now().Unix()

		// 3.解析价格
		for {
			select {
			case <-timer.C:
				go func() {
					err = con.WriteMessage(websocket.TextMessage, []byte(`{"id":999, "method":"server.ping", "params":[]}`))
					if err != nil {
						beego.Error(err.Error())
						beego.Info("【gate】websocket通信关闭.")
						con.Close()
					}
					timeTag = time.Now().Unix()
				}()
			default:
				if time.Now().Unix()-timeTag > 30 {
					beego.Info("【gate】心跳检活超时, 尝试重连.")
					con.Close()
					goto retryConnect
				}

				_, jsonData, err := con.ReadMessage()
				if err != nil {
					beego.Error(err)
					beego.Info("【gate】websocket通信关闭.")
					con.Close()
					goto retryConnect
				}
				beego.Debug("【gate】", string(jsonData))

				timeTag = time.Now().Unix()

				data, err := simplejson.NewJson(jsonData)
				if err != nil {
					beego.Error(err)
					continue
				}

				// 检测心跳
				id, _ := data.Get("id").Int()
				pong, _ := data.Get("result").String()
				if id == 999 && pong == "pong" {
					timer = time.NewTimer(time.Second * 2)
					beego.Debug("【gate】心跳检活成功.")
					continue
				}

				// 解析价格
				method, _ := data.Get("method").String()
				if method == "ticker.update" {
					symbol, _ := data.Get("params").GetIndex(0).String()
					symbolPair := strings.ToLower(strings.Replace(symbol, "_", "", -1))

					price, _ := data.Get("params").GetIndex(1).Get("last").String()
					prices <- symbolPair + ":" + price
				}
			}
		}
	retryConnect:
	}
}
