package server

import (
	"net/http"
	"flag"
	"github.com/gorilla/websocket"
	"time"
	"github.com/astaxie/beego"
	"strings"
)

type Market struct {
	Code     int          `json:"code"`
	Response []SymbolPair `json:"response"`
	Msg      string       `json:"msg"`
}

type SymbolPair struct {
	ExchangeId int    `json:"exchange_id"`
	Symbol     string `json:"symbol"`
	Price      int    `json:"price"`
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
			for _, v := range huobiUsdt {
				
			}

			err = con.WriteMessage(websocket.TextMessage, []byte("hello word"))
			if err != nil {
				beego.Warn("[行情服务]write:", err)
			}
		}()

	}
}

func MarketServer() {
	addr := flag.String("addr", "localhost:8088", "http service address")
	flag.Parse()
	http.HandleFunc("/market", market)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		beego.Error("[行情服务]", err)
	}
}
