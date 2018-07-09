package coin

import "bitsync/object"

const (
	EXCHANGE_HUOBI     = 1
	EXCHANGE_DRAGEONEX = 2
)

type Exchange struct {
	object.Base
	ExchangeId      int
	NameCn          string
	NameEn          string
	OfficialWebSite string
	Logo            string
}
