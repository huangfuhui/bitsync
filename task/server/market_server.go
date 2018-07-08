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

// 查询行情
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
		beego.Error("[行情服务]upgrade:", err)
		return
	}
	defer con.Close()

	huobiUsdt := strings.Split(beego.AppConfig.String("huobi::usdt_pair"), ",")
	dragonexUsdt := strings.Split(beego.AppConfig.String("dragonex::usdt_pair"), ",")

	for {
		_, _, err := con.ReadMessage()
		if err != nil {
			beego.Warn("[行情服务]read:", err)
			break
		}

		go func() {
			market := Market{
				Code: 200,
				Msg:  "",
			}

			for _, v := range huobiUsdt {
				key := "huobi:" + v + "usdt"
				redisCon := util.Redis.Con()
				priceStr, _ := util.Redis.Get(redisCon, key)
				price, _ := strconv.ParseFloat(priceStr, 64)
				market.Response = append(market.Response, SymbolPair{
					ExchangeId: 1,
					Symbol:     v + "/usdt",
					Price:      price,
				})
			}
			for _, v := range dragonexUsdt {
				key := "dragonex:" + v + "usdt"
				redisCon := util.Redis.Con()
				priceStr, _ := util.Redis.Get(redisCon, key)
				price, _ := strconv.ParseFloat(priceStr, 64)
				market.Response = append(market.Response, SymbolPair{
					ExchangeId: 2,
					Symbol:     v + "/usdt",
					Price:      price,
				})
			}

			err = con.WriteJSON(market)
			if err != nil {
				beego.Warn("[行情服务]write:", err)
			}
		}()

	}
}

func MarketServer() {
	http.HandleFunc("/market", market)
	err := http.ListenAndServe("localhost:8088", nil)
	if err != nil {
		beego.Error("[行情服务]", err)
	}
}
