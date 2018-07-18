package services

import (
	"github.com/astaxie/beego"
	"flag"
	"net/url"
	"github.com/gorilla/websocket"
	"strings"
)

type OkexService struct {
}

func (service *OkexService) WatchOkex() {
	okexScheme := beego.AppConfig.String("okex::ws_scheme")
	okexUrl := beego.AppConfig.String("okex::ws_url")
	okexPath := beego.AppConfig.String("okex::ws_path")

	addr := flag.String("addr", okexUrl, "http service address")
	flag.Parse()
	conUrl := url.URL{Scheme: okexScheme, Host: *addr, Path: okexPath}

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
			_, data, err := con.ReadMessage()
			if err != nil {
				beego.Error(err)
				return
			}
			beego.Info(data)
		}
	}
}
