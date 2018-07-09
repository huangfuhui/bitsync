package sms

import "bitsync/object"

const (
	DEVIATION_GT = 1
	DEVIATION_LT = 2
)

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
