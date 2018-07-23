package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"github.com/astaxie/beego"
	"strings"
	"bitsync/util"
	"strconv"
)

type Market struct {
	Code     int          `json:"code"`
	Response []SymbolPair `json:"response"`
	Msg      string       `json:"msg"`
}

type SymbolPair struct {
	ExchangeId int     `json:"exchange_id"`
	Symbol     string  `json:"symbol"`
	Price      float64 `json:"price"`
}

type SymbolRequest struct {
	ExchangeId int    `json:"exchange_id"`
	Symbol     string `json:"symbol"`
}

// 查询全部行情
func market(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 20 * time.Second,
		ReadBufferSize:   128,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		beego.Error("【行情服务】upgrade:", err)
		return
	}
	defer con.Close()

	exchange := map[string]int{
		"huobi":    1,
		"dragonex": 2,
		"okex":     3,
		"binance":  4,
		"gate":     5,
		"bithumb":  6,
	}
	marketConf := strings.Split(beego.AppConfig.String("watch::market"), ";")
	beego.Debug("【行情服务】行情配置", marketConf)

	for {
		_, _, err := con.ReadMessage()
		if err != nil {
			beego.Warn("【行情服务】read:", err)
			break
		}

		market := Market{
			Code: http.StatusOK,
			Msg:  "",
		}

		for _, v := range marketConf {
			exchangeSymbols := strings.Split(v, ":")
			exchangeId := exchange[exchangeSymbols[0]]
			for _, symbol := range strings.Split(exchangeSymbols[1], ",") {
				pairs := strings.Split(beego.AppConfig.String(exchangeSymbols[0]+"::"+symbol+"_pair"), ",")

				for _, sp := range pairs {
					key := exchangeSymbols[0] + ":" + sp + symbol

					redisCon := util.Redis.Con()
					priceStr, _ := util.Redis.Get(redisCon, key)
					price, _ := strconv.ParseFloat(priceStr, 64)

					market.Response = append(market.Response, SymbolPair{
						ExchangeId: exchangeId,
						Symbol:     sp + "/" + symbol,
						Price:      price,
					})
				}
			}
		}

		err = con.WriteJSON(market)
		if err != nil {
			beego.Warn("【行情服务】write:", err)
		}
	}
}

// 查询交易对价格
func symbol(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 20 * time.Second,
		ReadBufferSize:   128,
		WriteBufferSize:  512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		beego.Error("【行情服务】upgrade:", err)
		return
	}
	defer con.Close()

	for {
		symbol := SymbolRequest{}
		err = con.ReadJSON(&symbol)
		if err != nil {
			beego.Warn("【行情服务】read:", err)
			return
		}

		market := Market{
			Code: http.StatusBadRequest,
			Msg:  "交易对不存在",
		}

		symbolPairSli := strings.Split(symbol.Symbol, "/")

		exchange := map[int]string{
			1: "huobi",
			2: "dragonex",
			3: "okex",
			4: "binance",
			5: "gate",
			6: "bithumb",
		}
		if symbol.ExchangeId <= 6 && symbol.ExchangeId >= 1 {
			SymbolPairs := beego.AppConfig.String(exchange[symbol.ExchangeId] + "::" + symbolPairSli[1] + "_pair")
			if strings.Contains(SymbolPairs, symbolPairSli[0]) {
				key := exchange[symbol.ExchangeId] + ":" + strings.Replace(symbol.Symbol, "/", "", -1)
				redisCon := util.Redis.Con()
				priceStr, _ := util.Redis.Get(redisCon, key)
				price, _ := strconv.ParseFloat(priceStr, 64)

				market.Code = http.StatusOK
				market.Response = append(market.Response, SymbolPair{
					ExchangeId: 1,
					Symbol:     symbol.Symbol,
					Price:      price,
				})
				market.Msg = ""
			}
		} else {
			market.Msg = "交易所ID非法"
		}

		err = con.WriteJSON(market)
		if err != nil {
			beego.Warn("【行情服务】write:", err)
		}
	}
}

func MarketServer() {
	http.HandleFunc("/market", market)
	http.HandleFunc("/symbol", symbol)
	err := http.ListenAndServe("localhost:8088", nil)
	if err != nil {
		beego.Error("【行情服务】", err)
	}
}
