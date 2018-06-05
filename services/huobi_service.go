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
	"strconv"
)

// 与火币建立连接，并监控和解析通信数据
func Watch() {
	huobiScheme := beego.AppConfig.String("huobi_ws_scheme")
	huobiUrl := beego.AppConfig.String("huobi_ws_url")
	huobiPath := beego.AppConfig.String("huobi_ws_path")

	addr := flag.String("addr", huobiUrl, "http service address")
	flag.Parse()
	conUrl := url.URL{Scheme: huobiScheme, Host: *addr, Path: huobiPath}

	// 建立websocket通信
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

	done := make(chan struct{}, 1)

	// 订阅价格信息
	// con.WriteMessage(websocket.TextMessage, []byte(""))

	go func() {
		defer close(done)
		for {
			_, gzipData, err := con.ReadMessage()
			if err != nil {
				beego.Error(err)

				return
			}

			jsonData, err := parseGzip(gzipData)
			if err != nil {
				beego.Error(err)

				return
			}

			// 解析json数据
			data, err := simplejson.NewJson(jsonData)
			if err != nil {
				beego.Error(err)

				return
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
			}
		}
	}()



	for {
		select {
		case <-done:
			return
		}
	}
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
