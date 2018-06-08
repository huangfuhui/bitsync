package services

import (
	"flag"
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
func WatchHuobi() {
	huobiScheme := beego.AppConfig.String("huobi::ws_scheme")
	huobiUrl := beego.AppConfig.String("huobi::ws_url")
	huobiPath := beego.AppConfig.String("huobi::ws_path")

	addr := flag.String("addr", huobiUrl, "http service address")
	flag.Parse()
	conUrl := url.URL{Scheme: huobiScheme, Host: *addr, Path: huobiPath}

	// 1.建立websocket通信
	con, _, err := websocket.DefaultDialer.Dial(conUrl.String(), nil)
	if err != nil {
		beego.Error("【火币】dial: " + err.Error())
		return
	} else {
		beego.Info("【火币】websocket通信建立.")
	}
	defer func() {
		beego.Info("【火币】websocket通信关闭.")
	}()
	defer con.Close()

	// 2.价格信息订阅
	beego.Info("【火币】开始价格订阅.")
	pricePairs := beego.AppConfig.String("huobi::price_pairs")
	priceSlice := strings.Split(pricePairs, ",")
	for _, v := range priceSlice {
		// 发起订阅
		err := con.WriteMessage(websocket.TextMessage, []byte(`{"sub":"market.`+v+`.kline.1min","id":"`+v+`"}`))
		if err != nil {
			beego.Error(err.Error())
			return
		}

		// 获取订阅结果
		subRes, err := onlySubResult(con)
		if err != nil {
			beego.Error(err)
			return
		}

		if subRes.Status == "ok" {
			beego.Debug("【火币】" + subRes.Subbed + "订阅成功.")
			continue
		} else if subRes.Status == "error" {
			beego.Error("【火币】价格订阅失败,subbed: " + subRes.Subbed)
			return
		}
	}
	beego.Info("【火币】价格订阅成功.")

	prices := make(chan string, 1024)
	priceValidTime := beego.AppConfig.String("watch::price_valid_time")
	go func(priceValidTime string) {
		for {
			select {
			case priceSlice := <-prices:
				priceSli := strings.Split(priceSlice, ":")
				subbedSli := strings.Split(priceSli[0], ".")

				key := "huobi:" + subbedSli[1]
				value := priceSli[1]

				err := util.Redis.SetEx(key, value, priceValidTime)
				if err != nil {
					beego.Error(err)
				}
			}
		}
	}(priceValidTime)

	// 3.解析和更新本地价格信息
	for {
		jsonData, parseData, err := parseResponse(con)
		if err != nil {
			beego.Error(err)
			return
		} else if parseData == nil {
			continue
		}

		// 解析价格信息
		kLine := KLine{}
		err = json.Unmarshal(jsonData, &kLine)
		if err != nil {
			beego.Error(err)
			return
		}

		prices <- kLine.Ch + ":" + strconv.FormatFloat(kLine.Tick.Close, 'f', 4, 64)
	}
}

// 仅解析订阅结果
func onlySubResult(con *websocket.Conn) (SubSuccess, error) {
	var jsonData []byte
	var err error

	subSuccess := SubSuccess{}
	for {
		jsonData, _, err = parseResponse(con)
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
func parseResponse(con *websocket.Conn) ([]byte, *simplejson.Json, error) {
	// 读取数据
	_, gzipData, err := con.ReadMessage()
	if err != nil {
		beego.Error(err)
		return nil, nil, err
	}

	// 解压数据
	jsonData, err := parseGzip(gzipData)
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
		beego.Info("【火币】收到心跳检测: " + string(jsonData))
		pong := `{"pong": ` + strconv.FormatInt(ping, 10) + `}`
		err := con.WriteMessage(websocket.TextMessage, []byte(pong))
		if err != nil {
			beego.Info("【火币】回复心跳检测失败," + err.Error())
		} else {
			beego.Info("【火币】回复心跳检测成功.")
		}
		return jsonData, nil, nil
	}
	return jsonData, data, nil
}

// 解析Gzip数据
func parseGzip(data []byte) ([]byte, error) {
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
