package coin

import "bitsync/object"

const (
	EXCHANGE_HUOBI     = 1
	EXCHANGE_DRAGEONEX = 2
	EXCHANGE_OKEX      = 3
	EXCHANGE_BINANCE   = 4
	EXCHANGE_GATE      = 5
	EXCHANGE_BITHUMB   = 6
)

var EXCHANGE = map[int]string{
	1: "huobi",
	2: "dragonex",
	3: "okex",
	4: "binance",
	5: "gate",
	6: "bithumb",
}

type Exchange struct {
	object.Base
	ExchangeId      int
	NameCn          string
	NameEn          string
	OfficialWebSite string
	Logo            string
}
