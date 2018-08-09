package services

import (
	"net/url"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"bytes"
	"encoding/binary"
	"compress/gzip"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"strings"
	"strconv"
	"encoding/json"
	"bitsync/util"
)

type HuobiService struct {
}

// K线
type KLine struct {
	Ch string `json:"ch"`
	Ts int64  `json:"ts"`
	Tick struct {
		Id     int64   `json:"id"`
		Amount float64 `json:"amount"`
		Count  float64 `json:"count"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Vol    float64 `json:"vol"`
	}
}

// 成功订阅
type SubSuccess struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Subbed string `json:"subbed"`
	Ts     int64  `json:"ts"`
}

// 与火币建立连接，并监控和解析通信数据
func (service *HuobiService) WatchHuobi() {
	huobiScheme := beego.AppConfig.String("huobi::ws_scheme")
	huobiUrl := beego.AppConfig.String("huobi::ws_url")
	huobiPath := beego.AppConfig.String("huobi::ws_path")

	conUrl := url.URL{Scheme: huobiScheme, Host: huobiUrl, Path: huobiPath}

	prices := make(chan string, 1024)
	go func() {
		priceValidTime := beego.AppConfig.String("watch::price_valid_time")
		for {
			select {
			case priceSlice := <-prices:
				priceSli := strings.Split(priceSlice, ":")
				subbedSli := strings.Split(priceSli[0], ".")

				key := "huobi:" + subbedSli[1]
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
			beego.Error("【火币】dial: " + err.Error())
			continue
		} else {
			beego.Info("【火币】websocket通信建立.")
		}

		// 2.价格信息订阅
		beego.Info("【火币】开始价格订阅.")
		exchangeSymbols := strings.Split(beego.AppConfig.String("huobi::exchange_symbol"), ",")
		for _, v := range exchangeSymbols {
			pairs := beego.AppConfig.String("huobi::" + v + "_pair")
			pairsSlice := strings.Split(pairs, ",")
			for _, pair := range pairsSlice {
				// 发起订阅
				symbol := pair + v
				err := con.WriteMessage(websocket.TextMessage, []byte(`{"sub":"market.`+symbol+`.kline.1min","id":"`+symbol+`"}`))
				if err != nil {
					beego.Error(err.Error())
					beego.Info("【火币】websocket通信关闭.")
					con.Close()
					return
				}

				// 获取订阅结果
				subRes, err := service.onlySubResult(con)
				if err != nil {
					beego.Error(err)
					beego.Info("【火币】websocket通信关闭.")
					con.Close()
					return
				}

				if subRes.Status == "ok" {
					beego.Debug("【火币】" + subRes.Subbed + "订阅成功.")
				} else if subRes.Status == "error" {
					beego.Error("【火币】价格订阅失败,subbed: " + subRes.Subbed)
				}
			}
		}
		beego.Info("【火币】价格订阅成功.")

		// 3.解析和更新本地价格信息
		for {
			jsonData, parseData, err := service.parseResponse(con)
			if err != nil {
				beego.Error(err)
				beego.Info("【火币】websocket通信关闭.")
				con.Close()
				break
			} else if parseData == nil {
				continue
			}

			// 解析价格信息
			kLine := KLine{}
			err = json.Unmarshal(jsonData, &kLine)
			if err != nil {
				beego.Error(err)
				continue
			}

			prices <- kLine.Ch + ":" + strconv.FormatFloat(kLine.Tick.Close, 'f', 4, 64)
		}
		beego.Info("【火币】尝试重连websocket.")
	}
}

// 仅解析订阅结果
func (service *HuobiService) onlySubResult(con *websocket.Conn) (SubSuccess, error) {
	var jsonData []byte
	var err error

	subSuccess := SubSuccess{}
	for {
		jsonData, _, err = service.parseResponse(con)
		err = json.Unmarshal(jsonData, &subSuccess)
		if err != nil {
			return subSuccess, err
		}

		if subSuccess.Subbed != "" {
			break
		}
	}

	return subSuccess, nil
}

// 读取解析websocket响应
func (service *HuobiService) parseResponse(con *websocket.Conn) ([]byte, *simplejson.Json, error) {
	// 读取数据
	_, gzipData, err := con.ReadMessage()
	if err != nil {
		beego.Error(err)
		return nil, nil, err
	}

	// 解压数据
	jsonData, err := service.parseGzip(gzipData)
	if err != nil {
		beego.Error(err)
		return jsonData, nil, err
	}

	// 解析json数据
	data, err := simplejson.NewJson(jsonData)
	if err != nil {
		beego.Error(err)
		return jsonData, nil, err
	}

	// 回复心跳检测
	ping, _ := data.Get("ping").Int64()
	if ping != 0 {
		beego.Debug("【火币】收到心跳检测: " + string(jsonData))
		pong := `{"pong": ` + strconv.FormatInt(ping, 10) + `}`
		err := con.WriteMessage(websocket.TextMessage, []byte(pong))
		if err != nil {
			beego.Debug("【火币】回复心跳检测失败," + err.Error())
		} else {
			beego.Debug("【火币】回复心跳检测成功.")
		}
		return jsonData, nil, nil
	}
	return jsonData, data, nil
}

// 解析Gzip数据
func (service *HuobiService) parseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)

	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		res, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
