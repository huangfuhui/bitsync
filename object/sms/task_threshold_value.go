package sms

import "bitsync/object"

type TaskThresholdValue struct {
	object.Base
	UID            int `orm:"column(UID)"`
	CoinAId        int
	CoinBId        int
	SymbolPair     string
	ExchangeId     int
	ThresholdValue string
	BaseValue      string
	Deviation      int
}
