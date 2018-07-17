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
	"github.com/astaxie/beego/orm"
)

type TaskLogic struct {
	logic.BaseLogic
}

// 添加预警任务
func (l *TaskLogic) Add(taskType, exchangeId int, symbolPair string, deviation int, value string) orm.Params {
	UID := l.GetUID()

	symbolPairSli := strings.Split(symbolPair, "_")
	coinModel := coin.CoinModel{}
	coinA, _ := coinModel.GetByName(symbolPairSli[0])
	coinB, _ := coinModel.GetByName(symbolPairSli[1])
	if coinA.Id == 0 || coinB.Id == 0 {
		l.BadRequest("交易对不存在")
		return orm.Params{}
	}

	key := ""
	if exchangeId == coinObj.EXCHANGE_HUOBI {
		key = "huobi:" + coinA.Name + "usdt"
	} else if exchangeId == coinObj.EXCHANGE_DRAGEONEX {
		key = "dragonex:" + coinA.Name + "usdt"
	}

	// 查询当前价格
	redis := util.Cli{}
	db, _ := beego.AppConfig.Int("redis_db_price")
	redis.Select(db)
	currentPriceStr, _ := redis.Get(key)
	currentPrice, _ := strconv.ParseFloat(currentPriceStr, 64)
	taskValue, _ := strconv.ParseFloat(value, 64)

	if taskType == sms.TYPE_THRESOLD_VALUE {
		// 判断任务阈值有效性
		if deviation == sms.DEVIATION_GT && taskValue <= currentPrice {
			l.BadRequest("价格不能小于等于当前价格")
			return orm.Params{}
		} else if deviation == sms.DEVIATION_LT && taskValue >= currentPrice {
			l.BadRequest("价格不能大于等于当前价格")
			return orm.Params{}
		}

	}

	walletModel := model.SmsWalletModel{}
	balance, err := walletModel.Balance(UID)
	if err != nil {
		beego.Error(err)
		l.Warn("添加预警失败")
		return orm.Params{}
	}
	if balance <= 0 {
		l.Warn("短信钱包余额不足")
		return orm.Params{}
	}

	task := sms.TaskThresholdValue{
		UID:            UID,
		CoinAId:        coinA.Id,
		CoinBId:        coinB.Id,
		SymbolPair:     symbolPair,
		ExchangeId:     exchangeId,
		ThresholdValue: value,
		BaseValue:      currentPriceStr,
		Deviation:      deviation,
	}
	m := model.TaskThresholdValueModel{}
	err = m.Add(&task)
	if err != nil {
		beego.Error(err)
		l.Warn("添加预警失败")
		return orm.Params{}
	}

	return orm.Params{"task_id": task.Id}
}

// 查询某一交易对任务列表
func (l *TaskLogic) Get(types, exchangeId int, symbolPair string) []orm.Params {
	UID := l.GetUID()

	symbolPairSli := strings.Split(symbolPair, "_")
	coinModel := coin.CoinModel{}
	coinA, _ := coinModel.GetByName(symbolPairSli[0])
	coinB, _ := coinModel.GetByName(symbolPairSli[1])
	if coinA.Id == 0 || coinB.Id == 0 {
		l.BadRequest("交易对不存在")
		return []orm.Params{}
	}

	m := model.SmsTaskModel{}
	res, err := m.Get(UID, types, exchangeId, symbolPair)
	if err != nil {
		beego.Error(err)
		l.Warn("查询失败")
		return res
	}

	return res
}

// 查询任务列表
func (l *TaskLogic) List() []orm.Params {
	UID := l.GetUID()

	m := model.SmsTaskModel{}
	res, err := m.GetList(UID)
	if err != nil {
		beego.Error(err)
		l.Warn("查询失败")
		return res
	}

	return res
}

// 取消任务
func (l *TaskLogic) Cancel(taskId int) {
	UID := l.GetUID()

	m := model.SmsTaskModel{}

	// 判断任务状态
	status, err := m.Status(UID, taskId)
	if err != nil {
		beego.Error(err)
		l.Warn("取消任务失败")
		return
	} else if status != sms.STATUS_WAIT {
		l.BadRequest("无法取消任务")
		return
	}

	// 取消任务
	err = m.Cancel(UID, taskId)
	if err != nil {
		beego.Error(err)
		l.Warn("取消任务失败")
		return
	}
}
