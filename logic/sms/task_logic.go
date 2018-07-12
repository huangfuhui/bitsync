package sms

import (
	"bitsync/logic"
	"bitsync/object/sms"
	"strings"
	"bitsync/models/coin"
	"bitsync/util"
	"github.com/astaxie/beego"
	"strconv"
	model "bitsync/models/sms"
	coinObj "bitsync/object/coin"
)

type TaskLogic struct {
	logic.BaseLogic
}

// 添加预警任务
func (l *TaskLogic) Add(taskType, exchangeId int, symbolPair string, deviation int, value float64) {
	UID := l.GetUID()

	symbolPairSli := strings.Split(symbolPair, "_")
	coinModel := coin.CoinModel{}
	coinA, _ := coinModel.GetByName(symbolPairSli[0])
	coinB, _ := coinModel.GetByName(symbolPairSli[1])
	if coinA.Id == 0 || coinB.Id == 0 {
		l.BadRequest("交易对不存在")
		return
	}

	key := ""
	if exchangeId == coinObj.EXCHANGE_HUOBI {
		key = "huobi:" + coinA.Name + "usdt"
	} else if exchangeId == coinObj.EXCHANGE_DRAGEONEX {
		key = "dragonex:" + coinB.Name + "usdt"
	}

	// 查询当前价格
	redis := util.Cli{}
	db, _ := beego.AppConfig.Int("redis_db_price")
	redis.Select(db)
	currentPriceStr, _ := redis.Get(key)
	currentPrice, _ := strconv.ParseFloat(currentPriceStr, 64)

	if taskType == sms.TYPE_THRESOLD_VALUE {
		// 判断任务阈值有效性
		if deviation == sms.DEVIATION_GT && value <= currentPrice {
			l.BadRequest("价格不能小于等于当前价格")
			return
		} else if deviation == sms.DEVIATION_LT && value >= currentPrice {
			l.BadRequest("价格不能大于等于当前价格")
			return
		}

	}

	walletModel := model.SmsWalletModel{}
	balance, err := walletModel.Balance(UID)
	if err != nil {
		beego.Error(err)
		l.Warn("添加预警失败")
		return
	}
	if balance <= 0 {
		l.Warn("短信钱包余额不足")
		return
	}

	task := sms.TaskThresholdValue{
		UID:            UID,
		CoinAId:        coinA.Id,
		CoinBId:        coinB.Id,
		SymbolPair:     symbolPair,
		ExchangeId:     exchangeId,
		ThresholdValue: strconv.FormatFloat(value, 'f', 15, 64),
		BaseValue:      currentPriceStr,
		Deviation:      deviation,
	}
	m := model.TaskThresholdValueModel{}
	err = m.Add(&task)
	if err != nil {
		beego.Error(err)
		l.Warn("添加预警失败")
		return
	}
}

func (l *TaskLogic) Get() {

}

func (l *TaskLogic) List() {

}

func (l *TaskLogic) Cancel() {

}
